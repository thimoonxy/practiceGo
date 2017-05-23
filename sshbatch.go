package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	sum := 30
	onetime := 13
	stdout := make(chan string, sum)
	wg := sync.WaitGroup{}

	// main runs
	if sum <= onetime {
		onetime = sum
	}
	for i := 0; i < sum/onetime; i++ {
		wg.Add(onetime)
		for j := i * onetime; j < onetime*(i+1); j++ {
			go func() {
				stdout <- SSH("root", "xxx", "192.168.141.67:22", &wg)
			}()
		}

		wg.Wait()
		for j := i * onetime; j < onetime*(i+1); j++ {
			fmt.Println(<-stdout)
		}
	}

	// remainder runs
	if sum%onetime != 0 {
		stdout2 := make(chan string, sum%onetime)
		wg2 := sync.WaitGroup{}
		wg2.Add(sum % onetime)
		for x := sum / onetime * onetime; x < sum; x++ {
			go func() {
				stdout2 <- SSH("root", "xxx", "192.168.141.67:22", &wg2)
			}()
		}
		wg2.Wait()
		for x := sum / onetime * onetime; x < sum; x++ {
			fmt.Println(<-stdout2)
		}
	}
}

func SSH(user, password, ip_port string, wg *sync.WaitGroup) string {

	var r *bufio.Reader
	PassWd := []ssh.AuthMethod{ssh.Password(password)}
	Conf := ssh.ClientConfig{User: user, Auth: PassWd}
	Client, err := ssh.Dial("tcp", ip_port, &Conf)
	if err != nil {
		wg.Done()
		return err.Error()
	}
	defer Client.Close()

	if len(os.Args) >= 2 {
		s := strings.NewReader(os.Args[1])
		r = bufio.NewReader(s)
	} else {
		r = bufio.NewReader(os.Stdin)
	}

	for {
		b, _, e := r.ReadLine()
		if e == io.EOF {
			wg.Done()
			return ""
		}

		if e != nil {
			fmt.Println("no stdin or args.")
			wg.Done()
			return e.Error()
		}
		command := string(b)
		if session, err := Client.NewSession(); err == nil {
			defer session.Close()
			// session.Stdout = os.Stdout
			// session.Stderr = os.Stderr
			var stdoutBuf, stderrBuf bytes.Buffer
			session.Stdout = &stdoutBuf
			session.Stderr = &stderrBuf
			session.Run(command)
			if len(stderrBuf.String()) > 0 {
				wg.Done()
				return stdoutBuf.String() + stderrBuf.String()
			} else {
				wg.Done()
				return stdoutBuf.String()
			}
		}
	}
}
