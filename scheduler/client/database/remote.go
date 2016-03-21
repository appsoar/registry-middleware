package database

import (
	"encoding/json"
	"io/ioutil"
	"scheduler/client/common"
	"scheduler/log"
	"sync"
)

type RemoteClient struct {
	client common.BaseClient
	m      *sync.RWMutex
}

type response struct {
	Content json.RawMessage
	Message string
	Result  int
}

func init() {
	/*
		Url := os.Getenv("DBURL")
		accessKey := os.Getenv("ACCESSKEY")
		secretKey := os.Getenv("SECRETKEY")
		strTimeout := os.Getenv("TIMEOUT")

		if len(Url) == 0 {
			//出错处理
			panic("not sp database server")
		}
		timeout := 0
		if len(strTimeout) != 0 {
			timeout, err := strconv.Atoi(strTimeout)
			if err != nil {
				log.Logger.Error("set database timeout fail: " + err.Error())
				log.Logger.Error("set timeout default to 0")
				timeout = 0
			}
		}

	*/
	opts := &common.ClientOpts{
		Url:       "http://192.168.12.112:8080",
		AccessKey: "",
		SecretKey: "",
		Timeout:   0,
	}

	remoteClient := &RemoteClient{
		client: common.BaseClient{Opts: opts},
		m:      new(sync.RWMutex),
	}

	RegisterDatabaseClient("remote", remoteClient)

}

func (c *RemoteClient) doGet(url string) (content json.RawMessage, err error) {

	c.m.RLock()
	defer c.m.RUnlock()
	rp := response{}

	resp, err := c.client.DoAction(url, common.Get)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(byteContent, &rp)
	if err != nil {
		return
	}

	if rp.Result != 0 {
		err = EDatabase{Code: rp.Result, Msg: rp.Message}
		return
	}
	content = rp.Content
	return
}

func (c *RemoteClient) doPost(url string, byteData []byte) (content json.RawMessage, err error) {

	c.m.Lock()
	defer c.m.Unlock()

	var rp response

	log.Logger.Debug("request body:" + string(byteData))
	resp, err := c.client.DoPost(url, byteData)
	if err != nil {
		log.Logger.Error(err.Error())
		return
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Logger.Error("ioutil Read All fail")
		return
	}

	err = json.Unmarshal(byteContent, &rp)
	if err != nil {
		log.Logger.Error("json decoded fail")
		return
	}

	if rp.Result != 0 {
		err = EDatabase{Code: rp.Result, Msg: rp.Message}
		return
	}
	content = rp.Content
	return
}

func (c *RemoteClient) GetInfo() (json.RawMessage, error) {

	url := "/api/info"
	rp, err := c.doGet(url)
	return rp, err
}

func (c *RemoteClient) GetRepos() (json.RawMessage, error) {
	url := "/api/repositories"
	rp, err := c.doGet(url)
	return rp, err

}

func (c *RemoteClient) ListRepoTags(name string, repo string) (json.RawMessage, error) {

	if len(repo) == 0 {
		panic("invalid argment")
	}

	var url string
	if len(name) != 0 {
		url = "/api/repository/" + repo
	} else {
		url = "/api/repository/" + name + "/" + repo
	}

	rp, err := c.doGet(url)
	return rp, err

}

func (c *RemoteClient) GetUserRepos(user string) (json.RawMessage, error) {

	if len(user) == 0 {
		panic("invalid argment")
	}

	url := "/api/repositories/user/" + user
	rp, err := c.doGet(url)
	return rp, err
}

func (c *RemoteClient) GetNsRepos(ns string) (json.RawMessage, error) {

	if len(ns) == 0 {
		panic("invalid argment")
	}

	url := "/api/repositories/" + ns
	rp, err := c.doGet(url)
	return rp, err
}

func (c *RemoteClient) GetTagImage(name string, repo string, tag string) (json.RawMessage, error) {

	if len(repo) == 0 || len(tag) == 0 {
		panic("invalid arguments")
	}

	var url string
	if len(name) == 0 {
		url = "/api/tag/" + name + "/" + repo + "/" + tag
	} else {
		url = "/api/tag/" + repo + "/" + tag
	}

	rp, err := c.doGet(url)
	return rp, err
}

func (c *RemoteClient) GetNamespaces() (json.RawMessage, error) {

	url := "/api/namespaces"
	rp, err := c.doGet(url)
	return rp, err
}

func (c *RemoteClient) GetSpecificNamespace(ns string) (json.RawMessage, error) {
	if len(ns) == 0 {
		panic("invalid arguments")
	}

	url := "/api/namespace/" + ns
	rp, err := c.doGet(url)
	return rp, err
}

func (c *RemoteClient) AddNamespace(ns Namespace) (json.RawMessage, error) {
	if len(ns.Id) == 0 {
		panic("invalid arguments")
	}

	byteData, err := json.Marshal(ns)
	if err != nil {
		panic(err)
	}

	url := "/api/namespace"
	rp, err := c.doPost(url, byteData)
	return rp, err
}

func (c *RemoteClient) GetNsUgroup(ns string) (json.RawMessage, error) {
	if len(ns) == 0 {
		panic("invalid arguments")
	}

	url := "/api/grp/" + ns
	rp, err := c.doGet(url)
	return rp, err
}

func (c *RemoteClient) AddUgroup(ug UserGroup) (json.RawMessage, error) {
	if len(ug.GroupName) == 0 {
		panic("invalid argument")
	}

	byteData, err := json.Marshal(ug)
	if err != nil {
		panic(err)
	}

	url := "/api/grp"
	rp, err := c.doPost(url, byteData)
	return rp, err
}

func (c *RemoteClient) ListAccounts() (json.RawMessage, error) {

	url := "/api/accounts"
	rp, err := c.doGet(url)
	return rp, err
}

func (c *RemoteClient) AddUserAccount(user UserInfo) (json.RawMessage, error) {
	if len(user.Id) == 0 || len(user.Password) == 0 {
		log.Logger.Error("User Account have empty Id or Password")
		panic("invalid arguments")
	}
	byteData, err := json.Marshal(user)
	if err != nil {
		log.Logger.Error("Json encoded fail")
		panic(err)
	}
	url := "/api/account"
	rp, err := c.doPost(url, byteData)
	if err != nil {
		log.Logger.Error(err.Error())
	}
	return rp, err

}

func (c *RemoteClient) GetAccountInfo(user string) (json.RawMessage, error) {
	if len(user) == 0 {
		panic("invalid argument")
	}
	url := "/api/account/" + user
	rp, err := c.doGet(url)
	return rp, err
}
