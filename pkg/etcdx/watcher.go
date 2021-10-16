package etcdx

import (
	"context"
	"errors"
	"fmt"

	"go.etcd.io/etcd/client/v3"

	"cca/pkg/registry"
)

type etcdWatcher struct {
	watchChan clientv3.WatchChan
	cancel    context.CancelFunc
	ctx       context.Context
}

func (e *etcdWatcher) Next() ([]*registry.ServiceInstance, error) {
	select {
	case resp := <-e.watchChan: // todo 优化
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
		return serviceList, nil
	case <-e.ctx.Done():
		return nil, e.ctx.Err()
	}
}

func (e *etcdWatcher) Stop() error {
	e.cancel()
	return nil
}

func newEtcdWatcher(ctx context.Context, r *etcdRegistry, key string) registry.Watcher {
	newCtx, cancel := context.WithCancel(ctx)
	return &etcdWatcher{
		cancel:    cancel,
		ctx:       newCtx,
		watchChan: r.client.Watch(newCtx, key, clientv3.WithPrefix(), clientv3.WithPrevKV()),
	}
}
