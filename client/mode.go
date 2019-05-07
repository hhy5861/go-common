package client

type (
	IResponse interface {
		NewResponse(body string, code ...int) IResponse

		GetHttpStatus() int

		GetSuccess() int

		GetData() interface{}

		GetMessage() string

		GetStruct(data interface{}) error

		GetBody() string
	}
)
