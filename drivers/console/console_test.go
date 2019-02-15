package console

import (
	"testing"

	"github.com/efureev/tracefall"
	uuid "github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConsoleDriver(t *testing.T) {

	Convey("Console Driver Tests", t, func() {

		db, err := tracefall.Open(`console`, GetDefaultConnParams())

		Convey("Open Instance", func() {
			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)
			So(db, ShouldHaveSameTypeAs, &tracefall.DB{})
			So(db.Driver(), ShouldHaveSameTypeAs, &DriverConsole{})
		})

		l := tracefall.NewLog(`Test Message`)

		Convey("Create Log", func() {
			So(l, ShouldHaveSameTypeAs, &tracefall.Log{})
		})

		Convey("Send Log", func() {
			resp, err := db.Send(l)
			So(err, ShouldBeNil)

			So(resp, ShouldHaveSameTypeAs, tracefall.ResponseCmd{})
			So(resp.ID, ShouldNotBeNil)
			So(resp.ID, ShouldNotBeNil)
			So(resp.ID, ShouldHaveSameTypeAs, *new(string))
			So(resp.Result, ShouldBeTrue)
			So(resp.Request(), ShouldEqual, l.String())
		})

		Convey("Get Log", func() {
			uid, _ := uuid.NewV4()
			resp, err := db.GetLog(uid)
			So(resp.Result, ShouldBeTrue)
			So(resp.Log.Name, ShouldEqual, `console`)
			So(resp.Error, ShouldBeNil)
			So(err, ShouldBeNil)
		})

		Convey("Truncate", func() {
			resp, err := db.Truncate(`test`)

			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, tracefall.ResponseCmd{})
			So(resp.Error, ShouldBeNil)
			So(resp.Result, ShouldBeTrue)
			So(resp.Request(), ShouldHaveSameTypeAs, *new(string))
			So(resp.Request(), ShouldEqual, `Method not worked on Console Driver.. Don't use it!`)
		})

		Convey("Remove Thread", func() {
			id, _ := uuid.NewV4()
			resp, err := db.RemoveThread(id)

			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, tracefall.ResponseCmd{})
			So(resp.Error, ShouldBeNil)
			So(resp.Result, ShouldBeTrue)
			So(resp.Request(), ShouldHaveSameTypeAs, uuid.UUID{})
			So(resp.Request(), ShouldEqual, id)
		})

		Convey("Get Thread", func() {
			id, _ := uuid.NewV4()
			resp, err := db.GetThread(id)

			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, tracefall.ResponseThread{})
			So(resp.Error, ShouldBeNil)
			So(resp.Result, ShouldBeTrue)
			So(resp.Request(), ShouldHaveSameTypeAs, uuid.UUID{})
			So(resp.Thread, ShouldHaveSameTypeAs, tracefall.Thread{})
			So(resp.Request(), ShouldEqual, id)
		})

		Convey("Remove By Tags", func() {
			resp, err := db.RemoveByTags(tracefall.Tags{`tag 1`})

			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, tracefall.ResponseCmd{})
			So(resp.Error, ShouldBeNil)
			So(resp.Result, ShouldBeTrue)
			So(resp.Request(), ShouldHaveSameTypeAs, tracefall.Tags{})
		})
	})

}
