// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/4/3

package etcdx

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kratos/kratos/v2/registry"
	"go.etcd.io/etcd/client/v3"
)

type etcdWatcher struct {
	watchChan clientv3.WatchChan
	cancel    context.CancelFunc
}

func (e *etcdWatcher) Next() ([]*registry.ServiceInstance, error) {
	for resp := range e.watchChan {
		if resp.Err() != nil {
			return nil, resp.Err()
		}
		if resp.Canceled {
			return nil, errors.New("could not get next")
		}
		serviceList := make([]*registry.ServiceInstance, 0, len(resp.Events))
		for _, ev := range resp.Events {
			service := decode(ev.Kv.Value)
			var action string
			switch ev.Type {
			case clientv3.EventTypePut:
				if ev.IsCreate() {
					action = "create"
				} else if ev.IsModify() {
					action = "update"
				}
			case clientv3.EventTypeDelete:
				action = "delete"

				// get service from prevKv
				service = decode(ev.PrevKv.Value)
			}
			fmt.Println(action)
			if service == nil {
				continue
			}
			serviceList = append(serviceList, service)
		}
	}
	return nil, errors.New("could not get next")
}

func (e *etcdWatcher) Stop() error {
	e.cancel()
	return nil
}

func newEtcdWatcher(ctx context.Context, r *etcdRegistry, key string) registry.Watcher {
	newCtx, cancel := context.WithCancel(ctx)
	return &etcdWatcher{
		cancel:    cancel,
		watchChan: r.client.Watch(newCtx, key, clientv3.WithPrefix(), clientv3.WithPrevKV()),
	}
}
