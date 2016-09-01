package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
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
	scheme := parse[1]
	host := parse[2]
	if len(scheme) == 0 {
		url = "http://" + url
	}
	ip, _ := net.LookupIP(host)

	tr := &http.Transport{
		DisableCompression: true,
	}

	client := &http.Client{Transport:tr}
	req, err := http.NewRequest("GET", url, nil)
	req.Proto = "HTTP/1.1"

	req.Header.Add("User-Agent",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36")
	req.Header.Del("Accept-Encoding")
	//	req.Header.Add("Accept",
	//		"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	//	req.Header.Add("Accept-Encoding",
	//	"gzip, deflate, sdch, br")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Resoved: ", ip)
	dump, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(dump))

	dump, _ = httputil.DumpResponse(res, false)
	fmt.Println(string(dump))

	//	fmt.Println(req.ContentLength)
	outputs, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if os.Args[2] == "-p" || os.Args[2] == "--print" {
		fmt.Println(string(outputs))
	}
}
