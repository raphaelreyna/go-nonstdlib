package structt

import "reflect"

// ZeroFields zeroes out the values for the passed in field names.
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
