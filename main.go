package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/alcortesm/demo-dlock/worker"
	"github.com/alcortesm/demo-dlock/worker/unsafe"
)

const (
	nRuns          = 10
	nWriters       = 2
	tempFilePrefix = "demo-dlock-"
)

func main() {
	safe := false
	nGarbled := 0
	fmt.Printf("Running %d experiments with %d writers each\n",
		nRuns, nWriters)
	for i := 0; i < nRuns; i++ {
		garbled, path, err := run(safe)
		if err != nil {
			log.Fatalf("run %d: %s\n", i, err)
		}
		if garbled {
			fmt.Fprintf(os.Stderr,
				"run %d: text is garbled: %s\n", i, path)
			nGarbled++
		}
	}
	if nGarbled != 0 {
		fmt.Printf("FAILED: %d garbled resources\n", nGarbled)
		os.Exit(1)
	}
	fmt.Println("OK: all experiments were succesfull")
	os.Exit(0)
}

// returns if the text is garbled after the running several workers in
// parallel and the path to the temporal file used as a shared resource.
func run(safe bool) (garbled bool, path string, err error) {
	file, err := tempFile()
	if err != nil {
		return false, "", fmt.Errorf("cannot create temp file: ", err)
	}

	done := make(chan bool, nWriters)

	for i := 0; i < nWriters; i++ {
		var w worker.Worker
		if safe {
			return false, file.Name(),
				fmt.Errorf("TODO safe workers not implemented yet")
		} else {
			w = unsafe.NewWorker(i, file)
		}
		go w.Work(done)
	}

	for i := 0; i < nWriters; i++ {
		<-done
	}

	if _, err := file.Write([]byte{'\n'}); err != nil {
		return false, file.Name(), err
	}
	if err := file.Close(); err != nil {
		return false, file.Name(), err
	}

	garbled, err = worker.IsGarbled(file.Name())
	if err != nil {
		return false, file.Name(), fmt.Errorf(
			"%s: cannot check if garbled: %s", file.Name(), err)
	}
	return garbled, file.Name(), nil
}

func tempFile() (*os.File, error) {
	const useDefaultTempDir = ""
	path, err := ioutil.TempFile(useDefaultTempDir, tempFilePrefix)
	if err != nil {
		return nil, err
	}

	return path, nil
}
