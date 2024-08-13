package hush

// Option is a function type for configuring hushOptions.
type Option func(*hushOptions)

// hushOptions holds the configuration options for the Hush operation.
type hushOptions struct {
	separator      string
	maskFunc       func(string) string
	includePrivate bool
	prefix         string
	hushType       HushType
}

// WithSeparator sets the separator used for nested field names.
func WithSeparator(sep string) Option {
	return func(o *hushOptions) {
		o.separator = sep
	}
}

// WithMaskFunc sets a custom masking function.
func WithMaskFunc(f func(string) string) Option {
	return func(o *hushOptions) {
		o.maskFunc = f
	}
}

// WithPrivateFields sets whether to include private fields.
func WithPrivateFields(include bool) Option {
	return func(o *hushOptions) {
		o.includePrivate = include
	}
}

// WithOptions sets all options at once
func WithOptions(options *hushOptions) Option {
	return func(o *hushOptions) {
		*o = *options
	}
}
