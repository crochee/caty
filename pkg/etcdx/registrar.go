// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/2

package etcdx

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"hash/fnv"
	"net"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/json-iterator/go"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"

	"obs/pkg/registry"
)

type Option struct {
	Prefix    string
	AddrList  []string
	Timeout   time.Duration
	Secure    bool
	TLSConfig *tls.Config
	zapConfig *zap.Config
	Context   context.Context
	Username  string
	Password  string
	TTl       time.Duration
}

func NewEtcdRegistry(opts ...func(*Option)) (*etcdRegistry, error) {
	e := &etcdRegistry{
		Option: Option{
			Prefix: "/micro/registry/",
		},
		register: make(map[string]uint64),
		leases:   make(map[string]clientv3.LeaseID),
	}
	cf := clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Context:   context.Background(),
	}

	for _, opt := range opts {
		opt(&e.Option)
	}
	if e.Option.Context != nil {
		cf.Context = e.Option.Context
	}

	if e.Option.Timeout == 0 {
		e.Option.Timeout = 5 * time.Second
	}

	if e.Option.Secure {
		if e.Option.TLSConfig == nil {
			e.Option.TLSConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		cf.TLS = e.Option.TLSConfig
	}

	cf.Username = e.Option.Username
	cf.Password = e.Option.Password
	cf.LogConfig = e.Option.zapConfig

	var addrList []string

	for _, address := range e.Option.AddrList {
		if len(address) == 0 {
			continue
		}
		addr, port, err := net.SplitHostPort(address)
		if ae, ok := err.(*net.AddrError); ok && ae.Err == "missing port in address" {
			port = "2379"
			addr = address
			addrList = append(addrList, net.JoinHostPort(addr, port))
		} else if err == nil {
			addrList = append(addrList, net.JoinHostPort(addr, port))
		}
	}
	// if we got addrList then we'll update
	if len(addrList) > 0 {
		cf.Endpoints = addrList
	}

	cli, err := clientv3.New(cf)
	if err != nil {
		return nil, err
	}
	e.client = cli
	return e, nil
}

type etcdRegistry struct {
	client *clientv3.Client
	Option
	sync.RWMutex
	register map[string]uint64
	leases   map[string]clientv3.LeaseID
}

func (e *etcdRegistry) Register(ctx context.Context, service *registry.ServiceInstance) error {
	if len(service.Endpoints) == 0 {
		return errors.New("require at least one node")
	}

	// check existing lease cache
	e.RLock()
	leaseID, ok := e.leases[service.Name+service.ID]
	e.RUnlock()

	newCtx, cancel := context.WithTimeout(ctx, e.Option.Timeout)
	defer cancel()
	// 查询leaseId 不存在就查询etcd
	if !ok {
		// missing lease, check if the key exists
		// look for the existing key
		rsp, err := e.client.Get(newCtx, e.nodePath(service.Name, service.ID), clientv3.WithSerializable())
		if err != nil {
			return err
		}

		// get the existing lease
		for _, kv := range rsp.Kvs {
			if kv.Lease > 0 {
				leaseID = clientv3.LeaseID(kv.Lease)

				// decode the existing node
				srv := decode(kv.Value)
				if srv == nil || len(srv.Endpoints) == 0 {
					continue
				}

				// create hash of service; uint64
				h, err := Hash(srv.Endpoints)
				if err != nil {
					continue
				}

				// save the info
				e.Lock()
				e.leases[service.Name+service.ID] = leaseID
				e.register[service.Name+service.ID] = h
				e.Unlock()
				break
			}
		}
	}

	var leaseNotFound bool

	// renew the lease if it exists
	if leaseID > 0 {
		e.client.GetLogger().Debug(fmt.Sprintf("renewing existing lease for %s %d", service.Name, leaseID))
		if _, err := e.client.KeepAliveOnce(newCtx, leaseID); err != nil {
			if err != rpctypes.ErrLeaseNotFound {
				return err
			}
			e.client.GetLogger().Debug(fmt.Sprintf("lease not found for %s %d", service.Name, leaseID))
			// lease not found do register
			leaseNotFound = true
		}
	}

	// create hash of service; uint64
	h, err := Hash(service.Endpoints)
	if err != nil {
		return err
	}

	// get existing hash for the service node
	e.Lock()
	v, ok := e.register[service.Name+service.ID]
	e.Unlock()

	// the service is unchanged, skip registering
	if ok && v == h && !leaseNotFound {
		e.client.GetLogger().Debug(fmt.Sprintf("service %s node %s unchanged skipping registration",
			service.Name, service.ID))
		return nil
	}

	var lgr *clientv3.LeaseGrantResponse
	if second := e.Option.TTl.Seconds(); second > 0 {
		// get a lease used to expire keys since we have a ttl
		lgr, err = e.client.Grant(newCtx, int64(second))
		if err != nil {
			return err
		}
		e.client.GetLogger().Debug(fmt.Sprintf("Registering %s id %s with lease %v and leaseID %v and ttl %v",
			service.Name, service.ID, lgr, lgr.ID, e.Option.TTl))
	}

	// create an entry for the node
	if lgr != nil {
		_, err = e.client.Put(newCtx, e.nodePath(service.Name, service.ID),
			encode(service), clientv3.WithLease(lgr.ID))
	} else {
		_, err = e.client.Put(newCtx, e.nodePath(service.Name, service.ID), encode(service))
	}
	if err != nil {
		return err
	}

	e.Lock()
	// save our hash of the service
	e.register[service.Name+service.ID] = h
	// save our leaseID of the service
	if lgr != nil {
		e.leases[service.Name+service.ID] = lgr.ID
	}
	e.Unlock()
	return nil
}

func (e *etcdRegistry) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	if len(service.Endpoints) == 0 {
		return errors.New("require at least one node")
	}

	e.Lock()
	// delete our hash of the service
	delete(e.register, service.Name+service.ID)
	// delete our lease of the service
	delete(e.leases, service.Name+service.ID)
	e.Unlock()

	newCtx, cancel := context.WithTimeout(ctx, e.Option.Timeout)
	defer cancel()

	e.client.GetLogger().Debug(fmt.Sprintf("deregistering %s id %s", service.Name, service.ID))

	_, err := e.client.Delete(newCtx, e.nodePath(service.Name, service.ID))
	if err != nil {
		return err
	}
	return nil
}

func (e *etcdRegistry) nodePath(s, id string) string {
	service := strings.Replace(s, "/", "-", -1)
	node := strings.Replace(id, "/", "-", -1)
	return path.Join(e.Option.Prefix, service, node)
}

func (e *etcdRegistry) servicePath(s string) string {
	return path.Join(e.Option.Prefix, strings.Replace(s, "/", "-", -1))
}

func encode(s *registry.ServiceInstance) string {
	b, _ := jsoniter.ConfigFastest.MarshalToString(s)
	return b
}

func decode(ds []byte) *registry.ServiceInstance {
	var s registry.ServiceInstance
	_ = jsoniter.ConfigFastest.Unmarshal(ds, &s)
	return &s
}

func Hash(value interface{}) (uint64, error) {
	data, err := jsoniter.Marshal(value)
	if err != nil {
		return 0, err
	}
	w := fnv.New64()
	if _, err = w.Write(data); err != nil {
		return 0, err
	}
	return w.Sum64(), nil
}
