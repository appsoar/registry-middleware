package database

import (
	"encoding/json"
	"io/ioutil"
	"scheduler/client/common"
	"sync"
)

type RemoteClient struct {
	client common.BaseClient
	m      *sync.RWMutex
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

func doGet(url string) (Response, error) {

	c.m.RLock()
	defer c.m.RUnlock()
	rp := Response{}

	resp, err := c.client.DoAction(url, common.Get)
	if err != nil {
		return rp, err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rp, err
	}

	err = json.Unmarshal(byteContent, &rp)
	if err != nil {
		return rp, err
	}
	return rp, nil
}

func doPost(url string, byteData []byte) (Response, error) {

	c.m.Lock()
	defer c.m.Unlock()

	var rp Response

	resp, err := c.client.DoPost(url, byteData)
	if err != nil {
		return rp, err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return rp, err
	}

	err = json.Unmarshal(byteContent, &rp)
	if err != nil {
		return rp, err
	}
	return rp, nil
}

func (c *RemoteClient) GetInfo() (Response, error) {

	url := "/api/info"
	rp, err := doGet(url)
	return rp, err
}

func (c *RemoteClient) GetRepos() (Response, error) {
	url := "/api/repositories"
	rp, err := doGet(url)
	return rp, err

}

func (c *RemoteClient) ListRepoTags(name string, repo string) (Response, error) {

	if len(repo) == 0 {
		panic("invalid argment")
	}

	var url string
	if len(name) != 0 {
		url = "/api/repository/" + repo
	} else {
		url = "/api/repository/" + name + "/" + repo
	}

	rp, err := doGet(url)
	return rp, err

}

func (c *RemoteClient) GetUserRepos(user string) (Response, error) {

	if len(user) == 0 {
		panic("invalid argment")
	}

	url := "/api/repositories/user/" + user
	rp, err := goGet(url)
	return rp, err
}

func (c *RemoteClient) GetNsRepos(ns string) (Response, error) {

	if len(ns) == 0 {
		panic("invalid argment")
	}

	url := "/api/repositories/" + ns
	rp, err := goGet(url)
	return rp, err
}

func (c *RemoteClient) GetTagImage(name string, repo string, tag string) (Response, error) {

	if len(repo) == 0 || len(tag) == 0 {
		panic("invalid arguments")
	}

	var url string
	if len(name) == 0 {
		url = "/api/tag/" + name + "/" + repo + "/" + tag
	} else {
		url = "/api/tag/" + repo + "/" + tag
	}

	rp, err := doGet(url)
	return rp, err
}

func (c *RemoteClient) GetNamespaces() (Response, error) {

	url := "/api/namespaces"
	rp, err := doGet(url)
	return rp, err
}

func (c *RemoteClient) GetSpecificNamespace(ns string) (Response, error) {
	if len(ns) == 0 {
		panic("invalid arguments")
	}

	url := "/api/namespace/" + ns
	rp, err := doGet(url)
	return rp, err
}

func (c *RemoteClient) AddNamespace(ns Namespace) (Response, error) {
	if len(ns.Id) == 0 {
		panic("invalid arguments")
	}

	byteData, err := json.Marshal(ns)
	if err != nil {
		panic(err)
	}

	url := "/api/namespace"
	rp, err := doPost(url, byteData)
	return rp, err
}

func (c *RemoteClient) GetNsUgroup(ns string) (Response, error) {
	if len(ns) == 0 {
		panic("invalid arguments")
	}

	url := "/api/grp/" + ns
	rp, err := doGet(url)
	return rp, err
}

func (c *RemoteClient) AddUgroup(ug UserGroup) (Response, error) {
	if len(ug.GroupName) == 0 {
		panic("invalid argument")
	}

	byteData, err := json.Marshal(ug)
	if err != nil {
		panic(err)
	}

	url := "/api/grp"
	rp, err := doPost(url, byteData)
	return rp, err
}

func (c *RemoteClient) ListAccounts() (Response, error) {

	url := "/api/accounts"
	rp, err := doGet(url)
	return rp, err
}

func (c *RemoteClient) AddUserAccount(user UserInfo) (Response, error) {
	if len(user.Id) == 0 || len(user.Password) == 0 {
		panic("invalid arguments")
	}
	byteData, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}
	url := "/api/grp"
	rp, err := doPost(url, byteData)
	return rp, err

}

func (c *RemoteClient) GetAccountInfo(user string) (Response, error) {
	if len(user) == 0 {
		panic("invalid argument")
	}
	url := "/api/account/" + user
	rp, err := doGet(url)
	return rp, err
}
