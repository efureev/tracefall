package traceFall

import (
	"errors"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type DriverTest struct{}

func (d DriverTest) Send(l *Log) (Response, error) {
	r := l.String()
	return *NewResponse(r).GenerateID().Success(), nil
}

func (d DriverTest) Get(id uuid.UUID) (Response, error) {
	return *NewResponse(id).SetError(errors.New(`method worked on TEST Driver`)).GenerateID(), nil
}

func (d DriverTest) RemoveThread(id uuid.UUID) (Response, error) {
	return *NewResponse(id).SetData(ResponseData{`result`: true}).GenerateID().Success(), nil
}

func (d DriverTest) RemoveByTags(tags Tags) (Response, error) {
	return *NewResponse(tags).SetData(ResponseData{`result`: true}).GenerateID().Success(), nil
}

func (d DriverTest) GetThread(id uuid.UUID) (Response, error) {
	l1 := NewLog(`log 1`).ToLogJSON()
	l2 := NewLog(`log 2`).ToLogJSON()
	data := ResponseData{`thread`: Thread{&l1, &l2}}
	return *NewResponse(id).SetData(data).GenerateID().Success(), nil
}

func (d DriverTest) Truncate(ind string) (Response, error) {
	return *NewResponse(ind).SetError(errors.New(`method worked on TEST Driver`)).GenerateID(), nil
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
				So(r, ShouldHaveSameTypeAs, Response{})
				So(r.Error, ShouldBeNil)
				So(r.Result, ShouldBeTrue)
				So(r.Request(), ShouldHaveSameTypeAs, *new(string))
				So(r.Request(), ShouldEqual, l.String())
			})

			Convey("Get", func() {
				id := generateUUID()
				r, err := db.Get(id)
				So(err, ShouldBeNil)
				So(r, ShouldHaveSameTypeAs, Response{})
				So(r.Error, ShouldBeError)
				So(r.Result, ShouldBeFalse)
				So(r.Request(), ShouldHaveSameTypeAs, uuid.UUID{})
				So(r.Request(), ShouldEqual, id)
			})

			Convey("Truncate", func() {
				r, err := db.Truncate(`test`)

				So(err, ShouldBeNil)
				So(r, ShouldHaveSameTypeAs, Response{})
				So(r.Error, ShouldBeError)
				So(r.Result, ShouldBeFalse)
				So(r.Request(), ShouldHaveSameTypeAs, *new(string))
				So(r.Request(), ShouldEqual, `test`)
			})

			Convey("Get Thread", func() {
				id := generateUUID()
				r, err := db.GetThread(id)

				So(err, ShouldBeNil)
				So(r, ShouldHaveSameTypeAs, Response{})
				So(r.Error, ShouldBeNil)
				So(r.Result, ShouldBeTrue)
				So(r.Request(), ShouldHaveSameTypeAs, uuid.UUID{})
				So(r.Request(), ShouldEqual, id)
				So(r.Data[`thread`], ShouldHaveSameTypeAs, Thread{})
				So(len(r.Data[`thread`].(Thread)), ShouldEqual, 2)
			})

			Convey("Remove Thread", func() {
				id := generateUUID()
				r, err := db.RemoveThread(id)

				So(err, ShouldBeNil)
				So(r, ShouldHaveSameTypeAs, Response{})
				So(r.Error, ShouldBeNil)
				So(r.Result, ShouldBeTrue)
				So(r.Request(), ShouldHaveSameTypeAs, uuid.UUID{})
				So(r.Request(), ShouldEqual, id)
				So(r.Data[`result`], ShouldHaveSameTypeAs, *new(bool))
				So(r.Data[`result`].(bool), ShouldBeTrue)
			})

			Convey("Remove RemoveByTags", func() {
				tags := Tags{`first`, `two`}
				r, err := db.RemoveByTags(tags)

				So(err, ShouldBeNil)
				So(r, ShouldHaveSameTypeAs, Response{})
				So(r.Error, ShouldBeNil)
				So(r.Result, ShouldBeTrue)
				So(r.Request(), ShouldHaveSameTypeAs, Tags{})
				So(r.Request(), ShouldResemble, tags)
				So(r.Data[`result`], ShouldHaveSameTypeAs, *new(bool))
				So(r.Data[`result`].(bool), ShouldBeTrue)
			})

		})

		unregisterAllDrivers()
	})

}
