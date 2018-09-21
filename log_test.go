package traceFall

import (
	"errors"
	"fmt"
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

			So(log.ID, ShouldHaveSameTypeAs, uuid.UUID{})
			So(log.ID, ShouldNotBeNil)

			So(log.Thread, ShouldHaveSameTypeAs, uuid.UUID{})
			So(log.Thread, ShouldNotBeNil)

			So(log.Thread, ShouldEqual, log.ID)

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

		Convey("Parent", func() {
			So(log.Parent, ShouldBeNil)

			Convey("set parent ID", func() {
				parentID := generateUUID()

				log.SetParentID(parentID)

				So(log.Parent, ShouldHaveSameTypeAs, &Log{})
				So(reflect.ValueOf(log.Parent).Kind(), ShouldEqual, reflect.Ptr)

				So(log.Parent.ID, ShouldHaveSameTypeAs, uuid.UUID{})
				So(log.Parent.ID, ShouldNotBeNil)
				So(log.Parent.ID, ShouldEqual, parentID)

				So(log.Parent.Thread, ShouldHaveSameTypeAs, uuid.UUID{})
				So(log.Parent.Thread, ShouldNotBeNil)
				So(log.Parent.Thread.String(), ShouldEqual, log.Thread.String())
				So(log.Parent.Thread, ShouldEqual, log.Thread)

			})

			Convey("set parent Log Struct", func() {

				Convey("valid Parent set", func() {
					parenLog := NewLog(`ParenLog`)
					parenLog.Thread = log.Thread
					err := log.SetParent(parenLog)
					So(err, ShouldBeNil)

					So(log.Parent, ShouldHaveSameTypeAs, &Log{})
					So(reflect.ValueOf(log.Parent).Kind(), ShouldEqual, reflect.Ptr)

					So(log.Parent.ID, ShouldHaveSameTypeAs, uuid.UUID{})
					So(log.Parent.ID, ShouldNotBeNil)
					So(log.Parent.ID, ShouldEqual, parenLog.ID)

					So(log.Parent.Thread, ShouldHaveSameTypeAs, uuid.UUID{})
					So(log.Parent.Thread, ShouldNotBeNil)
					So(log.Parent.Thread, ShouldEqual, log.Thread)
				})

				Convey("invalid the Parent by Thread", func() {
					parenLog := NewLog(`ParenLog`)

					err := log.SetParent(parenLog)

					So(err, ShouldBeError)
					So(log.Parent, ShouldBeNil)
					So(err, ShouldEqual, ErrorParentThreadDiff)
				})

				Convey("invalid the Parent by Finish", func() {
					parenLog := NewLog(`ParenLog`)
					parenLog.Thread = log.Thread
					parenLog.ThreadFinish()

					err := log.SetParent(parenLog)

					So(err, ShouldBeError)
					So(log.Parent, ShouldBeNil)
					So(err, ShouldEqual, ErrorParentFinish)
				})
			})

		})

		Convey("String", func() {
			So(log.String(), ShouldEqual, fmt.Sprintf("[%s] %s", log.Time, log.Name))
		})

		Convey("Create Child: Success", func() {
			timeStart := time.Now()

			child, err := log.CreateChild(`child`)
			So(err, ShouldBeNil)

			So(child, ShouldHaveSameTypeAs, &Log{})
			So(reflect.ValueOf(child).Kind(), ShouldEqual, reflect.Ptr)

			So(child.Time, ShouldHappenAfter, timeStart)
			So(child.Time, ShouldHappenBefore, time.Now())

			So(child.ID, ShouldHaveSameTypeAs, uuid.UUID{})
			So(child.ID, ShouldNotBeNil)

			So(child.Thread, ShouldHaveSameTypeAs, uuid.UUID{})
			So(child.Thread, ShouldNotBeNil)

			So(child.Thread, ShouldEqual, log.Thread)

			So(child.Name, ShouldEqual, `child`)

			So(child.App, ShouldEqual, log.App)
			So(child.Environment, ShouldEqual, log.Environment)

			So(child.Data, ShouldHaveSameTypeAs, ExtraData{})
			So(child.Notes, ShouldHaveSameTypeAs, NoteGroups{})

			So(child.String(), ShouldEqual, fmt.Sprintf("[%s] %s", child.Time, child.Name))
		})

		Convey("Create Child: Fail", func() {
			log.ThreadFinish()

			child, err := log.CreateChild(`child`)
			So(err, ShouldBeError)
			So(child, ShouldBeNil)
			So(err, ShouldEqual, ErrorParentFinish)
		})

		Convey("To LogJson Struct", func() {
			log.Tags.Add(`tag 1`)
			log.Notes.Add(`group`, `note 1`)
			log.Data.Set(`key`, `val`)
			lJSON := log.ToLogJSON()

			So(lJSON, ShouldHaveSameTypeAs, &LogJSON{})

			So(lJSON.ID, ShouldEqual, log.ID)
			So(lJSON.Thread, ShouldEqual, log.Thread)
			So(lJSON.Time, ShouldEqual, log.Time.UnixNano())
			So(lJSON.Name, ShouldEqual, log.Name)
			So(lJSON.App, ShouldEqual, log.App)
			So(lJSON.Environment, ShouldEqual, log.Environment)
			So(lJSON.Tags, ShouldHaveSameTypeAs, []string{})
			So(lJSON.Tags, ShouldResemble, log.Tags.List())
			So(lJSON.Parent, ShouldBeNil)
			So(lJSON.TimeEnd, ShouldBeNil)
			So(lJSON.Finish, ShouldEqual, log.Finish)
			So(lJSON.Error, ShouldEqual, log.Error)
			So(lJSON.Result, ShouldEqual, log.Result)
			//So(lJSON.Step, ShouldEqual, log.Step)
			So(lJSON.Data, ShouldResemble, log.Data)
			So(lJSON.Notes, ShouldResemble, log.Notes.prepareToJSON())
		})

		Convey("To Json", func() {
			log.Tags.Add(`tag 1`)
			log.Notes.Add(`group first`, `note 1`)
			log.Data.Set(`key`, `val`)

			Convey("Simple log", func() {
				jsonBytes := log.ToJSON()

				expected := fmt.Sprintf(`{"id":"%s","thread":"%s","name":"test log","app":"%s","time":%d,"timeEnd":null,"result":false,"finish":false,"env":"%s","error":null,"data":{"key":"%s"},"notes":[{"notes":[{"t":%d,"v":"%s"}],"label":"%s"}],"tags":["%s"],"parent":null}`,
					log.ID,
					log.Thread,
					log.App,
					log.Time.UnixNano(),
					log.Environment,
					log.Data.Get(`key`),
					log.Notes.Get(`group first`).Notes[0].Time,
					log.Notes.Get(`group first`).Notes[0].Note,
					log.Notes.Get(`group first`).Label,
					log.Tags[0],
					//log.Step,
				)

				So(string(jsonBytes), ShouldEqual, expected)
			})

			Convey("Error Finish log", func() {
				log.Fail(errors.New(`fail`)).ThreadFinish()
				jsonBytes := log.ToJSON()

				expected := fmt.Sprintf(`{"id":"%s","thread":"%s","name":"test log","app":"%s","time":%d,"timeEnd":%d,"result":false,"finish":true,"env":"%s","error":"%s","data":{"key":"%s"},"notes":[{"notes":[{"t":%d,"v":"%s"}],"label":"%s"}],"tags":["%s"],"parent":null}`,
					log.ID,
					log.Thread,
					log.App,
					log.Time.UnixNano(),
					log.TimeEnd.UnixNano(),
					log.Environment,
					log.Error.Error(),
					log.Data.Get(`key`),
					log.Notes.Get(`group first`).Notes[0].Time,
					log.Notes.Get(`group first`).Notes[0].Note,
					log.Notes.Get(`group first`).Label,
					log.Tags[0],
					//log.Step,
				)

				So(string(jsonBytes), ShouldEqual, expected)
			})
		})

		Convey("Trees", func() {
			root := NewLog(`test log`)
			root.Tags.Add(`tag 1`)
			root.Notes.Add(`group first`, `note 1`)
			root.Data.Set(`key`, `val`)

			parent := root

			for i := 0; i <= 5; i++ {
				child, err := parent.CreateChild(fmt.Sprintf(`log #%d`, i+1))
				So(err, ShouldBeNil)

				parent = child

				So(child.GetLevel(), ShouldEqual, i+1)
			}
		})

	})

	Convey("Log Parent Shadow Struct", t, func() {

		log := NewLog(`test log`)

		Convey("To Shadow", func() {
			shadow := log.ToShadow()

			So(shadow, ShouldHaveSameTypeAs, &LogParentShadow{})
			So(reflect.ValueOf(shadow).Kind(), ShouldEqual, reflect.Ptr)

			So(shadow.ID, ShouldHaveSameTypeAs, uuid.UUID{})
			So(shadow.ID, ShouldNotBeNil)
			So(shadow.ID, ShouldEqual, log.ID)

			So(shadow.Thread, ShouldHaveSameTypeAs, uuid.UUID{})
			So(shadow.Thread, ShouldNotBeNil)
			So(shadow.Thread, ShouldEqual, log.Thread)

		})

		Convey("From Shadow", func() {

			parent := NewLog(`parent log`)
			shadow := parent.ToShadow()

			So(shadow.ID, ShouldEqual, parent.ID)
			So(shadow.Thread, ShouldEqual, parent.Thread)

			log.ParentFromShadow(shadow)

			So(log.Parent, ShouldHaveSameTypeAs, &Log{})
			So(reflect.ValueOf(log.Parent).Kind(), ShouldEqual, reflect.Ptr)

			So(log.Parent.ID, ShouldHaveSameTypeAs, uuid.UUID{})
			So(log.Parent.ID, ShouldNotBeNil)
			So(log.Parent.ID, ShouldEqual, parent.ID)

			So(log.Parent.Thread, ShouldHaveSameTypeAs, uuid.UUID{})
			So(log.Parent.Thread, ShouldNotBeNil)
			So(log.Parent.Thread, ShouldEqual, parent.Thread)
		})

	})
}
