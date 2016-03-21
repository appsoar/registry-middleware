package database

import (
	"encoding/json"
	"fmt"
)

var (
	databaseClients map[string]DatabaseClient
)

type UserInfo struct {
	Id       string  `json:"_id"`
	Password string  `json:"password"`
	NickName string  `json:"nick_name"`
	Avatar   string  `json:"avatar"`
	JoinTime float64 `json:"join_time"`
}

type UserGroup struct {
	Id         int     `json:"_id"`
	GroupName  string  `json:"group_name"`
	Namespace  string  `json:"namespace"`
	CreateTime float64 `json:"create_time"`
	Desc       string  `json:"desc"`
}

type Namespace struct {
	Id         string  `json:"_id"`
	OwnerId    string  `json:"owener_id"`
	Desc       string  `json:"desc"`
	Perms      string  `json:"permission"`
	CreateTime float64 `json:"create_time"`
}

type Response struct {
	Content json.RawMessage
	Message string
	Result  int
}

type DatabaseClient interface {
	/*user,namespace,repo statistic*/
	GetInfo() (Response, error)

	/*----repo*/
	GetRepos() (Response, error)
	ListRepoTags(string, string) (Response, error)
	GetUserRepos(string) (Response, error)
	GetNsRepos(string) (Response, error)
	GetTagImage(string, string, string) (Response, error)

	/*----Namespace*/
	GetNamespaces() (Response, error)
	GetSpecificNamespace(ns string) (Response, error)
	AddNamespace(Namespace) (Response, error)

	/*----UserGroup*/
	GetNsUgroup(string) (Response, error)
	AddUgroup(UserGroup) (Response, error)

	/*----Accounts*/
	ListAccounts() (Response, error)
	AddUserAccount(UserInfo) (Response, error)
	GetAccountInfo(string) (Response, error)
}

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
	name := "remote"

	if client, ok := databaseClients[name]; ok {
		return client, nil
	}
	return nil, fmt.Errorf("databaseClient not support.")
}

func init() {
	//挂载后端database client钩子
	//
}
