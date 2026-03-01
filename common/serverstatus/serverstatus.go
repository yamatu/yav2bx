// Package serverstatus generate the server system status
package serverstatus

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemStatus struct {
	CPU       float64
	MemTotal  uint64
	MemUsed   uint64
	SwapTotal uint64
	SwapUsed  uint64
	DiskTotal uint64
	DiskUsed  uint64
	Uptime    uint64
}

// GetSystemInfo get the system info of a given periodic
func GetSystemInfo() (Cpu float64, Mem float64, Disk float64, Uptime uint64, err error) {
	status, err := GetSystemStatus()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	Cpu = status.CPU
	Uptime = status.Uptime
	if status.MemTotal > 0 {
		Mem = float64(status.MemUsed) / float64(status.MemTotal) * 100
	}
	if status.DiskTotal > 0 {
		Disk = float64(status.DiskUsed) / float64(status.DiskTotal) * 100
	}

	return Cpu, Mem, Disk, Uptime, nil
}

func GetSystemStatus() (st *SystemStatus, err error) {
	st = &SystemStatus{}

	errorString := ""

	cpuPercent, err := cpu.Percent(0, false)
	// Check if cpuPercent is empty
	if len(cpuPercent) > 0 && err == nil {
		st.CPU = cpuPercent[0]
	} else {
		st.CPU = 0
		errorString += fmt.Sprintf("get cpu usage failed: %s ", err)
	}

	memUsage, err := mem.VirtualMemory()
	if err != nil {
		errorString += fmt.Sprintf("get mem usage failed: %s ", err)
	} else {
		st.MemTotal = memUsage.Total
		st.MemUsed = memUsage.Used
	}

	swapUsage, err := mem.SwapMemory()
	if err != nil {
		errorString += fmt.Sprintf("get swap usage failed: %s ", err)
	} else {
		st.SwapTotal = swapUsage.Total
		st.SwapUsed = swapUsage.Used
	}

	diskUsage, err := disk.Usage("/")
	if err != nil {
		errorString += fmt.Sprintf("get disk usage failed: %s ", err)
	} else {
		st.DiskTotal = diskUsage.Total
		st.DiskUsed = diskUsage.Used
	}

	uptime, err := host.Uptime()
	if err != nil {
		errorString += fmt.Sprintf("get uptime failed: %s ", err)
	} else {
		st.Uptime = uptime
	}

	if errorString != "" {
		err = fmt.Errorf("%s", errorString)
	}

	return st, err
}
