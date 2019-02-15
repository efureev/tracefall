package console

import (
	"github.com/efureev/tracefall"
	uuid "github.com/satori/go.uuid"
)

type DriverConsole struct{}

func (d DriverConsole) Send(l *tracefall.Log) (tracefall.ResponseCmd, error) {
	r := d.toString(l)
	println(r)
	return *tracefall.NewResponse(r).Success().ToCmd(), nil
}

func (d DriverConsole) RemoveThread(id uuid.UUID) (tracefall.ResponseCmd, error) {
	r := `method not worked on Console Driver.. Don't use it`
	println(r)
	return *tracefall.NewResponse(id).Success().ToCmd(), nil
}

func (d DriverConsole) RemoveByTags(tags tracefall.Tags) (tracefall.ResponseCmd, error) {
	r := `Method not worked on Console Driver.. Don't use it!`
	println(r)
	return *tracefall.NewResponse(tags).Success().ToCmd(), nil
}

func (d DriverConsole) GetLog(_ uuid.UUID) (tracefall.ResponseLog, error) {
	r := `Method not worked on Console Driver.. Don't use it!`
	println(r)
	return *tracefall.NewResponse(r).Success().ToLog(tracefall.NewLog(`console`).ToLogJSON()), nil
}
func (d DriverConsole) GetThread(id uuid.UUID) (tracefall.ResponseThread, error) {
	r := `Method not worked on Console Driver.. Don't use it!`
	println(r)
	return *tracefall.NewResponse(id).Success().ToThread(tracefall.Thread{}), nil
}

func (d DriverConsole) Truncate(id string) (tracefall.ResponseCmd, error) {
	r := `Method not worked on Console Driver.. Don't use it!`
	println(r)
	return *tracefall.NewResponse(r).Success().ToCmd(), nil
}

func (d DriverConsole) Open(map[string]string) (interface{}, error) {
	return nil, nil
}

func (d DriverConsole) toString(l *tracefall.Log) string {
	return l.String()
}

func init() {
	tracefall.Register("console", &DriverConsole{})
}

func GetDefaultConnParams() map[string]string {
	return make(map[string]string)
}
