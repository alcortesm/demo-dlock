package flock

import (
	"fmt"
	"math/rand"
	"time"

	original "github.com/theckman/go-flock"
)

const pollingMsecs = 10

type DLock struct {
	original.Flock
}

func NewDLock(path string) *DLock {
	return &DLock{Flock: *original.NewFlock(path + ".lock")}
}

func (dl *DLock) Lock() error {
	for {
		locked, err := dl.Flock.TryLock()
		if err != nil {
			return fmt.Errorf("trying to lock: %s", err)
		}
		if locked {
			break
		}
		randSleep(pollingMsecs)
	}
	return nil
}

func randSleep(msecs int32) {
	msec := rand.Int31n(msecs)
	time.Sleep(time.Duration(msec) * time.Millisecond)
}
