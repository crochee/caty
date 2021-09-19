// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/3

package etcdx

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (e *etcdRegistry) GetService(ctx context.Context, serviceName string) ([]*registry.ServiceInstance, error) {
	newCtx, cancel := context.WithTimeout(ctx, e.Option.Timeout)
	defer cancel()

	rsp, err := e.client.Get(newCtx, e.servicePath(serviceName)+"/", clientv3.WithPrefix(), clientv3.WithSerializable())
	if err != nil {
		return nil, err
	}

	if len(rsp.Kvs) == 0 {
		return nil, errors.New("service not found")
	}

	serviceMap := map[string]*registry.ServiceInstance{}

	for _, n := range rsp.Kvs {
		if sn := decode(n.Value); sn != nil {
			s, ok := serviceMap[sn.Version]
			if !ok {
				s = &registry.ServiceInstance{
					ID:        sn.ID,
					Name:      sn.Name,
					Version:   sn.Version,
					Metadata:  sn.Metadata,
					Endpoints: sn.Endpoints,
				}
				serviceMap[s.Version] = s
			}
		}
	}

	services := make([]*registry.ServiceInstance, 0, len(serviceMap))
	for _, service := range serviceMap {
		services = append(services, service)
	}
	return services, nil
}

func (e *etcdRegistry) Watch(ctx context.Context, serviceName string) (registry.Watcher, error) {
	return newEtcdWatcher(ctx, e, e.servicePath(serviceName)), nil
}
