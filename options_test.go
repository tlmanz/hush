package hush

import (
	"testing"
)

func TestOptions(t *testing.T) {
	opts := &hushOptions{}

	WithSeparator("_")(opts)
	if opts.separator != "_" {
		t.Errorf("WithSeparator() failed, got: %s, want: _", opts.separator)
	}

	customMask := func(s string) string { return "MASKED" }
	WithMaskFunc(customMask)(opts)
	if opts.maskFunc("test") != "MASKED" {
		t.Errorf("WithMaskFunc() failed, got: %s, want: MASKED", opts.maskFunc("test"))
	}

	WithPrivateFields(true)(opts)
	if !opts.includePrivate {
		t.Errorf("WithPrivateFields(true) failed, got: false, want: true")
	}
}
