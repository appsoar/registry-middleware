package sysinfo

import (
	//	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	//	"os"
	//	"errors"
	"scheduler/log"
	"time"
)

const (
	ProcStat = "/proc/stat"
)

var (
	filter = [...]string{"lo"}
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
	log.Logger.Debug("all:%vMB,free:%vMB", Total, Available)
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
	log.Logger.Debug("all:%vMB,free:%vMB", All, Free)
	return

}

func (c *LocalSysinfo) GetNetStat() ([]NetStat, error) {
	var netStat []NetStat
	//必须 通过make来进行初始化
	var netStatMap map[string]linuxproc.NetworkStat
	netStatMap = make(map[string]linuxproc.NetworkStat)

	networkStat1, err := linuxproc.ReadNetworkStat("/proc/net/dev")
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(networkStat1); i++ {
		netStatMap[networkStat1[i].Iface] = networkStat1[i]
	}
	//log.Logger.Debug(netStatMap)

	/*过滤掉lo等Iface*/
	for j := 0; j < len(filter); j++ {
		_, ok := netStatMap[filter[j]]
		if ok {
			delete(netStatMap, filter[j])
		}
	}

	time.Sleep(1 * time.Second)

	networkStat2, err := linuxproc.ReadNetworkStat("/proc/net/dev")
	if err != nil {
		panic(err)
	}

	for j := 0; j < len(networkStat2); j++ {
		iface := networkStat2[j].Iface
		//注,不能直接设置map中结构体元素中的元素
		//参考https://github.com/golang/go/issues/3117
		if v, ok := netStatMap[iface]; ok {
			v.RxBytes = networkStat2[j].RxBytes - netStatMap[iface].RxBytes
			v.TxBytes = networkStat2[j].TxBytes - netStatMap[iface].TxBytes
			netStatMap[iface] = v
		}
	}

	for _, v := range netStatMap {
		netStat = append(netStat, NetStat{Iface: v.Iface, RxBytes: v.RxBytes, TxBytes: v.TxBytes})
	}

	return netStat, nil

}
