package main

import (
	"C"
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
	// CPUの数を取得
	runtime.GOMAXPROCS(cpus)
	// 自分のいるディレクトリを取得
	// @TODO 引数検索するディレクトリを指定する
	dir, err := os.Getwd()
	if err != nil {
		//
		log.Fatal(err)
	}
	// ディレクトリ内を探索
	dirwalk(dir)
}

// export dirwalk
func dirwalk(dir string) string {
	// 指定されたディレクトリもある fileInfoの一覧を取得する
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	// CPU
	ch := make(chan int, cpus)
	wg := sync.WaitGroup{}
	// ファイルの一覧をループ
	for _, file := range files {
		// ディレクトリだったら
		if file.IsDir() {
			wg.Add(1)
			// 更に中を
			go deepDirwalk(dir, file, &wg, ch)
		} else {
			path := filepath.Join(dir, file.Name())
			fmt.Println(path)
		}
	}
	wg.Wait()
	return ""
}

func deepDirwalk(dir string, file os.FileInfo, wg *sync.WaitGroup, ch chan int) {
	defer wg.Done()
	ch <- 1
	path := filepath.Join(dir, file.Name())
	fmt.Fprintln(writer, path)
	dirwalk(path)
	<-ch
}
