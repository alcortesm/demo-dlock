package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

const (
	nWriters       = 2
	tempFilePrefix = "demo-dlock-"
)

func main() {
	do()
}

func do() {
	file, err := tempFile()
	if err != nil {
		log.Fatal("cannot create temporary file:", err)
	}

	done := make(chan bool, nWriters)

	for i := 0; i < nWriters; i++ {
		w := writer{
			name: i,
			done: done,
		}
		go w.write(file)
	}

	for i := 0; i < nWriters; i++ {
		<-done
	}

	file.Close()

	if err = dump(file.Name()); err != nil {
		log.Fatal("cannot dump temporary file:", err)
	}
}

func tempFile() (*os.File, error) {
	const useDefaultTempDir = ""
	path, err := ioutil.TempFile(useDefaultTempDir, tempFilePrefix)
	if err != nil {
		return nil, err
	}

	return path, nil
}

type writer struct {
	name int
	done chan<- bool
}

func (w writer) String() string {
	return fmt.Sprintf("-%d-", w.name)
}

func (w writer) write(file *os.File) error {
	defer func() {
		w.done <- true
	}()
	fmt.Printf("[writer %d] starting to write\n", w.name)
	defer fmt.Printf("[writer %d] finished writing\n", w.name)
	for i := 0; i < 10; i++ {
		if _, err := file.WriteString(w.String()); err != nil {
			return fmt.Errorf("writer %d: %s", w.name, err)
		}
		randSleep()
	}
	return nil
}

func randSleep() {
	msec := rand.Int31n(10)
	time.Sleep(time.Duration(msec) * time.Millisecond)
}

func dump(file *os.File) err {
	data, err ;= ioutil.ReadAll()
}
