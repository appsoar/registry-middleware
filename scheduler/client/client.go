package client

import (
	"errors"
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
	CpuUsage      int               `json:"cpuUsage"`
	TotalRam      uint64            `json:"totalRam"`
	AvailableRam  uint64            `json:"availableRam"`
	TotalDisk     uint64            `json:"totalDisk"`
	AvailableDisk uint64            `json:"availableDisk"`
	NetStat       []sysinfo.NetStat `json:"netStat"`
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
	var netchan = make(chan []sysinfo.NetStat)

	go func() {
		cpuUsage, err := c.sysInfo.GetCpuUsage()
		if err != nil {
			panic(err)
		}
		cpuchan <- cpuUsage
	}()

	go func() {
		totalRam, availableRam, err := c.sysInfo.GetRamStat()
		if err != nil {
			panic(err)
		}
		ramchan[0] <- totalRam
		ramchan[1] <- availableRam
	}()

	go func() {
		totalDisk, availableDisk, err := c.sysInfo.GetDiskStat()
		if err != nil {
			panic(err)
		}
		diskchan[0] <- totalDisk
		diskchan[1] <- availableDisk
	}()

	go func() {
		netStat, err := c.sysInfo.GetNetStat()
		if err != nil {
			panic(err)
		}
		netchan <- netStat
	}()

	info.CpuUsage = <-cpuchan
	info.TotalRam = <-ramchan[0]
	info.AvailableRam = <-ramchan[1]
	info.TotalDisk = <-diskchan[0]
	info.AvailableDisk = <-diskchan[1]
	info.NetStat = <-netchan

	log.Logger.Debug("cpu:%v, total:%v,avail:%v, total:%v,avail:%v\n", info.CpuUsage, info.TotalRam, info.AvailableRam, info.TotalRam, info.AvailableDisk)
	log.Logger.Debug("%+v", info.NetStat)

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
type UserStats struct {
	User       int `json:"user"`
	repository int `json:"repository"`
	Namespace  int `json:"namespace"`
}

func (c *Client) GetUserStats() (us UserStats, err error) {
	var resp json.RawMessage

	resp, err = c.database.GetInfo()
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &us)
	if err != nil {
		panic(err)
	}
	return
}

func (c *Client) GetAccounts() (resp []byte, err error) {

	resp, err = c.database.ListAccounts()
	return
	/*
		if e, ok := err.(database.Edatabase); ok {
			switch e.Code {
			case database.EPermission,
				database.EnoRecord,
				database.EMissingId,
				database.EInvalidFilter,
				database.EParameter,
				database.EIncompleteUserInfo,
				database.EUserExists,
				database.EIncompleteUserInfo,
				database.EGroupExits,
				database.EInvalidNsInfo,
				database.NsExists:
					errJson.

			}
		}
	*/
}

/*解析后的用户账号信息*/
func (c *Client) GetUserAccountDecoded(user string) (ui database.UserInfo, err error) {

	if len(user) == 0 {
		log.Logger.Error("invalid argument...")
		panic("invalid argument..")
	}
	log.Logger.Debug("get uesr account")
	resp, err := c.database.GetAccountInfo(user)
	if err != nil {
		log.Logger.Debug("GetAccount Info fail")
		return
	}
	//userinfo, ok := resp.Content.(UserInfo)
	err = json.Unmarshal(resp, &ui)
	return
}

func (c *Client) GetUserAccount(user string) (resp []byte, err error) {

	if len(user) == 0 {
		log.Logger.Error("invalid argument...")
		panic("invalid argument..")
	}
	log.Logger.Debug("get uesr account")
	resp, err = c.database.GetAccountInfo(user)
	return
}

func (c *Client) AddUserAccount(user database.UserInfo) (resp []byte, err error) {

	if len(user.Id) == 0 {
		log.Logger.Error("invalid argument...")
		panic("invalid argument..")
	}
	log.Logger.Debug("add uesr account")
	resp, err = c.database.AddUserAccount(user)
	return
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

func (c *Client) GetRepositories() (resp []byte, err error) {
	resp, err = c.database.GetRepos()
	return
}

func (c *Client) GetRepositoriesDecoded() (repo []Repository, err error) {
	var resp json.RawMessage
	resp, err = c.database.GetRepos()
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &repo)
	if err != nil {
		err = errors.New("json parse fail")
	}
	return
}

func (c *Client) ListRepoTags(usernameOrNamespace string, repoName string) (resp []byte, err error) {
	resp, err = c.database.ListRepoTags(usernameOrNamespace, repoName)

	return
}

func (c *Client) GetNsRepos(ns string) (resp []byte, err error) {
	resp, err = c.database.GetNsRepos(ns)
	return

}

func (c *Client) GetUserRepos(user string) (resp []byte, err error) {
	resp, err = c.database.GetUserRepos(user)
	return
}

/*===================镜像===================*/

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

func (c *Client) GetTagImageDecoded(usernameOrNamespace string, repoName string, tagName string) (tag TagInfo, err error) {
	var resp json.RawMessage
	resp, err = c.database.GetTagImage(usernameOrNamespace, repoName, tagName)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &tag)
	if err != nil {
		panic(err)
	}
	return
}

func (c *Client) GetTagImage(usernameOrNamespace string, repoName string, tagName string) (resp []byte, err error) {
	resp, err = c.database.GetTagImage(usernameOrNamespace, repoName, tagName)
	return
}

/*================ 命名空间 ==================*/

type Namespace struct {
	Id         string  `json:"_id"`
	OwnerId    string  `json:"_id"`
	Desc       string  `json:"desc"`
	Permission string  `json:"public"`
	CreateTime float64 `json:"create_time"`
}

func (c *Client) GetNamespacesDecoded() (ns []Namespace, err error) {
	var resp json.RawMessage
	resp, err = c.database.GetNamespaces()
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &ns)
	if err != nil {
		panic(err)
	}
	return
}

func (c *Client) GetNamespaces() (resp []byte, err error) {
	resp, err = c.database.GetNamespaces()
	return
}

func (c *Client) GetSpecificNamespace(ns string) (resp []byte, err error) {
	if len(ns) == 0 {
		panic("invalid ns")
	}
	resp, err = c.database.GetSpecificNamespace(ns)
	return
}
func (c *Client) AddNamespace(ns database.Namespace) (resp []byte, err error) {
	resp, err = c.database.AddNamespace(ns)
	return
}

/* ============== 用户组 ==================*/

func (c *Client) GetNsUgroup(ns string) (resp []byte, err error) {
	resp, err = c.database.GetNsUgroup(ns)
	return
}

func (c *Client) AddUgroup(ug database.UserGroup) (resp []byte, err error) {
	resp, err = c.database.AddUgroup(ug)
	return
}
