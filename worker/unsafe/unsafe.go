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

// Unsafe implements a worker that does not care about locking the shared
// resource.
type Unsafe struct {
	name   int
	writer io.Writer
}

// Returns a new unsafe worker, named after the given number (see the
// Strig method).  It will use the given writer as the shared resource.
func NewWorker(name int, writer io.Writer) *Unsafe {
	return &Unsafe{
		name:   name,
		writer: writer,
	}
}

// Implements fmt.Stringer.  Workers are identified by their name, which
// is a number for easier use; it is returned here as part of the
// identification string of each worker for debugging purposes.
func (us *Unsafe) String() string {
	return fmt.Sprintf("worker %d", us.name)
}

// Work implements Worker.
func (us *Unsafe) Work(done chan<- error) {
	var err error
	defer func() {
		done <- err
	}()
	//fmt.Printf("[%s] starting to work\n", us)
	//defer fmt.Printf("[%s] finished working\n", us)
	if _, err := us.writer.Write([]byte{'<'}); err != nil {
		err = fmt.Errorf("%s: writting '<': %s", us, err)
		return
	}
	randSleep(maxSleepMsecs)
	if _, err := us.writer.Write([]byte{'>'}); err != nil {
		err = fmt.Errorf("%s: writting '>': %s", us, err)
		return
	}
	return
}

func randSleep(msecs int32) {
	msec := rand.Int31n(msecs)
	time.Sleep(time.Duration(msec) * time.Millisecond)
}
