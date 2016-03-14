package client

import (
	//	"fmt"
	"scheduler/client/database"
	"scheduler/client/registry"
	"scheduler/client/sysinfo"
	"scheduler/log"
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
	CpuUsage      int               `json:"cpuUsage"`
	TotalRam      uint64            `json:"totalRam"`
	AvailableRam  uint64            `json:"availableRam"`
	TotalDisk     uint64            `json:"totalDisk"`
	AvailableDisk uint64            `json:"availableDisk"`
	NetStat       []sysinfo.NetStat `json:"netStat"`
}

func (c *Client) GetSysInfo() (SysInfo, error) {
	//	cpuUsage,err := c.sysInfo.
	//c.sysInfo.GetCpuUsage()
	info := new(SysInfo)
	var cpuchan = make(chan int)
	var ramchan = []chan uint64{
		make(chan uint64),
		make(chan uint64),
	}
	var diskchan = []chan uint64{
		make(chan uint64),
		make(chan uint64),
	}
	var netchan = make(chan []sysinfo.NetStat)

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

	go func() {
		netStat, err := c.sysInfo.GetNetStat()
		if err != nil {
			panic(err)
		}
		netchan <- netStat
	}()

	info.CpuUsage = <-cpuchan
	info.TotalRam = <-ramchan[0]
	info.AvailableRam = <-ramchan[1]
	info.TotalDisk = <-diskchan[0]
	info.AvailableDisk = <-diskchan[1]
	info.NetStat = <-netchan

	log.Logger.Debug("cpu:%v, total:%v,avail:%v, total:%v,avail:%v\n", info.CpuUsage, info.TotalRam, info.AvailableRam, info.TotalRam, info.AvailableDisk)
	log.Logger.Debug("%+v", info.NetStat)

	return *info, nil

}

func (c *Client) ListImages() ([]string, error) {
	images, err := c.registry.ListImages()
	if err != nil {
		return nil, err
	}
	list, ok := images.([]string)
	if !ok {
		panic("low-level registry api have changed")
	}

	return list, nil
}

func (c *Client) GetImageTags(image string) ([]string, error) {

	tags, err := c.registry.GetImageTags(image)
	if err != nil {
		return nil, err
	}
	list, ok := tags.([]string)
	if !ok {
		panic("low-level registry api have changed")
	}
	return list, nil
}

func (c *Client) GetImageDigest(image string, tag string) (string, error) {

	digest, err := c.registry.GetImageDigest(image, tag)
	if err != nil {
		return "", err
	}
	dg, ok := digest.(string)
	if !ok {
		panic("low-level registry api have changed")
	}
	return dg, nil
}

func (c *Client) DeleteImageDigest(image string, tag string) error {

	err := c.registry.DeleteImageTag(image, tag)
	return err
}

/*
type Comment struct {
	Time    string `json:"time"`
	User    string `json:"user"`
	Content string `json:"content"`
}

type ImageProperty struct {
	Name        string    `json:"name"`
	Public      bool      `json:"public"`
	Namespace   string    `json:"namespace"`
	Tags        []string  `json:"tags"`
	Download    uint      `json:"download"`
	Description string    `json:"description"`
	Comments    []Comment `json:"comments"`
}

func (c *Client) GetImageProperty(image string) ImageProperty {
	return nil
}

func (c *Client) SetImagePublic(Public bool) {
	return true
}

func (c *Client) SearchImage(image string) ImageProperty {
}*/
