package sysinfo

import ()

type LocalSysinfo struct {
}

func init() {
	local := &LocalSysinfo{}

	RegisterSysinfoClient("local", local)

}

func (c *LocalSysinfo) GetCpuUsage() int {
	//

}

func (c *LocalSysinfo) GetRamUsage() {
	//

}

func (c *LocalSysinfo) GetDiskUsage() {
	//

}
