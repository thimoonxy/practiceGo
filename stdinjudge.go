package main

import (
	"fmt"

	"io"
	"io/ioutil"
	"os"
	"time"
)

func main() {

	s := StdJudge(os.Stdin)

	fmt.Println("main done.", string(s))

}

func StdJudge(s io.Reader) []byte {
	stdinreader := make(chan []byte)
	go func() {
		for {
			select {
			case buf, ok := <-stdinreader:
				if !ok {
					stdinreader = nil
				} else {
					//fmt.Println("here:", string(buf), ok)
				}
				if len(buf) != 0 {
					close(stdinreader)
				}
			case <-time.After(1 * time.Second):
				fmt.Println("timeout")
				os.Exit(1)
			}
		}
	}()
	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		os.Exit(1)
	}

	stdinreader <- buf
	return buf
}
