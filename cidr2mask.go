package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	ip := "192.168.141.111"
	add := net.ParseIP(ip)
	mask := add.DefaultMask()
	m := cidr2mask(mask.String())
	fmt.Printf("mask %s = %s  \n", mask.String(), m)
	//mask ffffff00 = 255.255.255.0
}

func cidr2mask(hex_str string) string {

	each := make([]byte, 0)
	mask := make([]string, 0)
	for i := 0; i < len(hex_str); i++ {
		if i%2 != 1 || i == 0 {
			each = append(each, hex_str[i])
		} else if i != 0 {
			each = append(each, hex_str[i])
			r, _ := strconv.ParseInt(string(each[0])+string(each[1]), 16, 10)
			m := strconv.FormatInt(r, 10)
			mask = append(mask, m)
			if len(each) == 2 {
				each = make([]byte, 0)
			}
		}

	}
	return strings.Join(mask, ".")
}
