package console

import (
	"bitbucket.org/efureev/traceFall.2"
	"github.com/satori/go.uuid"
)

type DriverConsole struct{}

func (d DriverConsole) Send(l *traceFall.Log) (traceFall.Response, error) {
	r := d.toString(l)
	println(r)
	return *traceFall.NewResponse(r).GenerateId().Success(), nil
}

func (d DriverConsole) RemoveThread(id uuid.UUID) (traceFall.Response, error) {
	r := `method not worked on Console Driver.. Don't use it`
	println(r)
	return *traceFall.NewResponse(id).SetData(traceFall.ResponseData{`result`: true}).GenerateId().Success(), nil
}

func (d DriverConsole) RemoveByTags(_ traceFall.Tags) (traceFall.Response, error) {
	r := `Method not worked on Console Driver.. Don't use it!`
	println(r)
	return *traceFall.NewResponse(r).GenerateId().Success(), nil
}

func (d DriverConsole) Get(_ uuid.UUID) (traceFall.Response, error) {
	r := `Method not worked on Console Driver.. Don't use it!`
	println(r)
	return *traceFall.NewResponse(r).GenerateId().Success(), nil
}
func (d DriverConsole) GetThread(id uuid.UUID) (traceFall.Response, error) {
	r := `Method not worked on Console Driver.. Don't use it!`
	println(r)
	return *traceFall.NewResponse(r).GenerateId().Success(), nil
}

func (d DriverConsole) Truncate(_ string) (traceFall.Response, error) {
	r := `Method not worked on Console Driver.. Don't use it!`
	println(r)
	return *traceFall.NewResponse(r).GenerateId().Success(), nil
}

func (d DriverConsole) Open(map[string]string) (interface{}, error) {
	return nil, nil
}

func (d DriverConsole) toString(l *traceFall.Log) string {
	return l.String()
}

func init() {
	traceFall.Register("console", &DriverConsole{})
}

func GetDefaultConnParams() map[string]string {
	return make(map[string]string)
}
