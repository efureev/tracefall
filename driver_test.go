package tracefall

import (
	"errors"
	"testing"

	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

type DriverTest struct{}

func (d DriverTest) Send(l *Log) (ResponseCmd, error) {
	r := l.String()
	return *NewResponse(r).SetID(generateUUID().String()).Success().ToCmd(), nil
}

func (d DriverTest) GetLog(id uuid.UUID) (ResponseLog, error) {
	return *NewResponse(id).SetError(errors.New(`method worked on TEST Driver`)).ToLog(nil), nil
}

func (d DriverTest) RemoveThread(id uuid.UUID) (ResponseCmd, error) {
	return *NewResponse(id).Success().ToCmd(), nil
}

func (d DriverTest) RemoveByTags(tags Tags) (ResponseCmd, error) {
	return *NewResponse(tags).Success().ToCmd(), nil
}

func (d DriverTest) GetThread(id uuid.UUID) (ResponseThread, error) {
	l1 := NewLog(`log 1`).ToLogJSON()
	l2 := NewLog(`log 2`).ToLogJSON()

	return *NewResponse(id).Success().ToThread(Thread{l1, l2}), nil
}

func (d DriverTest) Truncate(ind string) (ResponseCmd, error) {
	return *NewResponse(ind).SetError(errors.New(`method worked on TEST Driver`)).ToCmd(), nil
}

func (d DriverTest) Open(map[string]string) (interface{}, error) {
	return nil, nil
}

func TestDriver(t *testing.T) {

	Convey("Driver Register", t, func() {

		Register("test", &DriverTest{})

		//Register("test", nil)
		So(func() { Register("test", nil) }, ShouldPanic)
		So(func() { Register("test", &DriverTest{}) }, ShouldPanic)

		Convey("Driver Open", func() {
			params := map[string]string{`key`: `val`}
			db, err := Open(`test`, params)

			So(len(Drivers()), ShouldEqual, 1)
			So(Drivers(), ShouldResemble, []string{`test`})

			So(err, ShouldBeNil)
			So(db, ShouldHaveSameTypeAs, &DB{})

			dbFail, err := Open(`miss`, params)

			So(err, ShouldBeError)
			So(dbFail, ShouldBeNil)

			So(db.Driver(), ShouldEqual, &DriverTest{})

			Convey("Send", func() {
				l := NewLog(`test log`).ThreadFinish()
				r, err := db.Send(l)
				So(err, ShouldBeNil)
				So(r, ShouldHaveSameTypeAs, ResponseCmd{})
				So(r.Error, ShouldBeNil)
				So(r.Result, ShouldBeTrue)
				So(r.Request(), ShouldHaveSameTypeAs, *new(string))
				So(r.Request(), ShouldEqual, l.String())
			})

			Convey("Get", func() {
				id := generateUUID()
				r, _ := db.GetLog(id)

				So(err, ShouldBeError)
				So(r, ShouldHaveSameTypeAs, ResponseLog{})

				So(r.Error, ShouldBeError)
				So(r.Result, ShouldBeFalse)
				So(r.Request(), ShouldHaveSameTypeAs, uuid.UUID{})
				So(r.Request(), ShouldEqual, id)
				So(r.Log, ShouldHaveSameTypeAs, &LogJSON{})
				So(r.Log, ShouldBeNil)
			})

			Convey("Truncate", func() {
				r, err := db.Truncate(`test`)

				So(err, ShouldBeNil)
				So(r, ShouldHaveSameTypeAs, ResponseCmd{})
				So(r.Error, ShouldBeError)
				So(r.Result, ShouldBeFalse)
				So(r.Request(), ShouldHaveSameTypeAs, *new(string))
				So(r.Request(), ShouldEqual, `test`)
			})

			Convey("Get Thread", func() {
				id := generateUUID()
				r, err := db.GetThread(id)

				So(err, ShouldBeNil)
				So(r, ShouldHaveSameTypeAs, ResponseThread{})
				So(r.Error, ShouldBeNil)
				So(r.Result, ShouldBeTrue)
				So(r.Request(), ShouldHaveSameTypeAs, uuid.UUID{})
				So(r.Request(), ShouldEqual, id)
				So(r.Thread, ShouldHaveSameTypeAs, Thread{})
				So(len(r.Thread), ShouldEqual, 2)
			})

			Convey("Remove Thread", func() {
				id := generateUUID()
				r, err := db.RemoveThread(id)

				So(err, ShouldBeNil)
				So(r, ShouldHaveSameTypeAs, ResponseCmd{})
				So(r.Error, ShouldBeNil)
				So(r.Result, ShouldBeTrue)
				So(r.Request(), ShouldHaveSameTypeAs, uuid.UUID{})
				So(r.Request(), ShouldEqual, id)
			})

			Convey("Remove By Tags", func() {
				tags := Tags{`first`, `two`}
				r, err := db.RemoveByTags(tags)

				So(err, ShouldBeNil)
				So(r, ShouldHaveSameTypeAs, ResponseCmd{})
				So(r.Error, ShouldBeNil)
				So(r.Result, ShouldBeTrue)
				So(r.Request(), ShouldHaveSameTypeAs, Tags{})
				So(r.Request(), ShouldResemble, tags)
			})

		})

		unregisterAllDrivers()
	})

}
