package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"practiceGo/logana/src/logline"
	"practiceGo/logana/src/output"
	"reflect"
)

func FieldCounter(buf *bufio.Reader, starttime int64, endtime int64, fieldname string, depth int) (map[interface{}]float64, map[interface{}]float64, float64) {
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
	stat := map[interface{}]float64{}
	var sum float64
	for _, v := range counter {
		sum += v
	}
	for k, _ := range counter {
		percent := (counter[k]) / sum * 100
		//		stat[k] = strconv.FormatFloat(percent, 'f', 2, 64) + "%"
		stat[k] = percent
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

	depth := 2
	fieldname := "Hostname"
	counter, stat, sum := FieldCounter(buf, 1452599529, 1452599530, fieldname, depth)
	//	fmt.Println(counter)
	//	fmt.Println(stat)
	//	fmt.Println(sum)

	slice := new(output.Output_slice)
	slice.Output_slice_gen(counter, stat, sum)
	fmt.Printf("Percent        Count          Field - %s  \n", fieldname)

	fmt.Println("_______________________________________________")
	for _, record := range slice.Records {
		fmt.Printf("%6s        %6.0f          %s  \n", record.Fmt_Percent, record.Number, record.Name)
	}
	fmt.Println("_______________________________________________")
	fmt.Printf("Totally %.0f lines during the qurey time period.\n", slice.Sum)
}
