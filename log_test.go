package traceFall

import (
	"errors"
	"github.com/satori/go.uuid"
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
	"time"
)

func TestLog(t *testing.T) {

	Convey("Log Struct", t, func() {

		timeStart := time.Now()
		log := NewLog(`test log`)

		Convey("Init", func() {
			So(log, ShouldHaveSameTypeAs, &Log{})
			So(reflect.ValueOf(log).Kind(), ShouldEqual, reflect.Ptr)

			So(log.Time, ShouldHappenAfter, timeStart)
			So(log.Time, ShouldHappenBefore, time.Now())

			So(log.Id, ShouldHaveSameTypeAs, uuid.UUID{})
			So(log.Id, ShouldNotBeNil)

			So(log.Thread, ShouldHaveSameTypeAs, uuid.UUID{})
			So(log.Thread, ShouldNotBeNil)

			So(log.Thread, ShouldEqual, log.Id)

			So(log.Name, ShouldEqual, `test log`)

			So(log.App, ShouldEqual, `App`)
			So(log.Environment, ShouldEqual, EnvironmentDev)
			So(log.Result, ShouldBeFalse)

			So(log.Data, ShouldHaveSameTypeAs, ExtraData{})
			So(log.Notes, ShouldHaveSameTypeAs, NoteGroups{})
		})

		Convey("Set Name", func() {
			So(log.Name, ShouldEqual, `test log`)
			log.SetName(`test 2`)
			So(log.Name, ShouldEqual, `test 2`)

		})

		Convey("Add Data", func() {
			log.Tags.Add(`first`)
			log.Tags.Add(`second`)

			So(log.Tags.List(), ShouldResemble, []string{`first`, `second`})
			So(log.Tags, ShouldResemble, Tags{`first`, `second`})
		})

		Convey("Add Note", func() {
			log.Notes.Add(`first group`, `first note`)
			log.Notes.Add(`second group`, `first note`)
			log.Notes.Add(`first group`, `second note`)
			log.Notes.Add(`first group`, `third note`)

			So(log.Notes.Count(), ShouldEqual, 2)
			So(log.Notes.Get(`first group`).Count(), ShouldEqual, 3)

			log.Notes.Get(`second group`).Clear()
			So(log.Notes.Get(`second group`).Count(), ShouldEqual, 0)

			log.Notes.Get(`first group`).Add(`adding`)
			So(log.Notes.Get(`first group`).Count(), ShouldEqual, 4)
		})

		Convey("Time", func() {
			So(log.Time, ShouldNotBeNil)
			So(log.Time, ShouldHaveSameTypeAs, time.Time{})
			So(reflect.ValueOf(log.Time).Kind(), ShouldEqual, reflect.Struct)
			So(log.Time, ShouldHappenBefore, time.Now())
		})

		Convey("TimeEnd", func() {
			So(log.TimeEnd, ShouldBeNil)
			log.FinishTimeEnd()
			So(log.TimeEnd, ShouldNotBeNil)
			So(log.TimeEnd, ShouldHaveSameTypeAs, &time.Time{})
			So(reflect.ValueOf(log.TimeEnd).Kind(), ShouldEqual, reflect.Ptr)
			So(*log.TimeEnd, ShouldHappenBefore, time.Now())
			log.FinishTimeEnd().FinishTimeEnd().FinishTimeEnd().FinishTimeEnd()
			So(*log.TimeEnd, ShouldHappenBefore, time.Now())
		})

		Convey("Thread Finish", func() {
			So(log.Finish, ShouldBeFalse)
			log.ThreadFinish()
			So(log.Finish, ShouldBeTrue)
			So(reflect.ValueOf(log.Finish).Kind(), ShouldEqual, reflect.Bool)
		})

		Convey("Result", func() {

			Convey("Success", func() {
				So(log.Result, ShouldBeFalse)
				So(log.TimeEnd, ShouldBeNil)

				log.Success()

				So(log.TimeEnd, ShouldNotBeNil)
				So(log.TimeEnd, ShouldHaveSameTypeAs, &time.Time{})
				So(reflect.ValueOf(log.TimeEnd).Kind(), ShouldEqual, reflect.Ptr)
				So(*log.TimeEnd, ShouldHappenBefore, time.Now())

				So(log.Result, ShouldBeTrue)
			})

			Convey("Fail", func() {
				So(log.Result, ShouldBeFalse)
				So(log.TimeEnd, ShouldBeNil)
				So(log.Error, ShouldBeNil)

				log.Fail(errors.New(`test errors`))

				So(log.TimeEnd, ShouldNotBeNil)
				So(log.TimeEnd, ShouldHaveSameTypeAs, &time.Time{})
				So(reflect.ValueOf(log.TimeEnd).Kind(), ShouldEqual, reflect.Ptr)
				So(*log.TimeEnd, ShouldHappenBefore, time.Now())

				So(log.Result, ShouldBeFalse)
				So(log.Error, ShouldBeError)
				So(log.Error.Error(), ShouldEqual, `test errors`)

				log.Success()

				log.Fail(nil)
				So(log.TimeEnd, ShouldNotBeNil)
				So(log.TimeEnd, ShouldHaveSameTypeAs, &time.Time{})
				So(reflect.ValueOf(log.TimeEnd).Kind(), ShouldEqual, reflect.Ptr)
				So(*log.TimeEnd, ShouldHappenBefore, time.Now())
				So(log.Error, ShouldBeNil)
			})
		})

		Convey("Environment", func() {
			So(log.Environment, ShouldEqual, EnvironmentDev)
			log.SetEnvironment(EnvironmentProd)
			So(log.Environment, ShouldEqual, EnvironmentProd)
			log.SetEnvironment(EnvironmentTest)
			So(log.Environment, ShouldEqual, EnvironmentTest)
		})
	})
}
