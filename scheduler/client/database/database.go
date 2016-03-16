package database

import (
	"fmt"
	"time"
)

var (
	databaseClients map[string]DatabaseClient
)

/*
type UserStatInfo struct {
	Content InfoContent `json:"content"`
	Message string      `json:"message"`
	Result  int         `json:"result"`
}

type InfoContent struct {
	NamespaceNum  int `json:"namespace"`
	UserNum       int `json:"user"`
	RepositoryNum int `json:"repository"`
}

type Respository struct {
	Content RespositoryContent `json:"content"`
	Message string             `json:"message"`
	Result  int                `json:"result"`
}

type RespositoryContent struct {
	PushTime  time.Time `json:"push_time"`
	UserId    string    `json:"user_id"`
	Namespace string    `json:"namespace"`
	IsPublic  bool      `json:"is_public"`
	Desc      string    `json:"desc"`
	Id        string    `json:"_id"`
	Delete    time.Time `json:"delete"`
}

type Account struct {
	Content AccountContent `json:"content"`
	Message string         `json:"message"`
	Result  int            `json:"result"`
}

type AccountContent struct {
	NickName string    `json:"nick_name"`
	UserID   string    `json:"user_id"`
	Avatar   string    `json:"avatar"` //头像
	JoinTime time.Time `json:"join_time"`
	Password string    `json:"password"`
}
*/
type Message struct {
	Content map[interface{}]interface{}
	Message string
	Result  int
}

type DatabaseClient interface {
	GetUserStatInfo() (Message, error)
	GetAccount(string) (Message, error)
	//AddUser(UserInfo) error
	//DelUser(UserInfo) error
	//GetNamespaceInfo() (Namespace, error)
	//AddNamespace(Namespace) error
	//DelNamespace(Namespace) error

	GetInfo(string) (Message, error)
	//	DelImageTag(string) error

	GetRes
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
