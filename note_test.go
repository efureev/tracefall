package traceFall

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestNote(t *testing.T) {

	Convey("Note Struct", t, func() {

		timeStart := time.Now()
		note := NewNote(`test Note`)

		So(note, ShouldHaveSameTypeAs, &Note{})
		So(reflect.ValueOf(note).Kind(), ShouldEqual, reflect.Ptr)

		So(note.Time, ShouldBeGreaterThanOrEqualTo, timeStart.UnixNano())
		So(note.Time, ShouldBeLessThan, time.Now().UnixNano())
	})
}

func TestNoteGroup(t *testing.T) {
	Convey("Group Note Struct", t, func() {

		group := NewNoteGroup(`test group`)

		Convey("Init", func() {
			So(group, ShouldHaveSameTypeAs, &NoteGroup{})
			So(group.Label, ShouldEqual, `test group`)
			So(group.Count(), ShouldBeZeroValue)
			So(reflect.ValueOf(group).Kind(), ShouldEqual, reflect.Ptr)
		})

		Convey("Add", func() {
			group.Add(`first note`)
			So(group.Count(), ShouldEqual, 1)
			So(group.Notes, ShouldHaveSameTypeAs, Notes{})
			group.Add(`second note`)
			So(group.Count(), ShouldEqual, 2)
		})

		Convey("Clear", func() {
			group.Clear()
			So(group.Count(), ShouldBeZeroValue)
			So(group.Label, ShouldEqual, `test group`)
		})
	})
}

func TestNoteGroups(t *testing.T) {
	Convey("Groups Note Struct", t, func() {

		groups := NewNotesGroups()

		Convey("Init", func() {
			So(groups, ShouldHaveSameTypeAs, NoteGroups{})
			So(groups.Count(), ShouldBeZeroValue)
			So(reflect.ValueOf(groups).Kind(), ShouldEqual, reflect.Map)
		})

		Convey("Add", func() {
			groups.Add(`first group`, `first note`)
			So(groups.Count(), ShouldEqual, 1)

			groups.Add(`first group`, `second note`)
			So(groups.Count(), ShouldEqual, 1)

			groups.Add(`second group`, `first note`)
			So(groups.Count(), ShouldEqual, 2)
		})

		Convey("AddGroup", func() {
			groups.AddGroup(`first group`, []string{`first note`, `second note`})
			So(groups.Count(), ShouldEqual, 1)
			So(groups.Get(`first group`).Count(), ShouldEqual, 2)
			So(groups.Get(`first group`).Notes[0].Note, ShouldEqual, `first note`)
			So(groups.Get(`first group`).Notes[1].Note, ShouldEqual, `second note`)

			groups.AddGroup(`first group`, []string{`third note`, `next note`, `next2 note`})
			So(groups.Get(`first group`).Count(), ShouldEqual, 5)
			So(groups.Get(`first group`).Notes[3].Note, ShouldEqual, `next note`)
		})

		Convey("Get & Clear", func() {
			groups.Add(`first group`, `first note`)
			So(groups.Count(), ShouldEqual, 1)

			groups.Add(`first group`, `second note`)
			So(groups.Count(), ShouldEqual, 1)

			groups.Add(`second group`, `first note`)
			So(groups.Count(), ShouldEqual, 2)

			groupFirst := groups.Get(`first group`)
			So(groupFirst.Count(), ShouldEqual, 2)
			So(groupFirst, ShouldHaveSameTypeAs, &NoteGroup{})
			So(reflect.ValueOf(groupFirst).Kind(), ShouldEqual, reflect.Ptr)

			groupFirst.Add(`adding 1`)

			groupFirst2 := groups.Get(`first group`)
			So(groupFirst2.Count(), ShouldEqual, 3)

			groupFirst.Add(`adding 2`)
			So(groupFirst2.Count(), ShouldEqual, 4)

			groupFirst2.Add(`adding 3`)
			So(groupFirst.Count(), ShouldEqual, 5)

			groupFirst2.Clear()
			So(groupFirst.Count(), ShouldBeZeroValue)
			So(groupFirst2.Count(), ShouldBeZeroValue)

			groupFail := groups.Get(`fail group`)
			So(groupFail, ShouldBeNil)
		})

		Convey("Remove", func() {
			groups.
				Add(`first group`, `first note`).
				Add(`second group`, `first note`)
			So(groups.Count(), ShouldEqual, 2)

			groups.Remove(`first group`)
			So(groups.Count(), ShouldEqual, 1)

			groups.Remove(`absent group`)
			So(groups.Count(), ShouldEqual, 1)

			groups.Remove(`second group`)
			So(groups.Count(), ShouldBeZeroValue)
		})

		Convey("ToJson", func() {
			groups.
				Add(`first group`, `first note`)
			note := groups.Get(`first group`).Notes[0]

			expected := fmt.Sprintf(`[{"notes":[{"t":%d,"v":"first note"}],"label":"first group"}]`, note.Time)
			So(groups.ToJsonString(), ShouldEqual, expected)
		})
		/*
				Convey("FromJson", func() {

					t1 := time.Now().UnixNano() - 2131
					t2 := time.Now().UnixNano()
					jsonString := fmt.Sprintf(`[{"notes":[{"t":%d,"v":"first json json note 1"},{"t":%d,"v":"first json json note 2"}],"label":"first json group"}]`, t1, t2)

					err:= groups.
						FromJson(jsonString)
					So(err, ShouldBeNil)

					So(groups.Count(), ShouldEqual, 1)

					noteGroup := groups.Get(`first json group`)
					So(noteGroup.Count(), ShouldEqual, 2)

					So(noteGroup.Notes[0].Time, ShouldEqual, t1)
					So(noteGroup.Notes[1].Time, ShouldEqual, t2)

					So(noteGroup.Notes[0].Note, ShouldEqual, `first json json note 1`)
					So(noteGroup.Notes[1].Note, ShouldEqual, `first json json note 2`)


				})*/
	})
}
