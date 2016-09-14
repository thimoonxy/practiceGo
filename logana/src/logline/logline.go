package logline

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Varnishlog_line struct {
	/*
		http://www.varnish-cache.org/docs/trunk/reference/varnishncsa.html

		Specify the log format to use. If no format is specified the default log format is used:
		%h %l %u %t "%r" %s %b "%{Referer}i" "%{User-agent}i"

		e.g.
		1.2.6.7 - - [12/Jan/2016:06:52:02 -0500] "GET http://www.simon.com.cn:80/apps/570/icons/econ/items/nyx_back/e8aae23e506057.png HTTP/1.0" 200 57080 "http://api.simon.com/api/econ_items/econ_items_display/?steam_id=149232154&pkey=MTQ4NC41ODEzNzFqdmtxdnNwc3R2aWFo&phone_num=13717609020ndroid&os_version=5.0&version=3.1.0" "Mozilla/5.0 (Linux; Android 5.0; Lenovo K50-t5 Build/LRX21M) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/37.0.0.0 Mobile Safari/537.36"

		%h
			Remote host. Defaults to '-' if not known. In backend mode this is the IP of the backend server.

		%l
			Remote logname. Always '-'.

		%u
			Remote user from auth.

		%t
			In client mode, time when the request was received, in HTTP date/time format. In backend mode, time when the request was sent.

		%r
		The first line of the request. Synthesized from other fields, so it may not be the request verbatim.

		%s
			Status sent to the client. In backend mode, status received from the backend.

		%b
			In client mode, size of response in bytes, excluding HTTP headers. In backend mode, the number of bytes received from the backend, excluding HTTP headers. In CLF format, i.e. a '-' rather than a 0 when no bytes are sent.

		%reference

		%user-agent
	*/
	Remote_host, Log_name, Auth_user, URL, Protocol, Reference, User_agent, Method, raw_data string
	Time_stamp                                                                               time.Time
	Status                                                                                   int
	Response_size                                                                            int64
}

func (log *Varnishlog_line) logfield(line string) {
	log.raw_data = line

	// Common fields
	partten := regexp.MustCompile("([\\S]+)")
	match := partten.FindAllStringSubmatch(line, 20)
	if len(match) >= 1 {
		log.Remote_host = match[0][1]
		log.Log_name = match[1][1]
		log.Auth_user = match[2][1]
		log.Method = strings.Trim(match[5][1], "\"")
		log.URL = match[6][1]
		log.Protocol = strings.Trim(match[7][1], "\"")
		log.Status, _ = strconv.Atoi(match[8][1])
		log.Response_size, _ = strconv.ParseInt(match[9][1], 10, 64)
	}
	//reference
	reference := regexp.MustCompile("\"([\\S]+)\"")
	match = reference.FindAllStringSubmatch(line, 9)
	if len(match) >= 1 {
		log.Reference = match[0][1]
	}
	//user-agent
	ua := regexp.MustCompile("\" \"(\\S.*)\"\n$")
	if len(match) >= 1 {
		match = ua.FindAllStringSubmatch(line, 9)
		log.User_agent = match[0][1]
	}
	//timestamp
	stamp := regexp.MustCompile(" \\[(\\d+)/(\\w+)/(\\d+):([\\d+:]+\\d+)(.*)\\] \"")
	//	log.Time_stamp.Format("02/Jan/2006:15:04:05 -0700")
	match = stamp.FindAllStringSubmatch(line, 20)
	if len(match) >= 1 {
		date := match[0][1]
		month := match[0][2]
		year := match[0][3]
		time_value_string := year + " " + month + " " + date + " " + match[0][4] + match[0][5]
		log.Time_stamp, _ = time.Parse("2006 Jan 02 15:04:05 -0700", time_value_string)
	}

}
