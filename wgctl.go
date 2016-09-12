package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	sum := 6
	onetime := 3
	wg := sync.WaitGroup{}
	for i := 0; i < sum/onetime; i++ {
		wg.Add(onetime)
		for j := i * onetime; j < onetime*(i+1); j++ {
			go Curl(&wg, j)
		}
		wg.Wait()
	}
}

func Curl(wg *sync.WaitGroup, id int) {
	fmt.Println(id, time.Now().Format(time.RFC3339Nano), "start")
	time.Sleep(2 * time.Second)
	fmt.Println(id, time.Now().Format(time.RFC3339Nano), "end")
	wg.Done()
}
