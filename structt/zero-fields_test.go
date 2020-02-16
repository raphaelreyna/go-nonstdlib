package structt

import "testing"

func TestZeroFields(t *testing.T) {
	type A struct {
		First string
		Second string
		Third int
	}
	a := A{
		First: "firstvalue",
		Second: "secondvalue",
		Third: 3,
	}
	ZeroFields(&a, "First", "Third", "Fourth")
	if a.First != "" {
		t.Errorf(`first field: "%s" != ""`, a.First)
	}
	if a.Second != "secondvalue" {
		t.Errorf(`second field: "%s" != "secondvalue"`, a.Second)
	}
	if a.Third != 0 {
		t.Errorf(`third field: %d != 0`, a.Third)
	}
}
