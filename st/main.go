package st

import (
	"errors"
	"reflect"
	"strings"
)

// UpdateStruct updates the field values of struct 's' with values given 'u' for the tag given by string 'tag'.
// The key for 'u' should match the values for 'tag' for each field.
// Supports targeting child tags such as the child tag 'fizz' in the tag `foo:"bar:baz;fizz:buzz"` with the notation: "foo:fizz".
func UpdateStructByTag(s interface{}, u map[string]interface{}, tag string) error {
	v, err := extractStructValue(s)
	if err != nil {
		return err
	}
	t := reflect.TypeOf(v.Interface())
	fieldCount := v.NumField()
	for i := 0; i < fieldCount; i++ {
		parentChild := strings.Split(tag, ":")
		var tagValue string
		switch len(parentChild) {
		case 1:
			tagValue = t.Field(i).Tag.Get(tag)
			break
		case 2:
			tagValue = extractSubTag(t.Field(i), parentChild[0], parentChild[1])
			break
		default:
			return errors.New("bad tag")
		}
		if tagValue == "" || u[tagValue] == nil {
			continue
		}
		newValue := reflect.ValueOf(u[tagValue])
		oldValue := v.Field(i)
		oldValue.Set(newValue)
	}
	return nil
}

// ExtractStructValue returns the reflect.Value for the underlying struct of the pointer i.
// If i is not a pointer to a struct then error will be non-nil and returned value should be ignored.
func extractStructValue(i interface{}) (reflect.Value, error) {
	v := reflect.ValueOf(i)
	err := errors.New("wrong kind")
	if v.Kind() != reflect.Ptr {
		return v, err
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return v, err
	}
	return v, nil
}

func extractSubTag(f reflect.StructField, parentTag, childTag string) string {
	testTags := strings.Split(f.Tag.Get(parentTag), ";")
	for _, tag := range testTags {
		isColTag := strings.Contains(tag, childTag)
		if isColTag {
			return strings.Split(tag, ":")[1]
		}
	}
	return ""
}
