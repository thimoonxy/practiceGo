package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"os"
	"practiceGo/logana/src/logline"
	"practiceGo/logana/src/output"
	"reflect"
	"sort"
	"time"
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

func ArgsParse() (int, string, int64, int64, bool) {
	//default timestamps
	now := time.Now()
	default_end_str := now.Format("2006 Jan 02 15:04:05 -0700")
	default_start := now.Add(time.Duration(-60) * time.Second)
	default_start_str := default_start.Format("2006 Jan 02 15:04:05 -0700")
	//flags
	depth := flag.Int("depth", 1, "Depth of url or reference. Show whole URL when <=0 .")
	fieldname := flag.String("fieldname", "Hostname", "Avail Fieldnames: Remote_host, Log_name, Auth_user, URL, Protocol, Reference, User_agent, Method, Raw_data, Scheme, Hostname, ReferenceHostname, ReferenceScheme, Time_stamp, Status, Response_size ")
	starttime_str := flag.String("start_time", default_start_str, "query from this start_time")
	endtime_str := flag.String("end_time", default_end_str, "query till this end_time")
	all := flag.Bool("all", false, "Show all outputs.")
	flag.Parse()
	//parse time flags
	s_time, _ := time.Parse("2006 Jan 02 15:04:05 -0700", *starttime_str)
	e_time, _ := time.Parse("2006 Jan 02 15:04:05 -0700", *endtime_str)
	start_time, end_time := s_time.Unix(), e_time.Unix()
	return *depth, *fieldname, start_time, end_time, *all
}

func main() {
	r, _ := os.Open("c:\\Users\\simon\\Downloads\\varnishncsa.log.5.gz")
	defer r.Close()
	g, _ := gzip.NewReader(r)
	buf := bufio.NewReader(g)
	//	l := &io.LimitedReader{R: buf, N: 20}
	//	io.Copy(ioutil.Discard, l.R)

	depth, fieldname, start_time, end_time, all := ArgsParse()
	counter, stat, sum := FieldCounter(buf, start_time, end_time, fieldname, depth)
	//		./logana.exe -start_time "2016 Jan 12 06:52:02 -0500" -end_time "2016 Jan 12 06:54:12 -0500" -fieldname  Status

	slice := new(output.Output_slice)
	slice.Output_slice_gen(counter, stat, sum)
	fmt.Printf("Percent        Count          Field - %s  \n", fieldname)

	fmt.Println("_______________________________________________")
	sort.Sort(output.Records(slice.Records))
	Records := slice.Records
	if len(Records) <= 20 || all == true {
		for _, record := range Records {
			fmt.Printf("%6s        %6.0f          %v  \n", record.Fmt_Percent, record.Number, record.Name)
		}
	} else if len(Records) > 20 {
		for i := 0; i < 10; i++ {
			fmt.Printf("%6s        %6.0f          %v  \n", Records[i].Fmt_Percent, Records[i].Number, Records[i].Name)
		}
		fmt.Println("...... Records more than 20 ......")
		fmt.Println("...... Truncated .................")
		for i := len(Records) - 10; i < len(Records); i++ {
			fmt.Printf("%6s        %6.0f          %v  \n", Records[i].Fmt_Percent, Records[i].Number, Records[i].Name)
		}
		fmt.Printf("\nUse --all parm to show all records queried.\n\n")
	}
	fmt.Println("_______________________________________________")
	fmt.Printf("Totally %.0f log lines during the query time period.\n", slice.Sum)
}
