// Package safe defines a worker that lock the shared resource when
// working on it, so several cooperating workers of this type can
// work at the same time on the same resource with garbling it.
package safe

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	flock "github.com/theckman/go-flock"
)

const maxSleepMsecs = 100

// Safe implements a concurrent safe worker.
type Safe struct {
	name   int
	writer io.Writer
	lock   *flock.Flock
}

// Returns a new concurrent safe worker, named after the given number
// (see the Strig method).  It will use the given writer as the shared
// resource and the given lock to avoid data races on the resource.
func NewWorker(name int, writer io.Writer, resourceName string) *Safe {
	return &Safe{
		name:   name,
		writer: writer,
		lock:   flock.NewFlock(resourceName + ".lock"),
	}
}

// Implements fmt.Stringer.  Workers are identified by their name, which
// is a number for easier use; it is returned here as part of the
// identification string of each worker for debugging purposes.
func (s *Safe) String() string {
	return fmt.Sprintf("worker %d", s.name)
}

// Work implements Worker.
func (s *Safe) Work(done chan<- bool) error {
	defer func() {
		done <- true
	}()
	for {
		locked, err := s.lock.TryLock()
		if err != nil {
			return fmt.Errorf("%s: trying to lock: %s", s, err)
		}
		if locked {
			break
		}
		randSleep(10)
	}
	//fmt.Printf("[%s] starting to work\n", us)
	//defer fmt.Printf("[%s] finished working\n", us)
	if _, err := s.writer.Write([]byte{'<'}); err != nil {
		return fmt.Errorf("%s: writting '<': %s", s, err)
	}
	randSleep(maxSleepMsecs)
	if _, err := s.writer.Write([]byte{'>'}); err != nil {
		return fmt.Errorf("%s: writting '>': %s", s, err)
	}
	if err := s.lock.Unlock(); err != nil {
		return fmt.Errorf("%s: unlocking: %s", s, err)
	}
	return nil
}

func randSleep(msecs int32) {
	msec := rand.Int31n(msecs)
	time.Sleep(time.Duration(msec) * time.Millisecond)
}
