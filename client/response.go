package client

import (
	"encoding/json"
	"fmt"
	"github.com/hhy5861/go-common/logger"
	"net/http"
)

type (
	Response struct {
		Success    int         `json:"success"`
		Data       interface{} `json:"data"`
		Message    string      `json:"message"`
		body       string
		httpStatus int
	}
)

func (svc *Response) NewResponse(body string, code ...int) IResponse {
	status := http.StatusBadGateway
	if len(code) == 1 {
		status = code[0]
	}

	svc.body = body
	svc.httpStatus = status

	err := json.Unmarshal([]byte(body), &svc)
	fmt.Println("err", err)
	if err != nil {
		logger.Error(err)
	}

	return svc
}

func (svc *Response) GetHttpStatus() int {
	return svc.httpStatus
}

func (svc *Response) GetSuccess() int {
	return svc.Success
}

func (svc *Response) GetData() interface{} {
	return svc.Data
}

func (svc *Response) GetMessage() string {
	return svc.Message
}

func (svc *Response) GetStruct(data interface{}) error {
	return json.Unmarshal([]byte(svc.body), &data)
}

func (svc *Response) GetBody() string {
	return svc.body
}
