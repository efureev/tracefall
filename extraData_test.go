package traceFall

import (
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
)

func TestExtraParams(t *testing.T) {

	Convey("Extra Params Struct", t, func() {
		params := NewExtraData()

		Convey("Init", func() {
			So(params, ShouldHaveSameTypeAs, ExtraData{})
			So(reflect.ValueOf(params).Kind(), ShouldEqual, reflect.Map)
		})

		Convey("Set & Get", func() {
			params.Set(`key1`, `value`)

			So(params.Get(`key1`), ShouldEqual, `value`)
			So(params.Get(`key2`), ShouldBeNil)

			params.Set(`key1`, `value 2`)
			So(params.Get(`key1`), ShouldEqual, `value 2`)
		})

		Convey("ToJson", func() {
			params[`test`] = `value`
			params[`dig`] = 123

			So(string(params.ToJson()), ShouldEqual, `{"dig":123,"test":"value"}`)
		})
	})
}
