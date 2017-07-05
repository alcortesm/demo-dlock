// Package safe defines a worker that lock the shared resource when
// working on it, so several cooperating workers of this type can
// work at the same time on the same resource with garbling it.
package safe

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/alcortesm/demo-dlock/dlock"
	"github.com/alcortesm/demo-dlock/worker/unsafe"
	flock "github.com/theckman/go-flock"
)

const maxSleepMsecs = 100

// Safe implements a concurrent safe worker.
type Safe struct {
	unsafe *unsafe.Unsafe
	lock   dlock.DLock
}

// Returns a new concurrent safe worker, named after the given number
// (see the Strig method).  It will use the given writer as the shared
// resource and the given lock to avoid data races on the resource.
func NewWorker(name int, writer io.Writer, resourceName string) *Safe {
	return &Safe{
		unsafe: unsafe.NewWorker(name, writer),
		lock:   flock.NewFlock(resourceName + ".lock"),
	}
}

// Implements fmt.Stringer.  Workers are identified by their name, which
// is a number for easier use; it is returned here as part of the
// identification string of each worker for debugging purposes.
func (s *Safe) String() string {
	return s.unsafe.String()
}

// Work implements Worker.
func (s *Safe) Work(done chan<- error) {
	var err error
	defer func() {
		done <- err
	}()
	for {
		locked, err := s.lock.TryLock()
		if err != nil {
			err = fmt.Errorf("%s: trying to lock: %s", s, err)
			return
		}
		if locked {
			break
		}
		randSleep(10)
	}
	inner := make(chan error)
	go s.unsafe.Work(inner)
	if err = <-inner; err != nil {
		return
	}
	if err = s.lock.Unlock(); err != nil {
		err = fmt.Errorf("%s: unlocking: %s", s, err)
		return
	}
}

func randSleep(msecs int32) {
	msec := rand.Int31n(msecs)
	time.Sleep(time.Duration(msec) * time.Millisecond)
}
