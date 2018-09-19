package traceFall

import (
	"time"
)

type ResponseData map[string]interface{}
type Response struct {
	Id      string
	Error   error
	Data    ResponseData
	Result  bool
	Time    time.Time
	request interface{}
}

func NewResponse(request interface{}) *Response {
	return &Response{
		request: request,
		Data:    make(ResponseData),
		Time:    time.Now(),
	}
}

func (r *Response) Success() *Response {
	r.Result = true
	return r
}

func (r *Response) SetError(err error) *Response {
	if err != nil {
		r.Result = false
		r.Error = err
	}

	return r
}
func (r *Response) SetId(id string) *Response {
	r.Id = id
	return r
}

func (r *Response) SetData(data ResponseData) *Response {
	r.Data = data
	return r
}

func (r *Response) GenerateId() *Response {
	r.Id = generateUUID().String()
	return r
}

func (r *Response) Request() interface{} {
	return r.request
}
