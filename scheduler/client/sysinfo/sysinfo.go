package sysinfo

import (
	"fmt"
)

var (
	sysInfoClients map[string]SysInfoClient
)

type SysInfoClient interface {
	/*里面提供的方法再定*/
	GetCpuUsage() (int, error)

	GetRamStat() (uint64, uint64, error)
	GetDiskStat() (uint64, uint64, error)
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
