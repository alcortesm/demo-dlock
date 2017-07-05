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
)

const maxSleepMsecs = 100

// Safe implements a concurrent safe worker.
type Safe struct {
	unsafe *unsafe.Unsafe
	lock   dlock.DLocker
}

// Returns a new concurrent safe worker, named after the given number
// (see the Strig method).  It will use the given writer as the shared
// resource and the given lock to avoid data races on the resource.
func NewWorker(name int, writer io.Writer, lock dlock.DLocker) *Safe {
	return &Safe{
		unsafe: unsafe.NewWorker(name, writer),
		lock:   lock,
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
	err = s.lock.Lock()
	if err != nil {
		err = fmt.Errorf("%s: trying to lock: %s", s, err)
		return
	}
	defer func() {
		if errUnlock := s.lock.Unlock(); err == nil && errUnlock != nil {
			err = fmt.Errorf("%s: unlocking: %s", s, errUnlock)
		}
	}()

	inner := make(chan error)
	go s.unsafe.Work(inner)
	if err = <-inner; err != nil {
		return
	}
}

func randSleep(msecs int32) {
	msec := rand.Int31n(msecs)
	time.Sleep(time.Duration(msec) * time.Millisecond)
}
