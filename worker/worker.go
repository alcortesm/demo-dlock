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
	// Work will write first '<', then '>' to the shared resource.
	//
	// Note how if several non-concurrent safe workers try to work at
	// the same resource, the contents there will end up garbled
	// ("<<>>", instead of "<><>")
	Work() error
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
	lastWasLT := true
	for _, current := range data {
		if lastWasLT {
			if current != '<' {
				return true
			}
			lastWasLT = false
		} else {
			if current != '>' {
				return true
			}
			lastWasLT = true
		}
	}
	return false
}
