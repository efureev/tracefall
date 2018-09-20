package traceFall

import (
	"reflect"
	"testing"
)

func TestTags(t *testing.T) {

	Convey("Tags Struct", t, func() {

		tags := Tags{}

		Convey("Init", func() {
			So(tags, ShouldHaveSameTypeAs, Tags{})
			So(reflect.ValueOf(tags).Kind(), ShouldEqual, reflect.Slice)

			So(tags.List(), ShouldHaveSameTypeAs, []string{})
			So(len(tags.List()), ShouldBeZeroValue)
		})

		Convey("Add", func() {
			tags.Add(`first`)
			So(len(tags.List()), ShouldEqual, 1)
			tags.Add(`second`)
			So(len(tags.List()), ShouldEqual, 2)

			So(tags[0], ShouldEqual, `first`)
		})

		Convey("Clear", func() {
			tags.Add(`first`)
			So(len(tags.List()), ShouldEqual, 1)
			tags.Add(`second`)
			So(len(tags), ShouldEqual, 2)

			tags.Clear()
			So(len(tags), ShouldBeZeroValue)
		})
	})
}
