package sysinfo

import (
	"fmt"
)

var (
	sysInfoClients map[string]SysInfoClient
)

type NetStat struct {
	Iface   string `json:"iface"`
	TxBytes uint64 `json:"TxBytes"`
	RxBytes uint64 `json:"RxBytes"`
}

type SysInfoClient interface {
	GetCpuUsage() (int, error)

	GetRamStat() (uint64, uint64, error)
	GetDiskStat() (uint64, uint64, error)
	//	GetNetStat() ([]NetStat, error)
	GetNetIfs() (interface{}, error)
	GetNetIfStat(string) (interface{}, error)
}

func RegisterSysinfoClient(name string, client SysInfoClient) error {
	if sysInfoClients == nil {
		sysInfoClients = make(map[string]SysInfoClient)
	}

	if _, exists := sysInfoClients[name]; exists {
		return fmt.Errorf("SysinfoClient already registered")
	}

	sysInfoClients[name] = client
	return nil
}

func GetSysInfoClient() (SysInfoClient, error) {
	/*
		name := os.Getnenv("SysinfoClient")
		if name != nil {
			return nil, fmt.Errorf("sysinfoclient not support.")

	*/
	name := "local"

	if client, ok := sysInfoClients[name]; ok {
		return client, nil
	}
	return nil, fmt.Errorf("sysinfoClient[%s] not suppport", name)
}

func init() {
}
