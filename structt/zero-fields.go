package structt

import "reflect"

// ZeroFields zeroes out the values for the passed in field names.
// Panics if i is not a pointer to a struct.
func ZeroFields(i interface{}, fields... string) {
	iv := reflect.ValueOf(i).Elem()
	for _, fn := range fields {
		df := iv.FieldByName(fn)
		if !df.IsValid() {
			continue
		}
		if !df.IsZero() {
			zero := reflect.Zero(df.Type())
			df.Set(zero)
		}
	}
}

// HasZeroes returns true if i has at least one field with a zero value.
// Ignores zero value cases for field names given as exceptions.
// Panics if i is not a struct or pointer to a struct.
func HasZeroes(i interface{}, exceptions... string) bool {
	iv := reflect.ValueOf(i)
	if iv.Kind() == reflect.Ptr {
		iv = iv.Elem()
	}
	it := iv.Type()
	for j := 0 ; j < it.NumField() ; j++ {
		fv := iv.Field(j)
		if fv.IsZero() {
			fn := it.Field(j).Name
			var exception bool
			for _, e := range exceptions {
				if fn == e {
					exception = true
					break
				}
			}
			if !exception {
				return true
			}
		}
	}
	return false
}
