package database

import (
	"encoding/json"
	"fmt"
)

var (
	databaseClients map[string]DatabaseClient
)

//这里部分错误，应当做服务器内部错误返回
//其余当做用户无效请求Forbidden
const (
	EPermission  = 10108 //用户权限不够
	EDbException = 10501 //数据库异常
	ENoRecord    = 10503 // 记录不存在
	EMissingId   = 10509 //记录信息没有Id

	EInvalidFilter = 10510 //过滤条件不合法
	EParameter     = 11001 //参数错误
	ENotInterface  = 11002 //没有实现对应的接口

	EIncompleteUserInfo  = 12001 //用户信息不完整
	EUserExists          = 12002 //用户已经存在
	EIncompleteGroupInfo = 12003 //分组信息不完整
	EGroupExists         = 12004 //分组已存在
	EInvalidNsInfo       = 12005 //命名空间信息无效
	ENsExists            = 12006 //命名空间已经存在
)

type EDatabase struct {
	Code int
	Msg  string
}

func (e EDatabase) Error() string {
	return e.Msg
}

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

type DatabaseClient interface {
	/*user,namespace,repo statistic*/
	GetInfo() (json.RawMessage, error)

	/*----repo*/
	GetRepos() (json.RawMessage, error)
	ListRepoTags(string, string) (json.RawMessage, error)
	GetUserRepos(string) (json.RawMessage, error)
	GetNsRepos(string) (json.RawMessage, error)
	GetTagImage(string, string, string) (json.RawMessage, error)

	/*----Namespace*/
	GetNamespaces() (json.RawMessage, error)
	GetSpecificNamespace(ns string) (json.RawMessage, error)
	AddNamespace(Namespace) (json.RawMessage, error)

	/*----UserGroup*/
	GetNsUgroup(string) (json.RawMessage, error)
	AddUgroup(UserGroup) (json.RawMessage, error)

	/*----Accounts*/
	ListAccounts() (json.RawMessage, error)
	AddUserAccount(UserInfo) (json.RawMessage, error)
	GetAccountInfo(string) (json.RawMessage, error)
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
