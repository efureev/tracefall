package postgres

import (
	"github.com/efureev/traceFall"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPostgresDriverOpen(t *testing.T) {
	Convey("Postgres Driver Tests", t, func() {
		Convey("Open wrong db", func() {
			So(func() {
				params := GetConnParams("localhost:5432", "postgres", "tracerFake", `root`, ``)
				traceFall.Open(`postgres`, params)
			}, ShouldPanic)

			So(func() {
				params := GetConnParams("localhost:54321", "postgres", "tracer", `postgres`, `postgres`)
				traceFall.Open(`postgres`, params)
			}, ShouldPanic)
		})

		Convey("Open Instance", func() {
			params := GetConnParams("localhost:5432", "postgres", "tracer", `postgres`, `postgres`)
			db, err := traceFall.Open(`postgres`, params)

			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)
			So(db, ShouldHaveSameTypeAs, &traceFall.DB{})
			So(db.Driver(), ShouldHaveSameTypeAs, &DriverPostgres{})

			db1 := db.Driver().(*DriverPostgres)
			err = db1.DropTable()
			So(err, ShouldBeNil)

			err = db1.CreateTable()
			So(err, ShouldBeNil)

			resp, err := db1.Truncate(``)

			So(err, ShouldBeNil)
			So(resp, ShouldHaveSameTypeAs, traceFall.ResponseCmd{})
			So(resp.Error, ShouldBeNil)
			So(resp.Result, ShouldBeTrue)
			So(resp.Request(), ShouldHaveSameTypeAs, *new(string))
			So(resp.Request(), ShouldEqual, ``)

			respFail, err := db1.Truncate(`absent`)

			So(err, ShouldBeError)
			So(respFail, ShouldHaveSameTypeAs, traceFall.ResponseCmd{})
			So(respFail.Error, ShouldBeError)
			So(respFail.Result, ShouldBeFalse)
			So(respFail.Request(), ShouldHaveSameTypeAs, *new(string))
			So(respFail.Request(), ShouldEqual, `absent`)

		})

	})
}

func TestPostgresDriverCreateAndDrop(t *testing.T) {
	params := GetConnParams("localhost:5432", "postgres", "tracer", `postgres`, `postgres`)
	db, err := traceFall.Open(`postgres`, params)

	Convey("Send & Drop", t, func() {
		So(err, ShouldBeNil)
		So(db, ShouldNotBeNil)
		So(db, ShouldHaveSameTypeAs, &traceFall.DB{})
		So(db.Driver(), ShouldHaveSameTypeAs, &DriverPostgres{})

		l := traceFall.NewLog(`Test`)

		resp, err := db.Send(l)

		So(err, ShouldBeNil)
		So(resp, ShouldHaveSameTypeAs, traceFall.ResponseCmd{})
		So(resp.ID, ShouldNotBeNil)
		So(resp.ID, ShouldHaveSameTypeAs, *new(string))
		So(resp.ID, ShouldEqual, l.ID.String())
		So(resp.Result, ShouldBeTrue)

		l2, err := l.CreateChild(`Test2`)
		So(err, ShouldBeNil)

		l2.SetEnvironment(`prod`).Success().ThreadFinish().Tags.Add(`child`)
		l2.Notes.Add(`step`, `note1`).
			Add(`step`, `note2`).
			Add(`id`, l.ID.String())
		l2.Data.Set(`thread`, `thread: `+l.Thread.String())

		resp2, err := db.Send(l2)

		So(err, ShouldBeNil)
		So(resp2, ShouldHaveSameTypeAs, traceFall.ResponseCmd{})
		So(resp2.ID, ShouldNotBeNil)
		So(resp2.ID, ShouldHaveSameTypeAs, *new(string))
		So(resp2.ID, ShouldEqual, l2.ID.String())
		So(resp2.Result, ShouldBeTrue)

		lGet, err := db.Get(l.ID)
		So(err, ShouldBeNil)
		So(lGet, ShouldHaveSameTypeAs, traceFall.ResponseLog{})
		So(lGet.Log, ShouldHaveSameTypeAs, &traceFall.LogJSON{})
		So(lGet.ID, ShouldNotBeNil)
		So(lGet.Result, ShouldBeTrue)
		So(lGet.Error, ShouldBeNil)

		So(lGet.Log.ID.String(), ShouldEqual, l.ID.String())
		So(l.Environment, ShouldEqual, traceFall.EnvironmentDev)

		// fail
		uid, _ := uuid.NewV4()
		lGetFail, err2 := db.Get(uid)
		So(err2, ShouldBeError)
		So(lGetFail.Error, ShouldBeError)
		So(lGetFail, ShouldHaveSameTypeAs, traceFall.ResponseLog{})
		So(lGetFail.Result, ShouldBeFalse)
		So(lGetFail.Log, ShouldBeNil)

		resp3, err := db.RemoveByTags([]string{`child`})
		So(err, ShouldBeNil)
		So(resp3, ShouldHaveSameTypeAs, traceFall.ResponseCmd{})
		So(resp3.Result, ShouldBeTrue)
		So(resp3.Error, ShouldBeNil)

		So(err, ShouldBeNil)
		So(resp3, ShouldHaveSameTypeAs, traceFall.ResponseCmd{})
		So(resp3.Result, ShouldBeTrue)
		So(resp3.Error, ShouldBeNil)
	})

}

func TestPostgresDriverGetter(t *testing.T) {

	Convey("Postgres Driver Getter", t, func() {

		params := GetConnParams("localhost:5432", "postgres", "tracer", `postgres`, `postgres`)

		Convey("Create Params", func() {
			So(params, ShouldHaveSameTypeAs, map[string]string{})
		})

		db, err := traceFall.Open(`postgres`, params)

		Convey("Open Instance", func() {
			So(err, ShouldBeNil)
			So(db, ShouldNotBeNil)
			So(db, ShouldHaveSameTypeAs, &traceFall.DB{})
			So(db.Driver(), ShouldHaveSameTypeAs, &DriverPostgres{})
		})

		l := traceFall.NewLog(`Test Root`)
		Convey("Create Log", func() {
			So(l, ShouldHaveSameTypeAs, &traceFall.Log{})
		})

		Convey("Send Log", func() {
			resp, err := db.Send(l)
			So(err, ShouldBeNil)

			So(resp, ShouldHaveSameTypeAs, traceFall.ResponseCmd{})
			So(resp.ID, ShouldNotBeNil)
			So(resp.ID, ShouldHaveSameTypeAs, *new(string))

			So(resp.Error, ShouldBeNil)

			So(resp.ID, ShouldEqual, l.ID.String())
			So(resp.Result, ShouldBeTrue)

			l2, err := l.CreateChild(`Child`)
			So(err, ShouldBeNil)
			l2.SetEnvironment(`prod`).Success().ThreadFinish().Tags.Add(`child`)
			l2.Notes.Add(`step`, `note1`).Add(`step`, `note2`)
			l2.Data.Set(`id`, l2.ID).
				Set(`thread`, `thread: `+l2.Thread.String())

			resp2, err := db.Send(l2)
			So(err, ShouldBeNil)
			So(resp2.Error, ShouldBeNil)
			So(resp2.ID, ShouldEqual, l2.ID.String())
			So(resp2.Result, ShouldBeTrue)

			l3, err := l.CreateChild(`Child`)
			So(err, ShouldBeNil)

			l3.Success().Tags.
				Add(`child`).Add(`2`)
			l3.Data.Set(`step`, `note1`).Set(`id`, l3.ID)

			resp3, err := db.Send(l3)

			So(err, ShouldBeNil)
			So(resp3.Error, ShouldBeNil)
			So(resp3.ID, ShouldEqual, l3.ID.String())
			So(resp3.Result, ShouldBeTrue)

			l4, err := l3.CreateChild(`SubChild of Child`)
			So(err, ShouldBeNil)

			l4.Success().ThreadFinish().Tags.Add(`child`)

			resp4, err := db.Send(l4)

			So(err, ShouldBeNil)
			So(resp4.Error, ShouldBeNil)
			So(resp4.ID, ShouldEqual, l4.ID.String())
			So(resp4.Result, ShouldBeTrue)

			/// get
			Convey("Get Log", func() {

				lGet, err := db.Get(l.ID)
				So(err, ShouldBeNil)
				logRootGet := lGet.Log

				So(logRootGet.ID.String(), ShouldEqual, l.ID.String())
				So(logRootGet.Parent, ShouldBeNil)
				So(logRootGet.Time, ShouldEqual, l.Time.UnixNano())
				So(logRootGet.TimeEnd, ShouldBeNil)
				So(logRootGet.Error, ShouldEqual, l.Error)
				So(logRootGet.App, ShouldEqual, l.App)
				So(logRootGet.Name, ShouldEqual, l.Name)
				So(logRootGet.Environment, ShouldEqual, l.Environment)
				So(logRootGet.ID.String(), ShouldEqual, l.ID.String())
				So(logRootGet.Finish, ShouldEqual, l.Finish)
				So(logRootGet.Result, ShouldEqual, l.Result)
				So(logRootGet.Tags, ShouldResemble, l.Tags.List())

				l2Get, err := db.Get(l2.ID)
				So(err, ShouldBeNil)

				log2Get := l2Get.Log

				So(log2Get.ID.String(), ShouldEqual, l2.ID.String())
				So(log2Get.Parent, ShouldNotBeNil)
				So(*log2Get.Parent, ShouldEqual, logRootGet.ID.String())
				So(log2Get.Time, ShouldEqual, l2.Time.UnixNano())
				So(*log2Get.TimeEnd, ShouldEqual, l2.TimeEnd.UnixNano())
				So(log2Get.App, ShouldEqual, l2.App)
				So(log2Get.Name, ShouldEqual, l2.Name)
				So(log2Get.Environment, ShouldEqual, l2.Environment)
				So(log2Get.Finish, ShouldEqual, l2.Finish)
				So(log2Get.Result, ShouldEqual, l2.Result)
				So(log2Get.Tags, ShouldResemble, l2.Tags.List())

				l3Get, err := db.Get(l3.ID)
				So(err, ShouldBeNil)

				log3Get := l3Get.Log

				So(log3Get.ID.String(), ShouldEqual, l3.ID.String())
				So(log3Get.Parent, ShouldNotBeNil)
				So(*log3Get.Parent, ShouldEqual, logRootGet.ID.String())
				So(log3Get.Time, ShouldEqual, l3.Time.UnixNano())
				So(*log3Get.TimeEnd, ShouldEqual, l3.TimeEnd.UnixNano())
				So(log3Get.Error, ShouldBeNil)
				So(log3Get.App, ShouldEqual, l3.App)
				So(log3Get.Name, ShouldEqual, l3.Name)
				So(log3Get.Environment, ShouldEqual, l3.Environment)
				So(log3Get.Finish, ShouldEqual, l3.Finish)
				So(log3Get.Result, ShouldEqual, l3.Result)
				So(log3Get.Tags, ShouldResemble, l3.Tags.List())

				l4Get, err := db.Get(l4.ID)
				So(err, ShouldBeNil)
				log4Get := l4Get.Log

				So(log4Get.ID.String(), ShouldEqual, l4.ID.String())
				So(log4Get.Parent, ShouldNotBeNil)
				So(*log4Get.Parent, ShouldEqual, log3Get.ID.String())

				So(log4Get.Time, ShouldEqual, l4.Time.UnixNano())
				So(*log4Get.TimeEnd, ShouldEqual, l4.TimeEnd.UnixNano())
				So(log4Get.Error, ShouldBeNil)
				So(log4Get.App, ShouldEqual, l4.App)
				So(log4Get.Name, ShouldEqual, l4.Name)
				So(log4Get.Environment, ShouldEqual, l4.Environment)
				So(log4Get.Finish, ShouldEqual, l4.Finish)
				So(log4Get.Result, ShouldEqual, l4.Result)
				So(log4Get.Tags, ShouldResemble, l4.Tags.List())

				lThreadResp, err := db.GetThread(l4.Thread)
				So(err, ShouldBeNil)
				So(lThreadResp.Error, ShouldBeNil)

				lThread := lThreadResp.Thread
				So(len(lThread), ShouldEqual, 4)

				for _, v := range lThread {
					So(v, ShouldHaveSameTypeAs, &traceFall.LogJSON{})
				}


				respRemove, err := db.RemoveThread(l4.Thread)
				So(err, ShouldBeNil)
				So(respRemove.Error, ShouldBeNil)
				So(respRemove.Result, ShouldBeTrue)
				So(respRemove.Request(), ShouldEqual, l4.Thread)
			})
		})
	})
}
