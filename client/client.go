package client

import (
	"time"
)

type (
	Client struct {
		Debug   bool                `json:"debug" yaml:"debug"`
		Remotes map[string][]string `json:"remotes" yaml:"remotes"`
		Resp    IResponse
	}
)

var (
	c *Client
	r *request
)

// init client
func NewClient(client *Client) *Client {
	c = client

	return c
}

// get client
func GetClient() *Client {
	getRequest()

	return c
}

func getRequest(res ...IResponse) *request {
	var rsp IResponse

	rsp = &Response{}
	if len(res) == 1 {
		rsp = res[0]
	}

	if r == nil {
		r = NewRequest(c.Remotes, c.Debug, rsp)
	}

	return r
}

//request get
func (c *Client) Get(
	remote,
	path string,
	queryParams interface{}) IResponse {

	return getRequest(c.Resp).SetRemote(remote).SetPath(path).SetParam(queryParams).Get()
}

//request post
func (c *Client) Post(
	remote,
	path string,
	dataForm interface{}) IResponse {

	return getRequest(c.Resp).SetRemote(remote).SetPath(path).SetParam(dataForm).Post()
}

//request post
func (c *Client) PostUrlEncode(
	remote,
	path string,
	dataForm interface{}) IResponse {

	return getRequest(c.Resp).SetRemote(remote).SetPath(path).SetParam(dataForm).PostUrlEncode()
}

//request put
func (c *Client) Put(
	remote,
	path string,
	dataForm interface{}) IResponse {

	return getRequest(c.Resp).SetRemote(remote).SetPath(path).SetParam(dataForm).Put()
}

//request post json
func (c *Client) PostJson(
	remote,
	path string,
	dataJson interface{}) IResponse {

	return getRequest(c.Resp).SetRemote(remote).SetPath(path).SetParam(dataJson).PostJson()
}

func (c *Client) Delete(
	remote,
	path string,
	dataForm interface{}) IResponse {

	return getRequest(c.Resp).SetRemote(remote).SetPath(path).SetParam(dataForm).Delete()
}

// set request header data params
func (c *Client) SetHeader(data map[string]string) *Client {
	getRequest(c.Resp).SetHeader(data)

	return c
}

// set request time out params
func (c *Client) SetTimeOut(times time.Duration) *Client {
	getRequest(c.Resp).SetTimeOut(times)

	return c
}

func (c *Client) AddParams(key, value string) {
	getRequest(c.Resp).SuperAgent.QueryData.Add(key, value)
}
