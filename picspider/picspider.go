package main

import (
	//	"fmt"
	"homehtmls"
	"imgen"
	"os"
	"time"
)

func main() {
	//	this tool is used for dl images from html pages
	websit := os.Args[1]
	//	fmt.Println(websit)
	pages := homehtmls.Homehtmls(string(websit))
	for _, html := range pages {
		imgen.Imgen(html)
	}
	time.Sleep(60)
}
