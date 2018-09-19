package postgres

import (
	"github.com/efureev/traceFall"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
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

		params := GetConnParams("localhost:5432", "postgres", "tracer", `postgres`, `postgres`)
		db, err := traceFall.Open(`postgres`, params)

		Convey("Open Instance", func() {
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
			So(resp, ShouldHaveSameTypeAs, traceFall.Response{})
			So(resp.Error, ShouldBeNil)
			So(resp.Result, ShouldBeTrue)
			So(resp.Request(), ShouldHaveSameTypeAs, *new(string))
			So(resp.Request(), ShouldEqual, ``)

			respFail, err := db1.Truncate(`absent`)

			So(err, ShouldBeError)
			So(respFail, ShouldHaveSameTypeAs, traceFall.Response{})
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
	if err != nil {
		log.Fatal(err)
	}

	l := traceFall.NewLog(`Test`)

	resp, err := db.Send(l)

	assert.Nil(t, err)
	assert.NotNil(t, resp.Id)
	assert.Equal(t, l.Id.String(), resp.Id)
	assert.NotEmpty(t, resp.Id)
	assert.IsType(t, traceFall.Response{}, resp)
	assert.IsType(t, traceFall.ResponseData{}, resp.Data)
	assert.IsType(t, *new(string), resp.Id)

	l2, err := l.CreateChild(`Test2`)
	assert.Nil(t, err)

	l2.SetEnvironment(`prod`).Success().ThreadFinish().Tags.Add(`child`)
	l2.Notes.Add(`step`, `note1`).
		Add(`step`, `note2`).
		Add(`id`, l.Id.String())
	l2.Data.Set(`thread`, `thread: `+l.Thread.String())

	resp2, err := db.Send(l2)

	assert.Nil(t, err)
	assert.Equal(t, l2.Id.String(), resp2.Id)
	assert.Equal(t, l2.Parent.Id.String(), l.Id.String())
	assert.NotEmpty(t, resp2.Id)
	assert.IsType(t, traceFall.Response{}, resp2)
	assert.IsType(t, traceFall.ResponseData{}, resp2.Data)
	assert.IsType(t, *new(string), resp2.Id)

	lGet, err := db.Get(l.Id)
	assert.Nil(t, err)
	assert.IsType(t, traceFall.Response{}, lGet)

	assert.IsType(t, traceFall.Log{}, lGet.Data[`log`])
	assert.True(t, lGet.Result)
	assert.Nil(t, lGet.Error)

	logGet := lGet.Data[`log`].(traceFall.Log)
	assert.Equal(t, l.Id.String(), logGet.Id.String())
	assert.Equal(t, l.Environment, `dev`)

	resp3, err := db.RemoveByTags([]string{`child`})
	assert.Nil(t, err)
	assert.IsType(t, traceFall.Response{}, resp3)
	assert.IsType(t, traceFall.ResponseData{}, resp3.Data)
	assert.Nil(t, resp3.Error)

	resp4, err := db.RemoveThread(l.Thread)
	assert.Nil(t, err)
	assert.IsType(t, traceFall.Response{}, resp4)
	assert.IsType(t, traceFall.ResponseData{}, resp4.Data)
	assert.Nil(t, resp4.Error)
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

			So(resp, ShouldHaveSameTypeAs, traceFall.Response{})
			So(resp.Id, ShouldNotBeNil)
			So(resp.Id, ShouldHaveSameTypeAs, *new(string))
			So(resp.Data, ShouldHaveSameTypeAs, traceFall.ResponseData{})
			//So(resp.Request(), ShouldEqual, l.String())

			So(resp.Error, ShouldBeNil)

			So(resp.Id, ShouldEqual, l.Id.String())
			So(resp.Result, ShouldBeTrue)

			l2, err := l.CreateChild(`Child`)
			So(err, ShouldBeNil)
			l2.SetEnvironment(`prod`).Success().ThreadFinish().Tags.Add(`child`)
			l2.Notes.Add(`step`, `note1`).Add(`step`, `note2`)
			l2.Data.Set(`id`, l2.Id).
				Set(`thread`, `thread: `+l2.Thread.String())

			resp2, err := db.Send(l2)

			assert.Nil(t, err)
			assert.Nil(t, resp2.Error)
			assert.NotNil(t, resp2.Id)
			assert.Equal(t, l2.Id.String(), resp2.Id)
			assert.True(t, resp2.Result)

			l3, err := l.CreateChild(`Child`)
			assert.Nil(t, err)

			l3.Success().Tags.
				Add(`child`).Add(`2`)
			l3.Data.Set(`step`, `note1`).Set(`id`, l3.Id)

			resp3, err := db.Send(l3)

			assert.Nil(t, err)
			assert.Nil(t, resp3.Error)
			assert.NotNil(t, resp3.Id)
			assert.Equal(t, l3.Id.String(), resp3.Id)
			assert.True(t, resp3.Result)

			l4, err := l3.CreateChild(`SubChild of Child`)
			assert.Nil(t, err)
			l4.Success().ThreadFinish().Tags.Add(`child`)

			resp4, err := db.Send(l4)

			assert.Nil(t, err)
			assert.Nil(t, resp4.Error)
			assert.NotNil(t, resp4.Id)
			assert.Equal(t, l4.Id.String(), resp4.Id)
			assert.True(t, resp4.Result)

			/// get

			lGet, err := db.Get(l.Id)
			logRootGet := lGet.Data[`log`].(traceFall.Log)

			assert.Equal(t, logRootGet.Id.String(), l.Id.String())

			assert.Nil(t, logRootGet.Parent)
			assert.Equal(t, l.Time.UnixNano(), logRootGet.Time.UnixNano())
			assert.Equal(t, l.Time.Format(time.RFC3339Nano), logRootGet.Time.Format(time.RFC3339Nano))
			assert.Equal(t, l.TimeEnd, logRootGet.TimeEnd)
			assert.Equal(t, l.Error, logRootGet.Error)
			assert.Equal(t, l.App, logRootGet.App)
			assert.Equal(t, l.Name, logRootGet.Name)
			assert.Equal(t, l.Environment, logRootGet.Environment)
			assert.Equal(t, l.Id.String(), logRootGet.Id.String())
			assert.Equal(t, l.Finish, logRootGet.Finish)
			assert.Equal(t, l.Result, logRootGet.Result)
			assert.Equal(t, l.Tags, logRootGet.Tags)

			l2Get, err := db.Get(l2.Id)
			log2Get := l2Get.Data[`log`].(traceFall.Log)

			assert.Equal(t, log2Get.Id.String(), l2.Id.String())

			assert.NotNil(t, log2Get.Parent)
			assert.Equal(t, log2Get.Parent.Id, logRootGet.Id)
			assert.Equal(t, l2.Time.UnixNano(), log2Get.Time.UnixNano())
			assert.Equal(t, l2.Time.Format(time.RFC3339Nano), log2Get.Time.Format(time.RFC3339Nano))
			assert.Equal(t, l2.TimeEnd.UnixNano(), log2Get.TimeEnd.UnixNano())
			assert.Equal(t, l2.TimeEnd.Format(time.RFC3339Nano), log2Get.TimeEnd.Format(time.RFC3339Nano))
			assert.Equal(t, l2.Error, log2Get.Error)
			assert.Equal(t, l2.App, log2Get.App)
			assert.Equal(t, l2.Name, log2Get.Name)
			assert.Equal(t, l2.Environment, log2Get.Environment)
			assert.Equal(t, l2.Id.String(), log2Get.Id.String())
			assert.Equal(t, l2.Finish, log2Get.Finish)
			assert.Equal(t, l2.Result, log2Get.Result)
			assert.Equal(t, l2.Tags, log2Get.Tags)

			l3Get, err := db.Get(l3.Id)
			log3Get := l3Get.Data[`log`].(traceFall.Log)

			assert.Equal(t, log3Get.Id.String(), l3.Id.String())

			assert.NotNil(t, log3Get.Parent)
			assert.Equal(t, log3Get.Parent.Id, logRootGet.Id)
			assert.Equal(t, l3.Time.UnixNano(), log3Get.Time.UnixNano())
			assert.Equal(t, l3.Time.Format(time.RFC3339Nano), log3Get.Time.Format(time.RFC3339Nano))
			assert.Equal(t, l3.TimeEnd.UnixNano(), log3Get.TimeEnd.UnixNano())
			assert.Equal(t, l3.TimeEnd.Format(time.RFC3339Nano), log3Get.TimeEnd.Format(time.RFC3339Nano))
			assert.Equal(t, l3.Error, log3Get.Error)
			assert.Equal(t, l3.App, log3Get.App)
			assert.Equal(t, l3.Name, log3Get.Name)
			assert.Equal(t, l3.Environment, log3Get.Environment)
			assert.Equal(t, l3.Id.String(), log3Get.Id.String())
			assert.Equal(t, l3.Finish, log3Get.Finish)
			assert.Equal(t, l3.Result, log3Get.Result)
			assert.Equal(t, l3.Tags, log3Get.Tags)

			l4Get, err := db.Get(l4.Id)
			log4Get := l4Get.Data[`log`].(traceFall.Log)

			assert.Equal(t, log4Get.Id.String(), l4.Id.String())

			assert.NotNil(t, log4Get.Parent)
			assert.Equal(t, log4Get.Parent.Id, log3Get.Id)
			assert.Equal(t, l4.Time.UnixNano(), log4Get.Time.UnixNano())
			assert.Equal(t, l4.Time.Format(time.RFC3339Nano), log4Get.Time.Format(time.RFC3339Nano))
			assert.Equal(t, l4.TimeEnd.UnixNano(), log4Get.TimeEnd.UnixNano())
			assert.Equal(t, l4.TimeEnd.Format(time.RFC3339Nano), log4Get.TimeEnd.Format(time.RFC3339Nano))
			assert.Equal(t, l4.Error, log4Get.Error)
			assert.Equal(t, l4.App, log4Get.App)
			assert.Equal(t, l4.Name, log4Get.Name)
			assert.Equal(t, l4.Environment, log4Get.Environment)
			assert.Equal(t, l4.Id.String(), log4Get.Id.String())
			assert.Equal(t, l4.Finish, log4Get.Finish)
			assert.Equal(t, l4.Result, log4Get.Result)
			assert.Equal(t, l4.Tags, log4Get.Tags)

			lThreadResp, err := db.GetThread(l4.Thread)
			assert.Nil(t, err)
			assert.Nil(t, lThreadResp.Error)

			lThread := lThreadResp.Data[`list`].([]*traceFall.Log)
			assert.Equal(t, 4, len(lThread))

			for _, v := range lThread {
				assert.IsType(t, &traceFall.Log{}, v)
			}
		})
	})
}
