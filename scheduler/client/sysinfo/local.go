package sysinfo

import (
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	//	"os"
	"time"
)

const (
	ProcStat = "/proc/stat"
)

type LocalSysinfo struct {
}

func init() {
	local := &LocalSysinfo{}
	RegisterSysinfoClient("local", local)
}

func getAllTimeAndIdle() (alltime uint64, idletime uint64, err error) {
	s, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		return
	}

	alltime = s.CPUStatAll.User + s.CPUStatAll.Nice +
		s.CPUStatAll.System + s.CPUStatAll.Idle +
		s.CPUStatAll.IOWait + s.CPUStatAll.IRQ +
		s.CPUStatAll.SoftIRQ + s.CPUStatAll.Steal

	idletime = s.CPUStatAll.Idle + s.CPUStatAll.IOWait
	return
}

func (c *LocalSysinfo) GetCpuUsage() (int, error) {
	alltime1, idle1, err := getAllTimeAndIdle()
	if err != nil {
		return 0, err
	}
	time.Sleep(1 * time.Second)
	alltime2, idle2, err := getAllTimeAndIdle()
	if err != nil {
		return 0, err
	}
	usedPercent := float64((alltime2-idle2)-(alltime1-idle1)) / float64(alltime2-alltime1) * 100
	return int(usedPercent), nil

}

func (c *LocalSysinfo) GetRamStat() (Total uint64, Available uint64, err error) {
	mem, err := linuxproc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		return
	}

	//获取的是KB单位
	//Total = mem.MemTotal / 1024 / 1024
	//Available = mem.MemAvailable / 1024 / 1024
	Total = mem.MemTotal / 1024
	Available = mem.MemAvailable / 1024
	fmt.Println("all:%v,free:%v", Total, Available)
	return
}

//Mb为单位
func (c *LocalSysinfo) GetDiskStat() (All uint64, Free uint64, err error) {
	disk, err := linuxproc.ReadDisk("/")
	if err != nil {
		return
	}
	//获取的是字节单位
	All = disk.All / 1024 / 1024
	Free = disk.Free / 1024 / 1024
	fmt.Println("all:%v,free:%v", All, Free)
	return

}
