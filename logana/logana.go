package main

import (
	"fmt"
	//	"runtime"
	//	"sync"
	"bufio"
	"compress/gzip"
	//	"io"
	//	"io/ioutil"
	"os"
	"practiceGo/logana/src/logline"
	"reflect"
	"strconv"
)

func FieldCounter(buf *bufio.Reader, starttime int64, endtime int64, fieldname string, depth int) (interface{}, interface{}, interface{}) {
	counter := map[interface{}]float64{}

	r := reflect.ValueOf(counter)
	k := r.MapKeys()
	t := r.Interface()
	//counter section
	log := new(logline.Varnishlog_line)

	if depth <= 0 {
		depth = -1
	}

	for {
		line, err := buf.ReadBytes('\n')
		if err == nil {
			log.LogFieldGen(string(line))
			log.URLdepth(depth)
			log.ReferenceDepth(depth)
			if log.Time_stamp.Unix() >= starttime && endtime >= log.Time_stamp.Unix() {
				fieldvalue := log.FieldValue(fieldname)

				haskey := false
				for _, v := range k {
					if fieldvalue == v.Interface() {
						counter[v.Interface()] += 1
						haskey = true
					}
					fmt.Println(v.Interface(), t, counter[v.Interface()])

				}
				if haskey == false {
					counter[fieldvalue] += 1
				}
			}

		} else {
			break
		}
	}
	// stat section
	stat := map[interface{}]string{}
	var sum float64
	for _, v := range counter {
		sum += v
	}
	for k, _ := range counter {
		percent := (counter[k]) / sum * 100
		stat[k] = strconv.FormatFloat(percent, 'f', 2, 64) + "%"
	}

	return counter, stat, sum
}

func main() {
	r, _ := os.Open("c:\\Users\\simon\\Downloads\\varnishncsa.log.5.gz")
	defer r.Close()
	g, _ := gzip.NewReader(r)
	buf := bufio.NewReader(g)
	//	l := &io.LimitedReader{R: buf, N: 20}
	//	io.Copy(ioutil.Discard, l.R)
	counter, stat, sum := FieldCounter(buf, 1452599529, 1452599530, "Hostname", 0)
	fmt.Println(counter)
	fmt.Println(stat)
	fmt.Println(sum)
}
