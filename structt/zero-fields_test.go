package structt

import "testing"

func TestZeroFields(t *testing.T) {
	a := struct {
		First string
		Second string
		Third int
	}{
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

func TestHasZeroes(t *testing.T) {
	type A struct {
		First string
		Second int
		Third bool
		Fourth string
	}
	type test struct {
		description string
		a A
		exceptions []string
		want bool
	}
	tt := []test{
		test{
			description: "no zeroes, no exceptions",
			a: A{
				First: "firstvalue",
				Second: 5,
				Third: true,
				Fourth: "fourthvalue",
			},
			want: false,
		},
		test{
			description: "has zero, no exceptions",
			a: A{
				First: "firstvalue",
				Second: 5,
				Third: true,
			},
			want: true,
		},
		test{
			description: "no zeroes, has exception",
			a: A{
				First: "firstvalue",
				Second: 5,
				Third: true,
				Fourth: "fourthvalue",
			},
			exceptions: []string{
				"Fourth",
			},
			want: false,
		},
		test{
			description: "has zero and exception",
			a: A{
				First: "firstvalue",
				Second: 5,
				Third: true,
			},
			exceptions: []string{
				"Fourth",
			},
			want: false,
		},
		test{
			description: "no zeroes, exceptions has non existing field name",
			a: A{
				First: "firstvalue",
				Second: 5,
				Third: true,
				Fourth: "fourthvalue",
			},
			exceptions: []string{
				"Fifth",
			},
			want: false,
		},
		test{
			description: "has zero, exceptions has non existing field name",
			a: A{
				First: "firstvalue",
				Second: 5,
				Third: true,
			},
			exceptions: []string{
				"Fourth",
				"Fifth",
			},
			want: false,
		},
		test{
			description: "has zeroes and exception",
			a: A{
				First: "firstvalue",
				Third: true,
			},
			exceptions: []string{
				"Fourth",
			},
			want: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			hz := HasZeroes(tc.a, tc.exceptions...)
			if hz != tc.want {
				t.Errorf("got: %v; expected: %v", hz, tc.want)
			}
		})
	}
}

