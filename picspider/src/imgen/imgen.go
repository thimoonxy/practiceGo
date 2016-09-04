package imgen

// imgen pkg used for dl img from a html page
import (
	"fmt"
	//	"fmt"
	"homehtmls"

	"regexp"
)

func Imgen(url string) {
	defer func() {
		if err := recover(); err != nil {
		}
	}()

	out := Curl(url)
	urh := Htmlurh(url)
	Pages(out, urh)
}

func Curl(url string) []byte {
	img := regexp.MustCompile("<img src=\"(http://[\\S]*.jpg)")

	out := homehtmls.Curl(url, true)

	for _, v := range img.FindAllSubmatch(out, 2000) {
		//		fmt.Println(string(v[1]))
		go homehtmls.Curl(string(v[1]), false)
		//		homehtmls.Curl(string(v[1]), false)
		fmt.Println("<img src=\"" + string(v[1]) + "\">")
	}
	return out
}

func Pages(out []byte, urh string) {
	html := regexp.MustCompile("<a href=\"?'?(http://[\\S]+)*/([\\S]+.html)")
	for _, v := range html.FindAllSubmatch(out, 2000) {
		fmt.Println("url in imgen:", urh, string(v[2]))
		Curl(urh + string(v[2]))
		fmt.Println(urh + string(v[2]))
	}
}

func Htmlurh(url string) string {

	html := regexp.MustCompile("(http://\\S+/)[\\S]+.html")
	out := html.FindStringSubmatch(url)
	urh := string(out[1])
	//	fmt.Println(urh, "\n", fname)
	return urh
}
