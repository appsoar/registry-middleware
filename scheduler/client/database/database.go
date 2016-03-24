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

type UserStats struct {
	User       float64 `json:"user"`
	repository float64 `json:"repository"`
	Namespace  float64 `json:"namespace"`
}
type Repository struct {
	Id         string  `json:"_id"`
	Namespace  string  `json:"namespace"`
	User       string  `json:"user_id"`
	PushTime   float64 `json:"push_time"`
	Desc       string  `json:"desc"`
	Public     bool    `json:"is_public"`
	DeleteTime float64 `json:"delete"`
}

type TagInfo struct {
	Id          int    `json:"_id"`
	UserID      string `json:"user_id"`
	Respository string `json:"repository"`
	TagName     string `json:"tag_name"`
	Size        int    `json:"size"`
	Digest      string `json:"digest"`
	CreateTime  int    `json:"create_time"`
	Delete      int    `json:"delete"`
	PullNum     int    `json:"pull_num"`
}

type Namespace struct {
	Id         string  `json:"_id"`
	OwnerId    string  `json:"_id"`
	Desc       string  `json:"desc"`
	Permission string  `json:"public"`
	CreateTime float64 `json:"create_time"`
}

type Response struct {
	Content json.RawMessage
	Message string
	Result  int
}

type DatabaseClient interface {
	/* user,namespace,repos number statistic*/
	GetInfo() (interface{}, error) /*UserStats*/

	/*----repo ---- struct : Repository*/
	/*list all repos*/
	GetRepos() (interface{}, error)
	/*list all repos of user*/
	GetUserRepos(user string) (interface{}, error)
	/*list all repos under ns*/
	GetNsRepos(ns string) (interface{}, error)

	/*-----Tag ----- struct : TagInfo*/
	GetTagImage(string, string) (interface{}, error)
	/*list repo's tags or tags of specified ns|user's repo*/
	ListRepoTags(repo string) (interface{}, error)

	/*----Namespace --- struct : Namespace*/
	/*list all ns: return []Namespace*/
	GetNamespaces() (interface{}, error)
	/*get info of specific ns*/
	GetSpecificNamespace(ns string) (interface{}, error)
	DeleteNamespace(ns string) (interface{}, error)
	UpdateNamespace() (interface{}, error)
	/*add a new ns*/
	AddNamespace(Namespace) (interface{}, error)

	/*----UserGroup ---- struct : UserGroup*/
	/*get user groups under ns*/
	GetNsUgroup(ns string) (interface{}, error)
	/*add a user ugroup*/
	AddUgroup(UserGroup) (interface{}, error)
	GetUgroup(ugid string) (interface{}, error)
	UpdateUgroup() (interface{}, error)
	DeleteUgroup(ugid string) (interface{}, error)

	/*----Accounts -- struct : UserInfo*/
	/*list all user accounts*/
	ListAccounts() (interface{}, error)
	/*add a new user accounts*/
	AddUserAccount(UserInfo) (interface{}, error)
	/*get specific uesr account*/
	GetAccountInfo(string) (interface{}, error)
	UpdateAccount() (interface{}, error)
	DeleteAccount(string) (interface{}, error)

	GetLog(lo string) (interface{}, error)
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
