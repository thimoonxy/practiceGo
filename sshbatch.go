package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	global_timeout := 30 * time.Second
	std := StdJudge(os.Stdin, global_timeout)
	hostlist := HostListGen(std, 22)
	PassWd := []ssh.AuthMethod{ssh.Password("0225")}
	Conf := &ssh.ClientConfig{User: "root", Auth: PassWd, Timeout: 5 * time.Second}
	onetime := 6
	RunCtl(hostlist, onetime, Conf)

}

func RunCtl(hostlist []string, onetime int, conf *ssh.ClientConfig) {
	sum := len(hostlist)
	if sum <= onetime {
		onetime = sum
	}
	//fmt.Println(onetime)
	wg := &sync.WaitGroup{}
	stdout := make(chan string, sum)
	for i := 0; i < sum/onetime; i++ {
		wg.Add(onetime)
		for j := i * onetime; j < onetime*(i+1); j++ {
			go func(idx int) {
				stdout <- SSH(conf, hostlist[idx])
				wg.Done()
			}(j)
		}
		wg.Wait()
	}

	// remainder runs
	if sum%onetime != 0 {
		wg.Add(onetime)
		for x := sum / onetime * onetime; x < sum; x++ {
			go func(idx int) {
				stdout <- SSH(conf, hostlist[idx])
				wg.Done()
			}(x)
		}
	}
	// printout
	close(stdout)
	for output := range stdout {
		fmt.Println(output)
	}
}

func StdJudge(s io.Reader, timeout time.Duration) []byte {
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
			case <-time.After(timeout):
				fmt.Println("timeout")
				os.Exit(1)
			}
		}
	}()
	stdbuf, err := ioutil.ReadAll(s)
	if err != nil {
		os.Exit(1)
	}

	stdinreader <- stdbuf
	return stdbuf
}

func HostListGen(raw []byte, default_port int) []string {
	var sep []byte
	sep = append(sep, '\n')
	tmp := bytes.SplitN(raw, sep, len(raw))
	//port := ":" + strconv.Itoa(default_port)

	var result []string
	for _, v := range tmp {
		strvar := string(v)
		if !strings.Contains(strvar, ":") && len(strvar) > 0 {
			strvar = net.JoinHostPort(strvar, strconv.Itoa(default_port))
			result = append(result, strvar)
		} else if len(strvar) > 0 {
			result = append(result, strvar)
		}
	}
	return result
}

func SSH(Conf *ssh.ClientConfig, ip_port string) string {

	var r *bufio.Reader
	//PassWd := []ssh.AuthMethod{ssh.Password(password)}
	//Conf := ssh.ClientConfig{User: user, Auth: PassWd, Timeout: 5 * time.Second}
	Client, err := ssh.Dial("tcp", ip_port, Conf)
	if err != nil {
		return ip_port + ": " + err.Error()
	}
	defer Client.Close()

	if len(os.Args) >= 2 {
		s := strings.NewReader(os.Args[1])
		r = bufio.NewReader(s)
	} else {
		fmt.Println("No cmd in Args.")
		os.Exit(1)
	}

	for {
		b, _, e := r.ReadLine()
		if e == io.EOF {
			return ""
		}

		if e != nil {
			fmt.Println("no stdin or args.")
			return ip_port + ": " + e.Error()
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
				return stdoutBuf.String() + stderrBuf.String()
			} else {
				return stdoutBuf.String()
			}
		}
	}
}
