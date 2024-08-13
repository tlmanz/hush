package hush

import (
	"context"
	"reflect"
	"testing"
)

func TestProcessField(t *testing.T) {
	opts := &hushOptions{
		separator:      ".",
		maskFunc:       defaultMaskFunc,
		includePrivate: false,
	}

	tests := []struct {
		name      string
		fieldName string
		field     reflect.StructField
		value     interface{}
		opts      *hushOptions
		want      [][]string
		wantErr   bool
	}{
		{"String", "test", reflect.StructField{Tag: reflect.StructTag(`hush:"mask"`)}, "sensitive", opts, [][]string{{"test", "se*****ve"}}, false},
		{"Int", "number", reflect.StructField{}, 42, opts, [][]string{{"number", "42"}}, false},
		{"Slice", "slice", reflect.StructField{}, []string{"a", "b"}, opts, [][]string{{"slice[0]", "a"}, {"slice[1]", "b"}}, false},
		{"Map", "map", reflect.StructField{}, map[string]int{"a": 1, "b": 2}, opts, [][]string{{"map[a]", "1"}, {"map[b]", "2"}}, false},
		{"Struct", "struct", reflect.StructField{}, struct {
			Name string `hush:"mask"`
		}{"John"}, opts, [][]string{{"struct.Name", "****"}}, false},
		{"Pointer", "ptr", reflect.StructField{}, &struct {
			Name string `hush:"hide"`
		}{"John"}, opts, [][]string{{"ptr.Name", "HIDDEN"}}, false},
		{"Nil Pointer", "ptr", reflect.StructField{}, (*struct{ Name string })(nil), opts, [][]string{{"ptr", "nil"}}, false},
		{"Hidden Field", "hidden", reflect.StructField{Tag: reflect.StructTag(`hush:"hide"`)}, "secret", opts, [][]string{{"hidden", "HIDDEN"}}, false},
		{"Private Field", "private", reflect.StructField{PkgPath: "main"}, "private", opts, nil, false},
		{"Private Field Included", "private", reflect.StructField{Tag: reflect.StructTag(`hush:"mask"`), PkgPath: "main"}, "private", &hushOptions{includePrivate: true, maskFunc: defaultMaskFunc}, [][]string{{"private", "pr***te"}}, false},
		{"Float", "price", reflect.StructField{}, 3.14, opts, [][]string{{"price", "3.14"}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processField(context.Background(), tt.fieldName, tt.field, reflect.ValueOf(tt.value), tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("processField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessSimpleField(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		field     reflect.StructField
		value     interface{}
		hushTag   string
		opts      *hushOptions
		want      [][]string
		wantErr   bool
	}{
		{"Exported Int", "number", reflect.StructField{}, 42, "", &hushOptions{}, [][]string{{"number", "42"}}, false},
		{"Exported String", "name", reflect.StructField{}, "John", "", &hushOptions{}, [][]string{{"name", "John"}}, false},
		{"Exported Bool", "active", reflect.StructField{}, true, "", &hushOptions{}, [][]string{{"active", "true"}}, false},
		{"Exported Float", "price", reflect.StructField{}, 3.14, "hide", &hushOptions{}, [][]string{{"price", "HIDDEN"}}, false},
		{"Unexported Int", "age", reflect.StructField{PkgPath: "main"}, 30, "", &hushOptions{includePrivate: true}, [][]string{{"age", "30"}}, false},
		{"Unexported String", "secret", reflect.StructField{PkgPath: "main"}, "secret", "mask", &hushOptions{includePrivate: true, maskFunc: defaultMaskFunc}, [][]string{{"secret", "se**et"}}, false},
		{"Unexported Bool", "isAdmin", reflect.StructField{PkgPath: "main"}, false, "", &hushOptions{includePrivate: true}, [][]string{{"isAdmin", "false"}}, false},
		{"Unexported Float", "salary", reflect.StructField{PkgPath: "main"}, 50000.50, "mask", &hushOptions{includePrivate: true, maskFunc: defaultMaskFunc}, [][]string{{"salary", "50***.5"}}, false},
		{"Unexported Skipped", "skipped", reflect.StructField{PkgPath: "main"}, "skip me", "", &hushOptions{}, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processSimpleField(tt.fieldName, tt.field, reflect.ValueOf(tt.value), tt.hushTag, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("processSimpleField() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processSimpleField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessString(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		value     string
		hushTag   string
		opts      *hushOptions
		want      [][]string
	}{
		{
			name:      "No hush tag",
			fieldName: "field1",
			value:     "value1",
			hushTag:   "",
			opts:      &hushOptions{},
			want:      [][]string{{"field1", "value1"}},
		},
		{
			name:      "Mask tag",
			fieldName: "field2",
			value:     "sensitive",
			hushTag:   TagMask,
			opts:      &hushOptions{maskFunc: func(s string) string { return "***" }},
			want:      [][]string{{"field2", "***"}},
		},
		{
			name:      "Include private with mask func",
			fieldName: "field3",
			value:     "private",
			hushTag:   "",
			opts:      &hushOptions{includePrivate: true, maskFunc: func(s string) string { return "xxx" }},
			want:      [][]string{{"field3", "private"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := processString(tt.fieldName, tt.value, tt.hushTag, tt.opts)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertNonCompositeToString(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		want  string
	}{
		{"Bool", true, "true"},
		{"Int", 42, "42"},
		{"Float64", 3.14, "3.14"},
		{"String", "hello", "hello"},
		{"Complex128", complex(1, 2), "(1+2i)"},
		{"Unsupported", struct{}{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertNonCompositeToString(reflect.ValueOf(tt.value))
			if got != tt.want {
				t.Errorf("convertNonCompositeToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
