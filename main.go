package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var writer = os.Stdout
var cpus = runtime.NumCPU()

func main() {
	runtime.GOMAXPROCS(cpus)
	dirwalk(string("."))
}

func dirwalk(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	semaohorer := make(chan int, cpus)
	wg := sync.WaitGroup{}
	for _, file := range files {
		if file.IsDir() {
			wg.Add(1)
			go deepDirwalk(dir, file, &wg, semaohorer)
		}
	}
	wg.Wait()
}

func deepDirwalk(dir string, file os.FileInfo, wg *sync.WaitGroup, ch chan int) {
	defer wg.Done()
	ch <- 1
	path := filepath.Join(dir, file.Name())
	fmt.Fprintln(writer, path)
	dirwalk(path)
	<-ch
}
