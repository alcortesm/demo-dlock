package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alcortesm/demo-dlock/dlock/etcd"
	"github.com/alcortesm/demo-dlock/dlock/flock"
	"github.com/alcortesm/demo-dlock/worker"
	"github.com/alcortesm/demo-dlock/worker/safe"
	"github.com/alcortesm/demo-dlock/worker/unsafe"
)

const (
	nRuns          = 2
	nWriters       = 3
	tempFilePrefix = "demo-dlock-"
	argUnsafe      = "unsafe"
	argFlock       = "flock"
	argEtcd        = "etcd"
)

func main() {
	implementation := parseArgs()
	fmt.Printf("Running %d experiments with %d writers each, impl=",
		nRuns, nWriters)
	fmt.Println(implementation)
	nRunsFinished := 0
	nGarbled := 0
	for i := 0; i < nRuns; i++ {
		fmt.Printf("run %d: ", i)
		file, err := tempFile()
		if err != nil {
			fmt.Printf("creating temp file: %s\n", err)
			continue
		}
		fmt.Printf("%s: ", file.Name())

		err = run(implementation, file)
		if err != nil {
			fmt.Println("ERROR", err)
			continue
		}
		if err := file.Close(); err != nil {
			fmt.Printf("ERROR closing the temp file: %s\n", err)
			continue
		}
		garbled, err := worker.IsGarbled(file.Name())
		if err != nil {
			fmt.Printf("ERROR checking garbled: %s\n", err)
		}
		nRunsFinished++
		if garbled {
			fmt.Println("FAILED text is garbled")
			nGarbled++
		} else {
			fmt.Println("SUCCESS")
		}
	}
	if nGarbled != 0 {
		fmt.Printf("FAILED: %d garbled resources\n", nGarbled)
		os.Exit(1)
	}
	if nRunsFinished != nRuns {
		fmt.Printf("ERROR: only %d runs finished (out of %d)\n",
			nRunsFinished, nRuns)
		os.Exit(1)
	}
	fmt.Println("OK: all experiments were successfull")
	os.Exit(0)
}

func parseArgs() (implementation string) {
	if len(os.Args) != 2 {
		usage()
	}
	switch os.Args[1] {
	case argUnsafe, argFlock, argEtcd:
		return os.Args[1]
	default:
		usage()
	}
	panic("unreachable")
}

func usage() {
	fmt.Println("usage:")
	fmt.Printf("\t%s [%s, %s, %s]\n",
		os.Args[0], argUnsafe, argFlock, argEtcd)
	os.Exit(1)
}

func tempFile() (*os.File, error) {
	const useDefaultTempDir = ""
	path, err := ioutil.TempFile(useDefaultTempDir, tempFilePrefix)
	if err != nil {
		return nil, err
	}

	return path, nil
}

// runs several workers over the same sared resource of the given
// implementation.
func run(implementation string, shared *os.File) error {
	done := make(chan error, nWriters)
	for i := 0; i < nWriters; i++ {
		var w worker.Worker
		switch implementation {
		case argUnsafe:
			w = unsafe.NewWorker(i, shared)
		case argFlock:
			l := flock.NewDLock(shared.Name())
			w = safe.NewWorker(i, shared, l)
		case argEtcd:
			l := etcd.NewDLock(shared.Name())
			w = safe.NewWorker(i, shared, l)
		default:
			return fmt.Errorf("unkown implementation: %s", implementation)
		}
		go w.Work(done)
	}

	for i := 0; i < nWriters; i++ {
		if err := <-done; err != nil {
			return err
		}
	}

	// adding a EOL at the end helps when humans want to read the files
	if _, err := shared.Write([]byte{'\n'}); err != nil {
		return err
	}
	return nil
}
