package main

import (
	"fmt"

	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func main() {
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %v MB, Used:%v MB, Percent:%2.2f%%\n", v.Total/1024/1024, v.Used/1024/1024, v.UsedPercent)
	fmt.Println()
	c, _ := cpu.Percent(time.Duration(1)*time.Second, false)

	fmt.Printf("CPU: %.2f%%\n", c[0])
	fmt.Println()
	n0, _ := net.IOCounters(true)
	time.Sleep(time.Duration(1) * time.Second)
	n1, _ := net.IOCounters(true)
	for i, x := range n0 {
		if x.BytesRecv != 0 && x.BytesSent != 0 {
			rec := n1[i].BytesRecv - x.BytesRecv
			sent := n1[i].BytesSent - x.BytesSent
			name := x.Name
			if sent != 0 && rec != 0 {
				fmt.Printf("NIC: %v\n R:%vKbps\n T:%vKbps\n", name, sent*8/1024, rec*8/1024)
			}

		}

	}
}
