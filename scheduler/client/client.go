package client

import (
	"fmt"
	"scheduler/client/database"
	"scheduler/client/registry"
	"scheduler/client/sysinfo"
)

/*
type ClientOpts struct {
	Url       string
	AccessKey string
	SecretKey string
}

type BaseClient struct {
	Opts *ClientOpts
}
*/

//所有访问的总接口
//
type Client struct {
	database database.DatabaseClient
	registry registry.RegistryClient
	sysInfo  sysinfo.SysInfoClient
	//  lock - 数据同步时使用?还是在subClient中实现?
}

func constructClient() *Client {
	var err error
	client := new(Client)

	client.database, err = database.GetDatabaseClient()
	if err != nil {
		//		panic(err)
	}

	client.sysInfo, err = sysinfo.GetSysInfoClient()
	if err != nil {
		//		panic(err)
	}

	client.registry, err = registry.GetRegistryClient()
	if err != nil {
		panic(err)
	}

	//client.Registry = newRegistryClient(client)

	return client
}

func NewClient() (*Client, error) {
	client := constructClient()
	//	client.BaseClient.Opts = opts
	return client, nil
}

type SysInfo struct {
	CpuUsage      int    `json:"cpuUsage"`
	TotalRam      uint64 `json:"totalRam"`
	AvailableRam  uint64 `json:"availableRam"`
	TotalDisk     uint64 `json:"totalDisk"`
	AvailableDisk uint64 `json:"availableDisk"`
}

func (c *Client) GetSysInfo() (SysInfo, error) {
	//	cpuUsage,err := c.sysInfo.
	//c.sysInfo.GetCpuUsage()
	sysinfo := new(SysInfo)
	var cpuchan = make(chan int)
	var ramchan = []chan uint64{
		make(chan uint64),
		make(chan uint64),
	}
	var diskchan = []chan uint64{
		make(chan uint64),
		make(chan uint64),
	}

	go func() {
		cpuUsage, err := c.sysInfo.GetCpuUsage()
		if err != nil {
			panic(err)
		}
		cpuchan <- cpuUsage
	}()

	go func() {
		totalRam, availableRam, err := c.sysInfo.GetRamStat()
		if err != nil {
			panic(err)
		}
		ramchan[0] <- totalRam
		ramchan[1] <- availableRam
	}()

	go func() {
		totalDisk, availableDisk, err := c.sysInfo.GetDiskStat()
		if err != nil {
			panic(err)
		}
		diskchan[0] <- totalDisk
		diskchan[1] <- availableDisk
	}()

	sysinfo.CpuUsage = <-cpuchan
	sysinfo.TotalRam = <-ramchan[0]
	sysinfo.AvailableRam = <-ramchan[1]
	sysinfo.TotalDisk = <-diskchan[0]
	sysinfo.AvailableDisk = <-diskchan[1]

	fmt.Printf("cpu:%v, total:%v,avail:%v, total:%v,avail:%v\n", sysinfo.CpuUsage, sysinfo.TotalRam, sysinfo.AvailableRam, sysinfo.TotalRam, sysinfo.AvailableDisk)

	return *sysinfo, nil

}

func (c *Client) ListImages() (interface{}, error) {
	images, err := c.registry.ListImages()
	return images, err
}

func (c *Client) GetImageTags(image string) (interface{}, error) {

	tags, err := c.registry.GetImageTags(image)
	return tags, err
}

func (c *Client) GetImageDigest(image string, tag string) (interface{}, error) {

	digest, err := c.registry.GetImageDigest(image, tag)
	return digest, err
}

func (c *Client) DeleteImageDigest(image string, tag string) error {

	err := c.registry.DeleteImageTag(image, tag)
	return err
}
