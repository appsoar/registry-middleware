package registry

import (
	"io/ioutil"
	//"net/http"
	"scheduler/client/common"
	"sync"
)

type RemoteClient struct {
	//这里需要有锁进行同步
	client common.BaseClient
	m      *sync.RWMutex
}

func init() {
	//初始化registry 服务器的配置参数
	opts := &common.ClientOpts{
		Url:       "",
		AccessKey: "",
		SecretKey: "",
		Timeout:   0,
	}

	remoteClient := &RemoteClient{
		client: common.BaseClient{Opts: opts},
		m:      new(sync.RWMutex),
	}

	RegisterRegistryClient("remote", remoteClient)
}

func (c *RemoteClient) ListImages() (interface{}, error) {
	c.m.RLock()
	defer c.m.RUnlock()
	resp, err := c.client.DoAction("/images", common.Get)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return byteContent, nil
}

func (c *RemoteClient) GetImageTags(image string) (interface{}, error) {
	c.m.RLock()
	defer c.m.RUnlock()
	return nil, nil
}

func (c *RemoteClient) GetImageDigest(image string, tag string) (interface{}, error) {
	c.m.RLock()
	defer c.m.RUnlock()
	return nil, nil
}

func (c *RemoteClient) DeleteImageTag(image string, tag string) error {
	c.m.Lock()
	defer c.m.Unlock()
	return nil
}
