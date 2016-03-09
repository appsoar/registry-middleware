package registry

import (
	"io/ioutil"
	//	"net/http"
	"fmt"
	"scheduler/client/common"
	"sync"
)

type DirectClient struct {
	//这里需要有锁进行同步
	client common.BaseClient
	m      *sync.RWMutex
}

func init() {
	//初始化registry 服务器的配置参数
	opts := &common.ClientOpts{
		Url:       "http://192.168.4.32:5050",
		AccessKey: "",
		SecretKey: "",
		Timeout:   0,
	}

	directClient := &DirectClient{
		client: common.BaseClient{Opts: opts},
		m:      new(sync.RWMutex),
	}

	RegisterRegistryClient("direct", directClient)
}

func (c *DirectClient) ListImages() (interface{}, error) {
	c.m.RLock()
	fmt.Println("listImage locking.......")
	defer func() {
		fmt.Println("listImage unlocking.....")
		c.m.RUnlock()
	}()
	resp, err := c.client.DoAction("/v2/_catalog", common.Get)
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

func (c *DirectClient) GetImageTags(image string) (interface{}, error) {
	c.m.RLock()
	fmt.Println("GetImageTags Locking......")
	defer func() {
		fmt.Println("GetImageTags Unlock........")
		c.m.RUnlock()
	}()
	resp, err := c.client.DoAction("/v2/"+image+"/tags/list", common.Get)
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

func (c *DirectClient) GetImageDigest(image string, tag string) (interface{}, error) {
	c.m.RLock()
	defer c.m.RUnlock()

	resp, err := c.client.DoAction("/v2"+image+"/manifests/"+tag, common.Get)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	for k, v := range resp.Header {
		if k == "Docker-Content-Digest" {
			digest := v[0]
			return digest, nil
		}

	}
	return nil, fmt.Errorf("headers don't have `Docker-Content-Digest` field")

}

func (c *DirectClient) DeleteImageTag(image string, tag string) error {
	c.m.Lock()
	defer c.m.Unlock()

	digest, err := c.GetImageDigest(image, tag)
	if err != nil {
		return err
	}

	str, ok := digest.(string)
	if !ok {
		panic("digest isnot string")
	}
	resp, err := c.client.DoAction("/v2/"+image+"/manifests/"+str, common.Delete)
	if err != nil {
		return err
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != 202 {
		return fmt.Errorf("delelte image fail")
	}
	return nil
}
