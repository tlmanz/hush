package hush

import (
	"fmt"
	"reflect"
)

// Constants used throughout the package
const (
	TagMask = "mask"
	TagHide = "hide"

	DefaultSeparator = "."
	HiddenValue      = "HIDDEN"
)

// NewHush creates a new Husher instance.
// It accepts either a struct or a string and returns a Husher interface.
func NewHush(v interface{}) (Husher, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		return &hushType{value: rv, isStr: false}, nil
	case reflect.String:
		return &hushType{value: rv, isStr: true}, nil
	default:
		return nil, fmt.Errorf("expected struct or string, got %v", rv.Kind())
	}
}
