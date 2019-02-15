package tracefall

import (
	"time"
)

type Responsible interface {
	Success() Responsible
	//Fail() Responsible
	//SetError(err error) Responsible
	//GenerateID() Responsible
	//Request() interface{}
}

//type ResponseData map[string]interface{}

type BaseResponse struct {
	ID      string
	Error   error
	Result  bool
	Time    time.Time
	request interface{}
}

type ResponseCmd struct {
	BaseResponse
}

type ResponseThread struct {
	BaseResponse
	Thread Thread
}

type ResponseLog struct {
	BaseResponse
	Log *LogJSON
}

func (r *BaseResponse) Success() *BaseResponse {
	r.Result = true
	return r
}

func (r *BaseResponse) SetError(err error) *BaseResponse {
	if err != nil {
		r.Result = false
		r.Error = err
	}

	return r
}

func (r *BaseResponse) SetID(id string) *BaseResponse {
	r.ID = id
	return r
}

func (r *BaseResponse) ToCmd() *ResponseCmd {
	return &ResponseCmd{*r}
}

func (r *BaseResponse) ToThread(thread Thread) *ResponseThread {
	return &ResponseThread{*r, thread}
}

func (r *BaseResponse) ToLog(log *LogJSON) *ResponseLog {
	return &ResponseLog{*r, log}
}

func (r *BaseResponse) GenerateID() *BaseResponse {
	r.ID = generateUUID().String()
	return r
}

func (r *BaseResponse) Request() interface{} {
	return r.request
}

func NewResponse(request interface{}) *BaseResponse {
	return (&BaseResponse{
		request: request,
		Time:    time.Now(),
	}).GenerateID()
}
