package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/mitchellh/ioprogress"
)

var (
	src, dst *string
	err      error
)
var verbose = new(bool)
var human = new(bool)
var progress = new(bool)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	start := time.Now().Unix()
	sep := sep_perOS()
	pwd, _ := os.Getwd()
	src = flag.String("s", pwd, "src folder path")
	dst = flag.String("d", pwd, "dst folder path")
	//	*verbose = false
	flag.BoolVar(verbose, "v", false, "Print Verbose.")
	flag.BoolVar(human, "h", false, "Print human readable output .")
	flag.BoolVar(progress, "p", false, "Print progress status .")
	*src = strings.TrimRight(*src, sep)
	*dst = strings.TrimRight(*dst, sep)
	// seems golang only recognize '/'
	*src = strings.Replace(*src, "\\", "/", -1)
	*dst = strings.Replace(*dst, "\\", "/", -1)
	flag.Parse()
	fmt.Printf("src: %s\ndst: %s\nCopying...\n\n", *src, *dst)

	lsrc := flag.Lookup("s")
	if lsrc.Value.String() == lsrc.DefValue {
		flag.Usage()
		os.Exit(1)
	}
	ldst := flag.Lookup("d")
	if ldst.Value.String() == ldst.DefValue {
		flag.Usage()
		os.Exit(1)
	}

	walk(*src)
	//	time.Sleep(200 * time.Second)
	fmt.Printf("\n\nElapsed: %ds .", time.Now().Unix()-start)
}

func walk(path string) {
	i, _ := ioutil.ReadDir(path)
	sep := sep_perOS()
	for _, info := range i {
		perm := info.Mode()
		if info.IsDir() {
			next_src_path := path + sep + info.Name()
			next_dst_path := *dst + strings.Split(next_src_path, *src)[1]
			//			fmt.Println("mkdir ", next_dst_path)
			os.MkdirAll(next_dst_path, perm)
			walk(next_src_path)
		} else {
			src_fname := path + sep + info.Name()
			dst_fname := *dst + strings.Split(src_fname, *src)[1]
			dst_folder := *dst + strings.Split(path, *src)[1]
			os.MkdirAll(dst_folder, perm)
			//			fmt.Println("dst_folder ", dst_folder)
			cp(src_fname, dst_fname)
		}
	}

}

var byteUnits = []string{"B", "KB", "MB", "GB", "TB", "PB"}

func DrawTerminal(w io.Writer) ioprogress.DrawFunc {
	return ioprogress.DrawTerminalf(w, func(progress, total int64) string {
		return fmt.Sprintf("%s/%s", byteUnitStr(progress), byteUnitStr(total))
	})
}

func byteUnitStr(n int64) string {
	var unit string
	size := float64(n)
	for i := 1; i < len(byteUnits); i++ {
		if size < 1024 {
			unit = byteUnits[i-1]
			break
		}

		size = size / 1024
	}

	return fmt.Sprintf("%.2f %s", size, unit)
}

func cp(src_fname, dst_fname string) {
	dst, _ := os.Create(dst_fname)
	src, _ := os.Open(src_fname)
	src_stat, _ := src.Stat()
	st := os.Stdout
	if *verbose == true {
		fmt.Fprintf(st, "%s\n", dst_fname)
	}
	var draw ioprogress.DrawFunc
	if *human == true {
		draw = DrawTerminal(st)
	} else {
		draw = nil
	}

	progressR := &ioprogress.Reader{
		Reader:   src,
		Size:     src_stat.Size(),
		DrawFunc: draw,
	}
	if *progress == true {
		_, err = io.Copy(dst, progressR)
	} else {
		_, err = io.Copy(dst, src)
	}

	if err != nil {
		fmt.Fprintf(st, "# Failed copying %s to %s .\n", src_fname, dst_fname)
	}
}
func sep_perOS() (sep string) {
	os, has := os.LookupEnv("OS")
	if has {
		os = strings.ToLower(os)
		windows := regexp.MustCompile(".*(windows).*")
		if len(windows.FindStringSubmatch(os)) >= 1 {
			sep = "\\"
		} else {
			sep = "/"
		}

	}
	return
}
