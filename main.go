package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alcortesm/demo-dlock/worker"
	"github.com/alcortesm/demo-dlock/worker/safe"
	"github.com/alcortesm/demo-dlock/worker/unsafe"
)

const (
	nRuns          = 10
	nWriters       = 10
	tempFilePrefix = "demo-dlock-"
)

func main() {
	wantSafe := true
	wantEtcd := false
	nRunsFinished := 0
	nGarbled := 0

	fmt.Printf("Running %d experiments with %d writers each, want safe?=%t",
		nRuns, nWriters, wantSafe)
	if wantSafe {
		if wantEtcd {
			fmt.Println(", using etcd")
		} else {
			fmt.Println(", using flock")
		}
	} else {
		fmt.Println()
	}

	for i := 0; i < nRuns; i++ {
		fmt.Printf("run %d: ", i)
		file, err := tempFile()
		if err != nil {
			fmt.Printf("creating temp file: %s\n", err)
			continue
		}
		fmt.Printf("%s: ", file.Name())

		err = run(wantSafe, wantEtcd, file)
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

func tempFile() (*os.File, error) {
	const useDefaultTempDir = ""
	path, err := ioutil.TempFile(useDefaultTempDir, tempFilePrefix)
	if err != nil {
		return nil, err
	}

	return path, nil
}

// runs several workers over the same sared resource.  The wantSafe
// arguments controls wether to use safe workers or unsafe ones.
func run(wantSafe, wantEtcd bool, shared *os.File) error {
	done := make(chan error, nWriters)
	for i := 0; i < nWriters; i++ {
		var w worker.Worker
		if wantSafe {
			if wantEtcd {
				return fmt.Errorf("TODO no etcd implementation yet")
			} else {
				w = safe.NewWorker(i, shared, shared.Name())
			}
		} else {
			w = unsafe.NewWorker(i, shared)
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
