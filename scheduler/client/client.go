package client

import (
	"scheduler/client/database"
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
	//  Registry RegistryOperation
	//  SysInfo SysInfoOperation
	//  lock - 数据同步时使用?还是在subClient中实现?
}

func constructClient() *Client {
	var err error
	client := new(Client)
	client.database, err = database.GetDatabaseClient()
	if err != nil {
		panic(err)
	}

	//client.Registry = newRegistryClient(client)
	//client.SysInfo = newSysInfoClient(client)

	return client
}

func NewClient(opts ClientOpts) (*Client, error) {
	client := constructClient()
	//	client.BaseClient.Opts = opts
	return client, nil
}
