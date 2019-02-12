package main

import (
	"errors"
	"fmt"
	"github.com/efureev/traceFall"
	"github.com/efureev/traceFall/drivers/postgres"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var X *Hub

type Hub struct {
	ToTraceLog chan *traceFall.Log
}

func init() {
	tracerLogStart(`localhost`, `efureev`, ``, `test`, `tracer`)
	newHub()
}

func newHub() *Hub {
	X = &Hub{
		ToTraceLog: make(chan *traceFall.Log),
	}
	return X
}

var logStorage *traceFall.DB

func tracerLogStart(tracerHost, tracerUser, tracerPassword, tracerDbName, tracerTable string) {
	var err error

	logStorage, err = traceFall.Open(`postgres`, postgres.GetConnParams(tracerHost, tracerDbName, tracerTable, tracerUser, tracerPassword))
	if err != nil {
		panic(err)
	}
}

func listen() {
	go func() {
		for {
			select {
			case m := <-X.ToTraceLog:
				println(`pull msg: ` + m.ID.String())

				if logStorage == nil {
					break
				}

				_, err := logStorage.Send(m)
				if err != nil {
					fmt.Println(`[error sent to trace logs] -> ` + err.Error())
				}
			}
		}
	}()
}

func runWork() {
	logParent := traceFall.NewLog(`Start`).SetApplication(`micro#1`)
	logParent.Success().Data.Set(`key1`, `zvalue`)
	logParent.Tags.Add(`micro1`).Add(`root`)

	X.ToTraceLog <- logParent

	for i := 0; i < 3; i++ {
		logChildren, err := logParent.CreateChild(fmt.Sprintf(`Processing # %d`, i))
		if err != nil {
			panic(err)
		}

		if i%2 == 0 {
			logChildren.Fail(errors.New(`error in child`))
		} else {
			logChildren.Success()
		}

		logChildren.Notes.
			AddGroup(`proc 1`, []string{`step one`, `step two`}).
			Add(`proc 1`, `step three`).
			Add(`proc 2`, `finally`)

		X.ToTraceLog <- logChildren
		shadow := logChildren.ToShadow()

		// new log form other service:
		logOther := traceFall.NewLog(`Resulting`).SetApplication(`micro#2`)
		logOther.Tags.Add(`micro2`).Add(`finish`)
		logOther.ParentFromShadow(shadow).Success().ThreadFinish()
		X.ToTraceLog <- logOther
	}
}

func main() {
	listen()

	println(`started... wait work for every 10 seconds`)

	go func() {
		for {
			select {
			case <-time.After(10 * time.Second):
				println(`run Work`)
				runWork()
			}
		}

	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	select {
	// wait on kill signal
	case <-exit:
	}

	println(`exiting...`)
}
