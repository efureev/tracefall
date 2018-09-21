package traceFall

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestThread(t *testing.T) {

	Convey("Thread Struct", t, func() {

		Convey("Init", func() {
			thread := Thread{}
			thread.Add(NewLog(`test`).ToLogJSON())

			So(len(thread), ShouldEqual, 1)
			thread.Add(NewLog(`test 2`).ToLogJSON())
			So(len(thread), ShouldEqual, 2)
		})

		Convey("From List", func() {

			thread := ThreadFromList([]*LogJSON{
				NewLog(`1`).ToLogJSON(),
				NewLog(`2`).ToLogJSON(),
				NewLog(`3`).ToLogJSON(),
			})

			So(len(thread), ShouldEqual, 3)

		})

	})
}
