// Package unsafe defines a worker that does not lock the shared
// resource when working on it.
//
// As a consequence, when several workers of this type are asked to work
// cocurrently on the same resource, it will end up garbled.
package unsafe

import (
	"fmt"
	"io"
	"math/rand"
	"time"
)

const maxSleepMsecs = 100

// UnSafe implements a worker that does not care about locking the shared
// resource.
type UnSafe struct {
	name   int
	writer io.Writer
	done   chan<- bool
}

// Returns a new unsafe worker, named after the given number (see the
// Strig method).
//
// This worker will use the given writer as the shared resource and will
// notify when the work is done by sending true to the given channel.
func NewWorker(name int, writer io.Writer, done chan<- bool) *UnSafe {
	return &UnSafe{
		name:   name,
		writer: writer,
		done:   done,
	}
}

// Implements fmt.Stringer.  Workers are identified by their name, which
// is a number for easier use; it is returned here as part of the
// identification string of each worker for debugging purposes.
func (us *UnSafe) String() string {
	return fmt.Sprintf("worker %d", us.name)
}

// Work implements Worker.
func (us *UnSafe) Work() error {
	defer func() {
		us.done <- true
	}()
	//fmt.Printf("[%s] starting to work\n", us)
	//defer fmt.Printf("[%s] finished working\n", us)
	if _, err := us.writer.Write([]byte{'<'}); err != nil {
		return fmt.Errorf("%s: writting '<': %s", us, err)
	}
	randSleep(maxSleepMsecs)
	if _, err := us.writer.Write([]byte{'>'}); err != nil {
		return fmt.Errorf("%s: writting '>': %s", us, err)
	}
	return nil
}

func randSleep(msecs int32) {
	msec := rand.Int31n(msecs)
	time.Sleep(time.Duration(msec) * time.Millisecond)
}
