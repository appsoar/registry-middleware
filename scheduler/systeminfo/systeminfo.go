package systeminfo

import ()

type RamUsage struct {
	Used  int64
	Total int64
}

type StorageUsage struct {
	Used  int64
	Total int64
}

type CpuUsage struct {
	Used  int64
	Total int64
}

type SystemInfo struct {
	//	RamUsage
	//	StorageUsage
	//	CpuUsage
	TimeUpdate int //设置更新时间:second?
}

func NewSystemInfo(TimeUpdate int) error {
	//	return &
	return &SystemInfo{TimeUpdate: TimeUpdate}, nil
}

func (systeminfo SystemInfo) GetRamUsage() (RamUsage, error) {
}

func (systeminfo SystemInfo) GetStorageUsage() (StorageUsage, error) {

}

type SystmeInfo interface {
}
