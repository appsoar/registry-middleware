package database

import (
	"fmt"
	//	"os"
	//	 "scheduler/client/database/local"
)

type UserInfo struct {
	username string
	password string
}

var (
	databaseClients map[string]DatabaseClient
)

type DatabaseClient interface {
	GetUserInfo(username string) (UserInfo, error)
}

/*
type DatabaseClient struct {
	client *Client
	BaseClient
}

func (c *DatabaseClient) GetUseInfo(username string) UserInfo {
	//c.client.do****
	return UserInfo{username: username, password: "12345"} //test
}

func newDatabaseClient() *DatabaseClient {
	return &DatabaseClient{
		client: client,
	}
}
*/

//注册数据库客户端
func RegisterDatabaseClient(name string, client DatabaseClient) error {
	if databaseClients == nil {
		databaseClients = make(map[string]DatabaseClient)
	}

	if _, exists := databaseClients[name]; exists {
		return fmt.Errorf("databaseClient already registered")
	}

	databaseClients[name] = client
	return nil

}

func GetDatabaseClient() (DatabaseClient, error) {
	/*
		name := os.Getenv("DatabaseClient")
		if name != nil {
			return nil, fmt.Errorf("databaseClient not choose")
		}
	*/
	name := "local"

	if client, ok := databaseClients[name]; ok {
		return client, nil
	}
	return nil, fmt.Errorf("databaseClient not support.")
}

func init() {
	//挂载后端database client钩子
	//
}
