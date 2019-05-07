package client

import (
	"fmt"
	"strings"
	"time"

	"encoding/json"
	"github.com/parnurzeal/gorequest"
)

type (
	request struct {
		remote     string
		url        string
		SuperAgent *gorequest.SuperAgent
		remotes    map[string][]string
		timeOut    time.Duration
		param      interface{}
		header     map[string]string
		Response   IResponse
	}
)

var (
	t       *Utils
	TimeOut time.Duration = 3
)

func init() {
	t = NewUtil()
}

func NewRequest(remote map[string][]string, debug bool, response IResponse) *request {
	return &request{
		remotes:    remote,
		Response:   response,
		SuperAgent: gorequest.New().SetDebug(debug),
	}
}

func (svc *request) SetTimeOut(times time.Duration) {
	svc.timeOut = times
}

func (svc *request) GetTimeOut() time.Duration {
	if svc.timeOut <= 0 {
		svc.timeOut = TimeOut
	}

	return svc.timeOut * time.Second
}

func (svc *request) SetRemote(remote string) *request {
	remoteArray, ok := svc.remotes[remote]
	if ok {
		max := len(remoteArray)
		num := t.GenerateRangeNum(0, max)
		svc.remote = remoteArray[num]
	}

	return svc
}

func (svc *request) SetPath(path string) *request {
	path = strings.TrimRight(strings.TrimLeft(path, "/"), "/")

	svc.url = fmt.Sprintf("%s/%s",
		strings.Trim(svc.remote, "/"),
		path)

	return svc
}

func (svc *request) SetHeader(data map[string]string) *request {
	svc.header = data

	return svc
}

func (svc *request) SetParam(param interface{}) *request {
	svc.param = param

	return svc
}

func (svc *request) GetParam() *request {
	svc.SuperAgent.Send(svc.param)

	return svc
}

func (svc *request) GetHeader() *request {
	svc.SuperAgent.Header = svc.header

	return svc
}

func (svc *request) Get() IResponse {
	svc.SuperAgent.Timeout(r.GetTimeOut()).Get(r.url)
	svc.GetHeader()

	res, body, err := svc.SuperAgent.Query(svc.param).End()
	if err == nil {
		return svc.Response.NewResponse(body, res.StatusCode)
	}

	return svc.Response.NewResponse(body)
}

func (svc *request) Post() IResponse {
	svc.SuperAgent.Timeout(svc.GetTimeOut()).Post(svc.url)
	svc.GetHeader().GetParam()

	res, body, err := svc.SuperAgent.End()
	if err == nil {
		return svc.Response.NewResponse(body, res.StatusCode)
	}

	return svc.Response.NewResponse(body)
}

func (svc *request) PostUrlEncode() IResponse {
	svc.SuperAgent.Timeout(svc.GetTimeOut()).Post(svc.url)
	svc.GetHeader().GetParam()

	res, body, err := svc.SuperAgent.End()
	if err == nil {
		return svc.Response.NewResponse(body, res.StatusCode)
	}

	return svc.Response.NewResponse(body)
}

func (svc *request) Put() IResponse {
	svc.SuperAgent.Timeout(svc.GetTimeOut()).Put(svc.url)
	svc.GetParam()

	res, body, err := svc.SuperAgent.End()
	if err == nil {
		return svc.Response.NewResponse(body, res.StatusCode)
	}

	return svc.Response.NewResponse(body)
}

func (svc *request) PostJson() IResponse {
	svc.SuperAgent.Timeout(svc.GetTimeOut()).Post(svc.url)

	paramsJson, errMsg := json.Marshal(svc.param)
	if errMsg != nil {
		return svc.Response.NewResponse("{}", 406)
	}

	res, body, err := svc.SuperAgent.Send(paramsJson).End()
	if err == nil {
		return svc.Response.NewResponse(body, res.StatusCode)
	}

	return svc.Response.NewResponse(body)
}

func (svc *request) Delete() IResponse {
	svc.SuperAgent.Timeout(svc.GetTimeOut()).Delete(svc.url)
	svc.GetParam()

	res, body, err := svc.SuperAgent.End()
	if err == nil {
		return svc.Response.NewResponse(body, res.StatusCode)
	}

	return svc.Response.NewResponse(body)
}
