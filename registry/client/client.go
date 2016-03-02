package client

//package main

import (
	"fmt"
	"net/http"
	//	"net/url" 需要使用url解析字符串,判定url是否合法
	"bytes"
	"encoding/json"
	//	"errors"
	"io/ioutil"
	"strconv"
	//"strings"
	"time"
)

//默认超时时间
const (
	DefaultTimeOut time.Duration = time.Second * 10
)

var (
	debug = false
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

func doGet(opts ClientOpts, respObject interface{}) (err error) {
	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeOut
	}
	client := &http.Client{Timeout: opts.Timeout}

	req, err := http.NewRequest("GET", opts.Url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(opts.AccessKey, opts.SecretKey)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return newApiError(resp, opts.Url)
	}

	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if debug {
		fmt.Println("Respond <=" + string(byteContent))
	}
	return json.Unmarshal(byteContent, respObject)
	//return nil
}

func doGet2(opts ClientOpts) (resp *http.Response, err error) {
	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeOut
	}
	client := &http.Client{Timeout: opts.Timeout}

	req, err := http.NewRequest("GET", opts.Url, nil)
	if err != nil {
		return
	}

	req.SetBasicAuth(opts.AccessKey, opts.SecretKey)
	resp, err = client.Do(req)

	return

	//return nil
}

func doHeader(opts ClientOpts) (http.Header, error) {
	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeOut
	}
	client := &http.Client{Timeout: opts.Timeout}

	req, err := http.NewRequest("HEAD", opts.Url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(opts.AccessKey, opts.SecretKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, newApiError(resp, opts.Url)
	}
	return resp.Header, nil
}

func doDelete(opts ClientOpts) error {
	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeOut
	}
	client := &http.Client{Timeout: opts.Timeout}

	req, err := http.NewRequest("DELETE", opts.Url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(opts.AccessKey, opts.SecretKey)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return newApiError(resp, opts.Url)
	}
	return nil
}

/*检测registry Api版本*/
func CheckVersion(opts ClientOpts) error {
	var respObject map[string][]string
	opts.Url = opts.Url + "/v2"
	err := doGet(opts, &respObject)
	if err != nil {
		//检测是否http请求触发的错误.
		x, ok := interface{}(err).(ApiError)
		if ok {
			switch x.StatusCode {
			default:
				fmt.Println(err)
			case 401:
				//do something
				fmt.Println(":未授权")
			case 404:
				fmt.Println(err)
			}
		}
	}
	return err
}

/*该url禁用了header请求*/
/*
func GetImageDigest(opts ClientOpts, image string, tag string) (string, error) {
	var digest string

	opts.Url = opts.Url + "/" + "/v2/" + image + "/manifests/" + tag
	header, err := doGetHeader(opts)
	if err != nil {
		return digest, err
	}
	for k, v := range header {
		if k == "Docker-Content-Digest" {
			digest = strings.TrimLeft(v[0], "sha256:")
			break
		}
	}

	if digest == "" {
		return digest, errors.New("headers don't have `Docker-Content-Digest` field")
	}
	return digest, nil

}
*/

//根据打印列出指定数量的url
func ListRepositoriesPagination(opts ClientOpts, n int) ([]byte, error) {
	//	var respObject map[string][]string
	if n == 0 {
		opts.Url = opts.Url + "/v2/_catalog"
	} else {
		opts.Url = opts.Url + "/v2/_catalog?n=" + strconv.Itoa(n)
	}
	resp, err := doGet2(opts)
	if err != nil {
		debug_print(err)
		return nil, err
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		debug_print(err)
		return nil, err
	}

	return byteContent, err
}

func ListImageTags(opts ClientOpts, image string) ([]byte, error) {
	//	var respObject map[string]interface{}
	opts.Url = opts.Url + "/v2/" + image + "/tags/list"
	debug_print(opts.Url)
	resp, err := doGet2(opts)
	if err != nil {
		debug_print(err)
		return nil, err
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	byteContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		debug_print(err)
		return nil, err
	}

	return byteContent, err
}

func GetImageManifests(opts ClientOpts, image string, tag string) (Manifests, error) {
	var respObject Manifests
	opts.Url = opts.Url + "/v2/" + image + "/manifests/" + tag
	err := doGet(opts, &respObject)
	return respObject, err
}

/*
func DeleteImage(opts ClientOpts, image string, tag string) error {
	defaultOpts := opts
	digest, err := GetImageDigest(opts, image, tag)
	if err != nil {
		return errors.New("can't get digest")
	}

	opts = defaultOpts
	opts.Url = opts.Url + "/v2/" + image + "/manifests/" + digest
	err = doDelete(opts)
	return err
}
*/

func debug_print(x interface{}) {
	if debug {
		fmt.Println(x)
	}
}

func init() {
	debug = true
}
