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
