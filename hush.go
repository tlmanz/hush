package hush

import (
	"context"
	"reflect"
)

type HushType string

// Husher is the interface that wraps the Hush method.
type Husher interface {
	Hush(ctx context.Context, v interface{}, args ...interface{}) ([][]string, error)
}

type hushType struct{}

// Constants used throughout the package
const (
	TagMask HushType = "mask"
	TagHide HushType = "hide"

	DefaultSeparator = "."
	HiddenValue      = "HIDDEN"
)

// NewHush creates a new Husher instance.
// It accepts either a struct or a string and returns a Husher interface.
func NewHush() Husher {
	return &hushType{}
}

func (ht *hushType) Hush(ctx context.Context, v interface{}, args ...interface{}) ([][]string, error) {
	opts := &hushOptions{
		separator:      DefaultSeparator,
		maskFunc:       defaultMaskFunc,
		includePrivate: false,
		prefix:         "",
	}

	for _, option := range args {
		switch opt := option.(type) {
		case string:
			opts.prefix = opt
		case HushType:
			opts.hushType = opt
		case Option:
			opt(opts)
		}
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return ht.processValue(ctx, opts.prefix, reflect.StructField{}, rv, opts)
}
