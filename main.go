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
	// CPUの数を取得
	runtime.GOMAXPROCS(cpus)
	// 自分のいるディレクトリを取得
	// @TODO 引数検索するディレクトリを指定する
	dir, err := os.Getwd()
	if err != nil {
		// log出力して強制終了
		log.Fatal(err)
	}
	// ディレクトリ内を探索
	dirwalk(dir)
}

// export dirwalk
func dirwalk(dir string) {
	// 指定されたディレクトリもある fileInfoの一覧を取得する
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	// CPU 分のchannelを作る
	ch := make(chan int, cpus)
	wg := sync.WaitGroup{}
	// ファイルの一覧をループ
	for _, file := range files {
		if file.IsDir() {
			// ディレクトリだったら
			// WaitGroupを追加して、goroutine発行し、自分自身は処理を続行
			wg.Add(1)
			go deepDirwalk(dir, file, &wg, ch)
		} else {
			// ファイルだったら、標準出力して終わり
			path := filepath.Join(dir, file.Name())
			fmt.Fprintln(writer, path)
		}
	}
	// WaitGroupが0になるまで待つ
	wg.Wait()
}

func deepDirwalk(dir string, file os.FileInfo, wg *sync.WaitGroup, ch chan int) {
	// deepDirwalkが終わったら、WaitGroupをへらす
	defer wg.Done()
	// channelに値を入れる
	ch <- 1
	path := filepath.Join(dir, file.Name())
	fmt.Fprintln(writer, path)
	dirwalk(path)
	// channelの値を捨ててブロック解除
	<-ch
}
