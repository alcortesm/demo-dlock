package worker

import (
	"fmt"
	"io/ioutil"
)

// Workers will work on a resource that might be shared among several of
// them.
//
// This interface does not make any assumptions about how safe,
// concurrency wise, are the workers.
type Worker interface {
	// Work will write first '<', then '>' to the shared resource of the
	// worker.  At the end of a successfull execution it will send a nil
	// over the given channel, or an error otherwise.
	//
	// Note how if several non-cooperative workers try to work at the
	// same time with the same resource, the contents there will end up
	// garbled ("<<>>", instead of "<><>" for the case of two workers).
	Work(done chan<- error)
}

// IsGarbled returns if the contents of the file at path are the result
// of several workers interleaving their work or not.  It returns an
// error if the contents of the file cannot be accessed or nil
// otherwise.
func IsGarbled(path string) (bool, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	if len(data) == 0 {
		return false, fmt.Errorf("empty file")
	}
	return isGarbled(data), nil
}

func isGarbled(data []byte) bool {
	lastWasGT := true // when starting, assume last byte was '>'
	for _, current := range data {
		if current == '\n' {
			break
		}
		if lastWasGT {
			if current != '<' {
				return true
			}
			lastWasGT = false
		} else {
			if current != '>' {
				return true
			}
			lastWasGT = true
		}
	}
	return false
}
