package csv

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// StructScanner also does things
type StructScanner struct {
	*ColumnScanner
}

// NewStructScanner does things
func NewStructScanner(reader io.Reader, options ...Option) (*StructScanner, error) {
	inner, err := NewColumnScanner(reader, options...)
	if err != nil {
		return nil, err
	}
	return &StructScanner{ColumnScanner: inner}, nil
}

// Populate scans values into a struct
func (scanner *StructScanner) Populate(v interface{}) error {
	rtype := reflect.TypeOf(v)
	if rtype.Kind() != reflect.Ptr {
		return fmt.Errorf("provided value must be reflect.Ptr. You provided [%v] ([%v])", v, rtype.Kind())
	}

	value := reflect.ValueOf(v)
	if value.IsNil() {
		return fmt.Errorf("the provided value was nil. Please provide a non-nil pointer")
	}

	scanner.populate(rtype.Elem(), value.Elem())
	return nil
}

func (scanner *StructScanner) populate(rtype reflect.Type, value reflect.Value) {
	for x := 0; x < rtype.NumField(); x++ {
		column := rtype.Field(x).Tag.Get("csv")

		_, found := scanner.columnIndex[column]
		if !found {
			continue
		}

		field := value.Field(x)
		if field.Kind() == reflect.Int64 {
			intValue, err := strconv.Atoi(scanner.Column(column))
			if err != nil {
				continue
			}
			field.SetInt(int64(intValue))
			continue
		} else if field.Kind() == reflect.Float64 {
			floatValue, err := strconv.ParseFloat(scanner.Column(column), 64)
			if err != nil {
				continue
			}
			field.SetFloat(floatValue)
			continue
		} else if field.Kind() != reflect.String {
			continue
		} else if !field.CanSet() {
			continue // Future: return err?
		}

		field.SetString(scanner.Column(column))
	}
}
