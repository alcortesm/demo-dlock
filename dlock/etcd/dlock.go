package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
)

type DLock struct {
	session   *concurrency.Session
	mutex     *concurrency.Mutex
	lockLease time.Duration
}

// TODO use namespaces to avoid pfx collisions.

func NewDLock(client *clientv3.Client, pfx string, lockLease time.Duration) (*DLock, error) {
	session, err := concurrency.NewSession(client,
		concurrency.WithContext(client.Ctx()))
	if err != nil {
		return nil, fmt.Errorf("creating session for etcd client")
	}
	return &DLock{
		session:   session,
		mutex:     concurrency.NewMutex(session, pfx),
		lockLease: lockLease,
	}, nil
}

func (dl *DLock) Lock() error {
	ctx, cancel := context.WithTimeout(context.Background(), dl.lockLease)
	defer func() {
		cancel()
	}()
	return dl.mutex.Lock(ctx)
}

func (dl *DLock) Unlock() error {
	return dl.mutex.Unlock(context.Background())
}

// this should go into the DLock interface.
func (dl *DLock) Close() error {
	return dl.session.Close()
}
