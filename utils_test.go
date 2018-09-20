package traceFall

import (
	"github.com/satori/go.uuid"
	"testing"
)

func TestGenerateUUID(t *testing.T) {

	Convey("Generate UUID", t, func() {

		id := generateUUID()

		So(id, ShouldHaveSameTypeAs, uuid.UUID{})
		So(id.Version(), ShouldEqual, uint8(4))
	})
}

func TestRemoveDuplicatesFromSlice(t *testing.T) {

	Convey("Remove Duplicates From Slice", t, func() {

		slice := []string{`1`, `1`, `2`, `1`, `1`, `3`}
		uniqueSlice := removeDuplicatesFromSlice(slice)

		So(uniqueSlice, ShouldResemble, []string{`1`, `2`, `3`})
	})
}
