package reader

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func Locate(path string) FileSlice {
	// Testing env
	os, has := os.LookupEnv("OS")
	if has {
		os = strings.ToLower(os)
		windows := regexp.MustCompile(".*(windows).*")
		if len(windows.FindStringSubmatch(os)) >= 1 {
			path = "c:\\Users\\simon\\Downloads\\"
		}

	}
	//
	varnishlog := regexp.MustCompile(".*varnish.*\\.log.\\d*$")
	varnishgz := regexp.MustCompile(".*varnish.*\\.gz$")
	files, _ := ioutil.ReadDir(path)
	fileslice := make(FileSlice, 0)
	fi := FI{}
	for _, v := range files {
		if varnishlog.MatchString(v.Name()) || varnishgz.MatchString(v.Name()) {
			fi.ModTime = v.ModTime()
			fi.Name = path + v.Name()
			fi.Size = v.Size()
			fi.Type = "plain"
			if varnishgz.MatchString(v.Name()) {
				fi.Type = "gz"
			}
			fileslice = append(fileslice, fi)
		}
	}
	sort.Sort(fileslice)
	//	fmt.Println(fileslice)
	return fileslice
}

func (f FileSlice) Read(start_time, end_time int64) ([]*os.File, []*bufio.Reader) {
	var first_index, last_index int
	//	fmt.Println("len(fileslice):", len(f))
	/*
		for i, fi := range f {
			if len(f) == 1 {
				break
			}
			if len(f)-1 >= i+1 && start_time > f[i+1].ModTime.Unix() {
				last_index = i
			} else if len(f)-1 >= i+1 && start_time <= fi.ModTime.Unix() && start_time > f[i+1].ModTime.Unix() {
				last_index = i
			}
			if len(f)-1 >= i+1 && end_time <= fi.ModTime.Unix() && end_time > f[i+1].ModTime.Unix() {
				first_index = i
			}
		}
	*/
	timelines := make([]TimeLine, 5)
	// initializing
	for i, fi := range f {
		tl := TimeLine{}
		tl.Fname = fi.Name
		tl.LastLineTime = fi.ModTime.Unix()
		tl.Index = i
		timelines = append(timelines, tl)
	}
	// 1st line time = next file's lastline
	var ZeroTime time.Time
	for i, _ := range timelines {
		if i+1 == len(f) {
			timelines[i].FirstLineTime = ZeroTime.Unix()
			break
		} else {
			timelines[i].FirstLineTime = timelines[i+1].LastLineTime
		}
	}

	//judge if read this file
	for _, tl := range timelines {
		if tl.FirstLineTime < start_time && start_time <= tl.LastLineTime {
			last_index = tl.Index
		}
		if tl.FirstLineTime < end_time && end_time <= tl.LastLineTime {
			first_index = tl.Index
		}
	}

	iolist := make([]*os.File, 0)
	buflist := make([]*bufio.Reader, 0)
	//	fmt.Println("first_index,last_index:", first_index, last_index)
	namelist := make([]string, 0)

	for i := first_index; i <= last_index; i++ {
		fname := f[i].Name
		r, _ := os.Open(fname)
		namelist = append(namelist, fname)
		iolist = append(iolist, r)
		var buf *bufio.Reader
		if f[i].Type == "gz" {
			g, _ := gzip.NewReader(r)
			buf = bufio.NewReader(g)
		} else {
			buf = bufio.NewReader(r)
		}
		buflist = append(buflist, buf)
	}
	//	fmt.Println("iolist, buflist:", iolist, buflist)
	fmt.Println("Reading Files:", strings.Join(namelist, ", "))
	return iolist, buflist

}

type FI struct {
	ModTime time.Time
	Name    string
	Size    int64
	Type    string
}

type FileSlice []FI

func (s FileSlice) Less(i, j int) bool {
	return s[j].ModTime.Unix() < s[i].ModTime.Unix()
}

func (a FileSlice) Len() int {
	return len(a)
}
func (a FileSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type TimeLine struct {
	Fname                       string
	FirstLineTime, LastLineTime int64
	Index                       int
}
