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

func init() {
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

func (c *RemoteClient) GetUserAccount(user string) (Response, error) {

	if len(user) == 0 {
		panic("invalid argument")
	}
	log.Logger.Debug("get user account:" + user)
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

func (c *RemoteClient) GetAccounts() (Response, error) {
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

/*
	GetInfo(string) (Response, error)
	//	DelImageTag(string) error

	GetRepos() (Response, error)
	GetSpecificRepos(string) (Response, error)
	GetTagImage(string, string, string) (Response, error)
*/

/*
	GetInfo(string) (Response, error)
	//	DelImageTag(string) error

	GetRepos() (Response, error)
	GetSpecificRepos(string) (Response, error)
	GetTagImage(string, string, string) (Response, error)
*/

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

func (c *RemoteClient) AddNamespace(ns string) (Response, error) {
	if len(ns) == 0 {
		panic("invalid arguments")
	}

	c.m.Lock()
	defer c.m.Unlock()

	var rp Response
	resp, err := c.client.DoAction("/api/namespace/"+ns, common.Post)
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

func (c *RemoteClient) ListNsUgroup(ns string, ug string) (Response, error) {
	if len(ns) == 0 || len(ug) == 0 {
		panic("invalid arguments")
	}

	c.m.RLock()
	defer c.m.RUnlock()

	var rp Response
	resp, err := c.client.DoAction("/api/grp/"+ns+"/"+ug, common.Get)
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

func (c *RemoteClient) AddNsUgroup(ns string, ug string) (Response, error) {
	if len(ns) == 0 || len(ug) == 0 {
		panic("invalid arguments")
	}

	c.m.Lock()
	defer c.m.Unlock()

	var rp Response
	resp, err := c.client.DoAction("/api/grp/"+ns+"/"+ug, common.Post)
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

func (c *RemoteClient) AddAccount(user string) (Response, error) {
	if len(user) == 0 {
		panic("invalid arguments")
	}

	c.m.Lock()
	defer c.m.Unlock()

	var rp Response
	resp, err := c.client.DoAction("/api/account/"+user, common.Post)
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
