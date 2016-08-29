package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	//	"strconv"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			//				fmt.Printf("\n\n%v", err)
		}
	}()
	urh := regexp.MustCompile("(https?://)?(\\S+)/?")
	url := os.Args[1]
	parse := urh.FindStringSubmatch(url)
	proto := parse[1]
	host := parse[2]
	if len(proto) == 0 {
		url = "http://" + url
	}
	ip, _ := net.LookupIP(host)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Proto = "HTTP/1.1"

	req.Header.Add("User-Agent",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36")
	req.Header.Add("Accept",
		"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	//	req.Header.Add("Accept-Encoding",
	//		"gzip, deflate, sdch, br")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	outputs, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s %s\n", req.Method, req.URL, req.Proto)
	fmt.Printf("Host: %s\n", host)
	fmt.Printf("Address: %v\n", ip)
	for k, v := range req.Header {
		fmt.Println(string(k)+":", v[0])
	}
	fmt.Printf("\n")
	fmt.Printf("%s %s\n", res.Proto, res.Status)
	for k, v := range res.Header {
		fmt.Println(string(k)+":", v[0])
	}

	if os.Args[2] == "-p" || os.Args[2] == "--print" {
		fmt.Println(string(outputs))
	}
}
