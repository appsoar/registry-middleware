package client

import (
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

func (c *Client) GetCpuUsage() (interface{}, error) {
	//	cpuUsage,err := c.sysInfo.
	//c.sysInfo.GetCpuUsage()
	return nil, nil
}

func (c *Client) GetRamUsage() (interface{}, error) {
	//	cpuUsage,err := c.sysInfo.
	//c.sysInfo.GetRamUsage()
	return nil, nil
}

func (c *Client) GetDiskUsage() (interface{}, error) {
	//	cpuUsage,err := c.sysInfo.
	//c.sysInfo.GetRamUsage()
	return nil, nil
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
