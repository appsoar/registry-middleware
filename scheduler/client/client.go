package client

import (
	"errors"
	"fmt"
	"scheduler/client/database"
	"scheduler/client/registry"
	"scheduler/client/sysinfo"

	"encoding/json"
	"scheduler/log"
)

//UserChannel chan int
//ImageChannel chan int
//NamespaceChannel chan int

//所有访问的总接口
//在每个subClient实现lock
type Client struct {
	database database.DatabaseClient
	registry registry.RegistryClient
	sysInfo  sysinfo.SysInfoClient
}

func constructClient() *Client {
	var err error
	client := new(Client)

	client.database, err = database.GetDatabaseClient()
	if err != nil {
		panic(err)
	}

	client.sysInfo, err = sysinfo.GetSysInfoClient()
	if err != nil {
		panic(err)
	}

	client.registry, err = registry.GetRegistryClient()
	if err != nil {
		panic(err)
	}

	return client
}

func NewClient() (*Client, error) {
	client := constructClient()
	return client, nil
}

/*====================获取系统统计信息====================*/
type SysInfo struct {
	CpuUsage      int    `json:"cpuUsage"`
	TotalRam      uint64 `json:"totalRam"`
	AvailableRam  uint64 `json:"availableRam"`
	TotalDisk     uint64 `json:"totalDisk"`
	AvailableDisk uint64 `json:"availableDisk"`
	//	NetStat       []sysinfo.NetStat `json:"netStat"`
}

func (c *Client) GetNetIfs() ([]byte, error) {
	tmp, err := c.sysInfo.GetNetIfs()
	if err == nil {
		if ifs, ok := tmp.([]string); ok {
			ifsbytes, err := json.Marshal(ifs)
			if err != nil {
				return []byte{}, err
			}
			return ifsbytes, err
		}

	}
	return []byte{}, err
}

func (c *Client) GetNetIfStat(If string) ([]byte, error) {
	tmp, err := c.sysInfo.GetNetIfStat(If)
	if err == nil {
		if ifstat, ok := tmp.(sysinfo.NetStat); ok {
			ifbytes, err := json.Marshal(ifstat)
			if err != nil {
				return []byte{}, err
			}
			return ifbytes, nil
		} else {
			fmt.Println("type assertion fail")
		}

	}
	return []byte{}, err
}

func (c *Client) GetSysInfo() (SysInfo, error) {

	info := new(SysInfo)
	var cpuchan = make(chan int)
	var ramchan = []chan uint64{
		make(chan uint64),
		make(chan uint64),
	}
	var diskchan = []chan uint64{
		make(chan uint64),
		make(chan uint64),
	}
	//	var netchan = make(chan []sysinfo.NetStat)

	go func() {
		cpuUsage, err := c.sysInfo.GetCpuUsage()
		if err != nil {
			log.Logger.Error("getCpuUsage fail")
			panic(err)
		}
		cpuchan <- cpuUsage
	}()

	go func() {
		totalRam, availableRam, err := c.sysInfo.GetRamStat()
		if err != nil {
			log.Logger.Error("getRamStat fail")
			panic(err)
		}
		ramchan[0] <- totalRam
		ramchan[1] <- availableRam
	}()

	go func() {
		totalDisk, availableDisk, err := c.sysInfo.GetDiskStat()
		if err != nil {
			log.Logger.Error("get diskstat fail")
			panic(err)
		}
		diskchan[0] <- totalDisk
		diskchan[1] <- availableDisk
	}()
	/*
		go func() {
			netStat, err := c.sysInfo.GetNetStat()
			if err != nil {
				log.Logger.Error("get netstat fail")
				panic(err)
			}
			netchan <- netStat
		}()
	*/
	info.CpuUsage = <-cpuchan
	info.TotalRam = <-ramchan[0]
	info.AvailableRam = <-ramchan[1]
	info.TotalDisk = <-diskchan[0]
	info.AvailableDisk = <-diskchan[1]
	//	info.NetStat = <-netchan
	/*
		log.Logger.Debug("cpu:%v, total:%v,avail:%v, total:%v,avail:%v\n", info.CpuUsage, info.TotalRam, info.AvailableRam, info.TotalRam, info.AvailableDisk)
		log.Logger.Debug("%+v", info.NetStat)

	*/
	return *info, nil

}

/*

func (c *Client) ListImages() ([]string, error) {
	images, err := c.registry.ListImages()
	if err != nil {
		return nil, err
	}
	list, ok := images.([]string)
	if !ok {
		panic("low-level registry api have changed")
	}

	return list, nil
}

func (c *Client) GetImageTags(image string) ([]string, error) {

	tags, err := c.registry.GetImageTags(image)
	if err != nil {
		return nil, err
	}
	list, ok := tags.([]string)
	if !ok {
		panic("low-level registry api have changed")
	}
	return list, nil
}

func (c *Client) GetImageDigest(image string, tag string) (string, error) {

	digest, err := c.registry.GetImageDigest(image, tag)
	if err != nil {
		return "", err
	}
	dg, ok := digest.(string)
	if !ok {
		panic("low-level registry api have changed")
	}
	return dg, nil
}

func (c *Client) DeleteImageDigest(image string, tag string) error {

	err := c.registry.DeleteImageTag(image, tag)
	return err
}
*/

/*===============获取用户数量,命名空间,镜像数统计===============*/

func (c *Client) GetUserStats() (resp []byte, err error) {
	respRec, err := c.database.GetInfo()
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}

	/*
		var rp database.Response
		err = json.Unmarshal(resp, &rp)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(rp.Content, &us)
		if err != nil {
			panic(err)
		}
	*/

	return
}

func (c *Client) GetAccounts() (resp []byte, err error) {

	respRec, err := c.database.ListAccounts()
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

/*解析后的用户账号信息*/
func (c *Client) GetUserAccountDecoded(user string) (ui database.UserInfo, err error) {

	if len(user) == 0 {
		log.Logger.Error("invalid argument...")
		panic("invalid argument..")
	}
	log.Logger.Debug("get uesr account")
	respRec, err := c.database.GetAccountInfo(user)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}

	var rp database.Response
	err = json.Unmarshal(resp, &rp)
	if err == nil {
		log.Logger.Debug(string(resp))
		err = json.Unmarshal(rp.Content, &ui)

	}
	return
}

func (c *Client) GetUserAccount(user string) (resp []byte, err error) {

	if len(user) == 0 {
		log.Logger.Error("invalid argument...")
		panic("invalid argument..")
	}
	respRec, err := c.database.GetAccountInfo(user)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

func (c *Client) AddUserAccount(user database.UserInfo) (resp []byte, err error) {

	if len(user.Id) == 0 {
		log.Logger.Error("invalid argument...")
		panic("invalid argument..")
	}
	respRec, err := c.database.AddUserAccount(user)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

func (c *Client) UpdateUserAccount() (resp []byte, err error) {

	respRec, err := c.database.UpdateAccount()
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

func (c *Client) DeleteUserAccount(user string) (resp []byte, err error) {

	respRec, err := c.database.DeleteAccount(user)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

/*=============Repositories===================*/

func (c *Client) GetRepositories() (resp []byte, err error) {
	respRec, err := c.database.GetRepos()
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

func (c *Client) GetRepositoriesDecoded() (repo []database.Repository, err error) {
	respRec, err := c.database.GetRepos()
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}

	var rp database.Response
	err = json.Unmarshal(resp, &rp)
	if err != nil {
		err = errors.New("json parse fail")
	}

	err = json.Unmarshal(rp.Content, &repo)
	if err != nil {
		err = errors.New("json parse fail")
	}
	return
}

func (c *Client) ListRepoTags(repoName string) (resp []byte, err error) {
	respRec, err := c.database.ListRepoTags(repoName)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}

	return
}

func (c *Client) GetNsRepos(ns string) (resp []byte, err error) {
	respRec, err := c.database.GetNsRepos(ns)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return

}

func (c *Client) GetUserRepos(user string) (resp []byte, err error) {
	respRec, err := c.database.GetUserRepos(user)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

/*===================镜像===================*/

func (c *Client) GetTagImageDecoded(repoName string, tagName string) (tag database.TagInfo, err error) {
	var respUmr database.Response
	respRec, err := c.database.GetTagImage(repoName, tagName)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}

	err = json.Unmarshal(resp, &respUmr)
	if err != nil {
		log.Logger.Error("json parse fail")
		panic(err)
	}
	err = json.Unmarshal(respUmr.Content, &tag)
	return
}

func (c *Client) GetTagImage(repoName string, tagName string) (resp []byte, err error) {
	respRec, err := c.database.GetTagImage(repoName, tagName)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

/*================ 命名空间 ==================*/

func (c *Client) GetNamespacesDecoded() (ns []database.Namespace, err error) {
	var rp database.Response
	respRec, err := c.database.GetNamespaces()
	if err != nil {
		return
	}

	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}

	err = json.Unmarshal(resp, &rp)
	if err != nil {
		log.Logger.Error("json parse fail")
		panic(err)
	}
	err = json.Unmarshal(rp.Content, &ns)
	if err != nil {
		log.Logger.Error("json parse fail")
		panic(err)
	}
	return
}

func (c *Client) GetNamespaces() (resp []byte, err error) {
	respRec, err := c.database.GetNamespaces()
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

func (c *Client) GetSpecificNamespace(ns string) (resp []byte, err error) {
	if len(ns) == 0 {
		log.Logger.Error("invalid argument")
		panic("invalid ns")
	}
	respRec, err := c.database.GetSpecificNamespace(ns)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

func (c *Client) AddNamespace(ns database.Namespace) (resp []byte, err error) {
	respRec, err := c.database.AddNamespace(ns)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

func (c *Client) DeleteNamespace(ns string) (resp []byte, err error) {
	respRec, err := c.database.DeleteNamespace(ns)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

func (c *Client) UpdateNamespace() (resp []byte, err error) {
	respRec, err := c.database.UpdateNamespace()
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

/* ============== 用户组 ==================*/

func (c *Client) GetNsUgroup(ns string) (resp []byte, err error) {
	respRec, err := c.database.GetNsUgroup(ns)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

func (c *Client) AddUgroup(ug database.UserGroup) (resp []byte, err error) {
	respRec, err := c.database.AddUgroup(ug)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

func (c *Client) GetUgroup(gid string) (resp []byte, err error) {
	respRec, err := c.database.GetUgroup(gid)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

func (c *Client) UpdateUgroup() (resp []byte, err error) {
	respRec, err := c.database.UpdateUgroup()
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

func (c *Client) DeleteUgroup(gid string) (resp []byte, err error) {
	respRec, err := c.database.DeleteUgroup(gid)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}

func (c *Client) GetLog(lo string) (resp []byte, err error) {
	respRec, err := c.database.GetLog(lo)
	if err != nil {
		return
	}
	resp, ok := respRec.([]byte)
	if !ok {
		panic("type assertion fail")
	}
	return
}
