package database

import (
	"fmt"
)

type LocalClient struct {
	db map[string]string
}

func init() {
	//localclient := &LocalClient{}
	//RegisterDatabaseClient("local", localclient)

}

func (c *LocalClient) GetUserInfo(username string) (UserInfo, error) {
	if password, ok := c.db[username]; ok {
		return UserInfo{username, password}, nil
	}
	return UserInfo{}, fmt.Errorf("invalid username")
}
