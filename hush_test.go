package hush

import (
	"context"
	"reflect"
	"testing"
)

type nestedStruct struct {
	NestedField string `hush:"hide"`
}

type testStruct struct {
	PublicField  string `hush:"mask"`
	privateField string `hush:"mask"`
	NestedStruct nestedStruct
	IntField     int
	BoolField    bool
	FloatField   float64
}

func TestNewHush(t *testing.T) {
	h := NewHush()

	if h == nil {
		t.Error("NewHush() returned nil, want non-nil")
	}

	_, ok := h.(Husher)
	if !ok {
		t.Error("NewHush() did not return a Husher interface")
	}
}

func TestHushType_Hush(t *testing.T) {
	tests := []struct {
		name           string
		input          interface{}
		options        []interface{}
		want           [][]string
		wantErr        bool
		wantErrMessage string
	}{
		{
			name: "Simple struct without private fields",
			input: testStruct{
				PublicField:  "sensitive",
				privateField: "private",
				NestedStruct: nestedStruct{NestedField: "nested"},
				IntField:     42,
				BoolField:    true,
				FloatField:   3.14,
			},
			want: [][]string{
				{"BoolField", "true"},
				{"FloatField", "3.14"},
				{"IntField", "42"},
				{"NestedStruct.NestedField", "HIDDEN"},
				{"PublicField", "se*****ve"},
			},
			wantErr: false,
		},
		{
			name: "Simple struct with private fields",
			input: testStruct{
				PublicField:  "sensitive",
				privateField: "private",
				NestedStruct: nestedStruct{NestedField: "nested"},
				IntField:     42,
				BoolField:    true,
				FloatField:   3.14,
			},
			options: []interface{}{WithPrivateFields(true)},
			want: [][]string{
				{"BoolField", "true"},
				{"FloatField", "3.14"},
				{"IntField", "42"},
				{"NestedStruct.NestedField", "HIDDEN"},
				{"PublicField", "se*****ve"},
				{"privateField", "pr***te"},
			},
			wantErr: false,
		},
		{
			name:    "String",
			input:   "sensitive",
			options: []interface{}{TagMask},
			want:    [][]string{{"se*****ve"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHush()

			got, err := h.Hush(context.Background(), tt.input, tt.options...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Hush() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if err == nil || err.Error() != tt.wantErrMessage {
					t.Errorf("Hush() error message = %v, want %v", err, tt.wantErrMessage)
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hush() = %v, want %v", got, tt.want)
			}
		})
	}
}
