package homehtmls

// homehtmls gets all the html urls from homepage
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	//	"imgen"
	//	"log"
)

func Homehtmls(url string) []string {
	//	defer func() {
	//		if err := recover(); err != nil {
	//		}
	//	}()

	out := Curl(url, true)
	urh := Htmlurh(url)
	//	fmt.Println("homehtml_main:", urh)
	pages := Pages(out, urh)
	return pages
}

func Curl(url string, read bool) (out []byte) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}

	}()
	//	fmt.Println(url)

	tr := &http.Transport{DisableCompression: true}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("GET", url, nil)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	req.Proto = "HTTP/1.1"

	req.Header.Add("User-Agent",
		"Mozilla/5.0  AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.106 Safari/537.36")

	res, _ := client.Do(req)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	if read {
		out, _ = ioutil.ReadAll(res.Body)
	} else {
		out = nil
	}

	defer res.Body.Close()

	return out
}

func Pages(out []byte, urh string) (pages []string) {
	html := regexp.MustCompile("href=\"?'?(http://)?(\\w*\\.?\\w+\\.\\w+)?/?(\\S+/)*(\\S+.html)")
	for _, v := range html.FindAllSubmatch(out, 2000) {
		var url string
		if len(v[2]) > 0 {
			//			fmt.Println("v1: ", string(v[1]))
			//			fmt.Println("v2: ", string(v[2]))
			//			fmt.Println("v3: ", string(v[3]))
			//			fmt.Println("v4: ", string(v[4]))
			url = string(v[0])
		} else {
			url = urh + "/" + string(v[3]) + string(v[4])
		}

		//		fmt.Println("homehtmls: ", url)
		pages = append(pages, url)
	}
	return pages
}

func Htmlurh(url string) (urh string) {

	html := regexp.MustCompile("(http://)?(\\w*\\.?\\w+\\.\\w+)?/?(\\S+/)*(\\S+.html)?")
	out := html.FindStringSubmatch(url)
	//	fmt.Println(out)
	scheme := out[1]
	host := out[2]
	urh = scheme + host
	return
}
