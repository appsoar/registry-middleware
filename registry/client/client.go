package client

//package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"registry/debug"
	"strconv"
	"time"
)

//默认超时时间
const (
	DefaultTimeOut time.Duration = time.Second * 10
)

type Op struct {
	Name string
}

var (
	Get    = Op{Name: "GET"}
	Head   = Op{Name: "HEAD"}
	Delete = Op{Name: "DELETE"}
	Put    = Op{Name: "PUT"}
	Post   = Op{Name: "POST"}
)

/*客户端选项*/
type ClientOpts struct {
	Url       string //http url
	AccessKey string
	SecretKey string
	Timeout   time.Duration //http请求超时限制
}

/*Api错误,保存http请求返回的所有信息*/
type ApiError struct {
	Msg        string
	Status     string
	Url        string
	Body       string
	StatusCode int
}

/*实现error的接口*/
func (e *ApiError) Error() string {
	return e.Msg
}

func newApiError(resp *http.Response, url string) *ApiError {
	contents, err := ioutil.ReadAll(resp.Body)
	var body string
	if err != nil {
		body = "Unreadable body"
	} else {
		body = string(contents)
	}

	data := map[string]interface{}{}
	if json.Unmarshal(contents, &data) == nil {
		buf := &bytes.Buffer{}
		for k, v := range data {
			if v == nil {
				continue
			}

			if buf.Len() > 0 {
				buf.WriteString(",")
			}

			fmt.Fprintf(buf, "%s=%v", k, v)
		}
		body = buf.String()

	}
	formattedMsg := fmt.Sprintf("Bad response statusCode [%d]. Status [%s]. Body: [%s] from [%s]", resp.StatusCode, resp.Status, body, url)
	return &ApiError{
		Url:        url,
		Msg:        formattedMsg,
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       body,
	}
}

func doAction(opts ClientOpts, op Op) (resp *http.Response, err error) {
	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeOut
	}
	client := &http.Client{Timeout: opts.Timeout}

	req, err := http.NewRequest(op.Name, opts.Url, nil)
	if err != nil {
		panic(err.Error())
	}

	req.SetBasicAuth(opts.AccessKey, opts.SecretKey)
	resp, err = client.Do(req)

	return

	//return nil
}

/*检测registry Api版本*/
func CheckVersion(opts ClientOpts) error {
	//	var respObject map[string][]string
	opts.Url = opts.Url + "/v2"
	_, err := doAction(opts, Get)
	return err
}

func GetImageDigest(opts ClientOpts, image string, tag string) (docker_content_digest string, err error) {

	opts.Url = opts.Url + "/" + "/v2/" + image + "/manifests/" + tag
	resp, err := doAction(opts, Get)
	if err != nil {
		return
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	for k, v := range resp.Header {
		if k == "Docker-Content-Digest" {
			//			digest = strings.TrimLeft(v[0], "sha256:")
			docker_content_digest = v[0]
			return
		}
	}

	displayResp(*resp)
	errors.New("headers don't have `Docker-Content-Digest` field")
	return

}

//根据打印列出指定数量的url
func ListRepositoriesPagination(opts ClientOpts, n int) ([]byte, error) {
	//	var respObject map[string][]string
	if n == 0 {
		opts.Url = opts.Url + "/v2/_catalog"
	} else {
		opts.Url = opts.Url + "/v2/_catalog?n=" + strconv.Itoa(n)
	}
	resp, err := doAction(opts, Get)
	if err != nil {
		debug.Print(err)
		return nil, err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		debug.Print(err)
		return nil, err
	}

	return byteContent, err
}

func ListImageTags(opts ClientOpts, image string) ([]byte, error) {
	//	var respObject map[string]interface{}
	opts.Url = opts.Url + "/v2/" + image + "/tags/list"
	debug.Print(opts.Url)
	resp, err := doAction(opts, Get)
	if err != nil {
		debug.Print(err)
		return nil, err
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		debug.Print(err)
		return nil, err
	}

	return byteContent, err
}

func GetImageManifests(opts ClientOpts, image string, tag string) (respObject Manifests, err error) {
	var resp *http.Response
	opts.Url = opts.Url + "/v2/" + image + "/manifests/" + tag

	resp, err = doAction(opts, Get)
	if err != nil {
		debug.Print(err)
		return
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	if resp.StatusCode != 200 {
		err = newApiError(resp, opts.Url)
		return
	}
	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteContent, interface{}(respObject))
	return
}

func DeleteImage(opts ClientOpts, image string, tag string) error {
	defaultOpts := opts
	docker_content_digest, err := GetImageDigest(opts, image, tag)
	if err != nil {
		return err
	}
	debug.Print(docker_content_digest)

	//			digest = strings.TrimLeft(v[0], "sha256:")
	opts = defaultOpts
	opts.Url = opts.Url + "/v2/" + image + "/manifests/" + docker_content_digest
	resp, err := doAction(opts, Delete)
	if err != nil {
		return err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != 202 {
		return newApiError(resp, opts.Url)
	}
	return nil
}

func displayResp(resp http.Response) {
	if debug.Debug {
		debug.Print("Header <===")
		for k, v := range resp.Header {
			debug.Print(k, v)
		}
		byteContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			debug.Print("read resp fail")
		}
		debug.Print("Respond <=" + string(byteContent))
	}

}
