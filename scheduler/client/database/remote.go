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

/*
func get(url string) (Response, error) {

}*/

func (c *RemoteClient) GetInfo() (Response, error) {

	c.m.RLock()
	defer c.m.RUnlock()

	rp := Response{}

	resp, err := c.client.DoAction("/api/info", common.Get)
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

func (c *RemoteClient) GetRepos() (Response, error) {
	c.m.RLock()
	defer c.m.RUnlock()

	var rp Response
	resp, err := c.client.DoAction("/api/repositories", common.Get)
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

func (c *RemoteClient) ListRepoTags(name string, repo string) (Response, error) {

	if len(repo) == 0 {
		panic("invalid argment")
	}

	c.m.RLock()
	defer c.m.RUnlock()

	var rp Response
	var url string

	if len(name) != 0 {
		url = "/api/repository/" + repo

	} else {
		url = "/api/repository/" + name + "/" + repo
	}
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

func (c *RemoteClient) GetUserRepos(user string) (Response, error) {

	if len(user) == 0 {
		panic("invalid argment")
	}

	c.m.RLock()
	defer c.m.RUnlock()

	var rp Response
	resp, err := c.client.DoAction("/api/repositories/user/"+user, common.Get)
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

func (c *RemoteClient) GetNsRepos(ns string) (Response, error) {

	if len(ns) == 0 {
		panic("invalid argment")
	}

	c.m.RLock()
	defer c.m.RUnlock()

	var rp Response
	resp, err := c.client.DoAction("/api/repositories/"+ns, common.Get)
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

func (c *RemoteClient) GetTagImage(name string, repo string, tag string) (Response, error) {

	if len(repo) == 0 || len(tag) == 0 {
		panic("invalid arguments")
	}

	c.m.RLock()
	defer c.m.RUnlock()

	var rp Response
	var url string
	if len(name) == 0 {
		url = "/api/tag/" + name + "/" + repo + "/" + tag
	} else {
		url = "/api/tag/" + repo + "/" + tag
	}
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

func (c *RemoteClient) GetNamespaces() (Response, error) {

	c.m.RLock()
	defer c.m.RUnlock()

	var rp Response
	resp, err := c.client.DoAction("/api/namespaces", common.Get)
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

func (c *RemoteClient) GetSpecificNamespace(ns string) (Response, error) {
	if len(ns) == 0 {
		panic("invalid arguments")
	}

	c.m.RLock()
	defer c.m.RUnlock()

	var rp Response
	resp, err := c.client.DoAction("/api/namespace/"+ns, common.Get)
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

func (c *RemoteClient) AddNamespace(ns Namespace) (Response, error) {
	if len(ns.Id) == 0 {
		panic("invalid arguments")
	}

	c.m.Lock()
	defer c.m.Unlock()

	var rp Response

	byteData, err := json.Marshal(ns)
	if err != nil {
		panic(err)
	}

	resp, err := c.client.DoPost("/api/namespace", byteData)
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

func (c *RemoteClient) GetNsUgroup(ns string) (Response, error) {
	if len(ns) == 0 {
		panic("invalid arguments")
	}

	c.m.RLock()
	defer c.m.RUnlock()

	var rp Response
	resp, err := c.client.DoAction("/api/grp/"+ns, common.Get)
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

func (c *RemoteClient) AddUgroup(ug UserGroup) (Response, error) {
	if len(ug.GroupName) == 0 {
		panic("invalid argument")
	}

	c.m.Lock()
	defer c.m.Unlock()

	var rp Response

	byteData, err := json.Marshal(ug)
	if err != nil {
		panic(err)
	}

	resp, err := c.client.DoPost("/api/grp", byteData)
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

func (c *RemoteClient) ListAccounts() (Response, error) {

	c.m.RLock()
	defer c.m.RUnlock()

	var rp Response
	resp, err := c.client.DoAction("/api/accounts", common.Get)
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

func (c *RemoteClient) AddUserAccount(user UserInfo) (Response, error) {
	if len(user.Id) == 0 || len(user.Password) == 0 {
		panic("invalid arguments")
	}

	c.m.Lock()
	defer c.m.Unlock()

	var rp Response

	byteData, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	resp, err := c.client.DoPost("/api/account", byteData)
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
func (c *RemoteClient) GetAccountInfo(user string) (Response, error) {
	if len(user) == 0 {
		panic("invalid argument")
	}
	c.m.RLock()
	defer c.m.RUnlock()

	var rp Response
	resp, err := c.client.DoAction("/api/account/"+user, common.Get)
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
