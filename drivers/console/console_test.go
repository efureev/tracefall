package console

import (
	"github.com/efureev/traceFall"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConsoleDriver(t *testing.T) {

	Convey("Console Driver Tests", t, func() {

		db, err := traceFall.Open(`console`, GetDefaultConnParams())

		Convey("Open Instance", func() {
			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)
			So(db, ShouldHaveSameTypeAs, &traceFall.DB{})
			So(db.Driver(), ShouldHaveSameTypeAs, &DriverConsole{})
		})

		l := traceFall.NewLog(`Test Message`)

		Convey("Create Log", func() {
			So(l, ShouldHaveSameTypeAs, &traceFall.Log{})
		})

		Convey("Send Log", func() {
			resp, err := db.Send(l)
			So(err, ShouldBeNil)

			So(resp, ShouldHaveSameTypeAs, traceFall.Response{})
			So(resp.ID, ShouldNotBeNil)
			So(resp.ID, ShouldNotBeNil)
			So(resp.ID, ShouldHaveSameTypeAs, *new(string))
			So(resp.Data, ShouldHaveSameTypeAs, traceFall.ResponseData{})
			So(resp.Request(), ShouldEqual, l.String())
		})

		Convey("Truncate", func() {
			resp, err := db.Truncate(`test`)

			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, traceFall.Response{})
			So(resp.Error, ShouldBeNil)
			So(resp.Result, ShouldBeTrue)
			So(resp.Request(), ShouldHaveSameTypeAs, *new(string))
			So(resp.Request(), ShouldEqual, `Method not worked on Console Driver.. Don't use it!`)
		})

		Convey("RemoveThread", func() {
			id, _ := uuid.NewV4()
			resp, err := db.RemoveThread(id)

			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, traceFall.Response{})
			So(resp.Error, ShouldBeNil)
			So(resp.Result, ShouldBeTrue)
			So(resp.Request(), ShouldHaveSameTypeAs, uuid.UUID{})
			So(resp.Request(), ShouldEqual, id)
			So(resp.Data[`result`].(bool), ShouldBeTrue)
		})
	})

}
