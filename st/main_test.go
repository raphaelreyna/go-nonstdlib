package st

import (
	"fmt"
	"reflect"
	"testing"
)

func TestUpdateStructByTag(t *testing.T) {
	type testStruct struct {
		FieldOne   int    `foo:"bar:first;fizz:buzz" letters:"alpha:a"`
		FieldTwo   string `letters:"beta:b"`
		FieldThree bool   `foo:"bar:third"`
	}
	test := testStruct{
		FieldOne: 4,
		FieldTwo: "test",
	}
	control := testStruct{
		FieldOne:   5,
		FieldTwo:   "test",
		FieldThree: true,
	}
	updates := map[string]interface{}{
		"first":  5,
		"second": "failed",
		"third":  true,
	}
	err := UpdateStructByTag(&test, updates, "foo:bar")
	if err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(control, test) {
		msg := "test struct not updated correctly\nexpected: %v\ngot: %v\n"
		t.Error(fmt.Sprintf(msg, test, updates))
	}

	err = UpdateStructByTag([]string{}, updates, "foo:bar")
	if err == nil {
		t.Error("function accepted non struct kind for parameter 's'")
	}
}
