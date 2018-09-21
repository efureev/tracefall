package console

import (
	"github.com/efureev/traceFall"
	"github.com/satori/go.uuid"
)

type DriverConsole struct{}

func (d DriverConsole) Send(l *traceFall.Log) (traceFall.ResponseCmd, error) {
	r := d.toString(l)
	println(r)
	return *traceFall.NewResponse(r).Success().ToCmd(), nil
}

func (d DriverConsole) RemoveThread(id uuid.UUID) (traceFall.ResponseCmd, error) {
	r := `method not worked on Console Driver.. Don't use it`
	println(r)
	return *traceFall.NewResponse(id).Success().ToCmd(), nil
}

func (d DriverConsole) RemoveByTags(_ traceFall.Tags) (traceFall.ResponseCmd, error) {
	r := `Method not worked on Console Driver.. Don't use it!`
	println(r)
	return *traceFall.NewResponse(r).Success().ToCmd(), nil
}

func (d DriverConsole) GetLog(_ uuid.UUID) (traceFall.ResponseLog, error) {
	r := `Method not worked on Console Driver.. Don't use it!`
	println(r)
	return *traceFall.NewResponse(r).Success().ToLog(traceFall.NewLog(`console`).ToLogJSON()), nil
}
func (d DriverConsole) GetThread(id uuid.UUID) (traceFall.ResponseThread, error) {
	r := `Method not worked on Console Driver.. Don't use it!`
	println(r)
	return *traceFall.NewResponse(r).Success().ToThread(traceFall.Thread{}), nil
}

func (d DriverConsole) Truncate(_ string) (traceFall.ResponseCmd, error) {
	r := `Method not worked on Console Driver.. Don't use it!`
	println(r)
	return *traceFall.NewResponse(r).Success().ToCmd(), nil
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
