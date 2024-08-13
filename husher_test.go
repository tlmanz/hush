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

func TestHushType_Hush(t *testing.T) {
	tests := []struct {
		name           string
		input          interface{}
		options        []Option
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
			options: []Option{WithPrivateFields(true)},
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
			want:    [][]string{{"", "se*****ve"}},
			wantErr: false,
		},
		{
			name:           "Invalid input",
			input:          42,
			wantErr:        true,
			wantErrMessage: "expected struct or string, got int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := NewHush(tt.input)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("NewHush() unexpected error = %v", err)
				}
				if err.Error() != tt.wantErrMessage {
					t.Fatalf("NewHush() error = %v, wantErrMessage %v", err, tt.wantErrMessage)
				}
				return
			}

			got, err := h.Hush(context.Background(), "", tt.options...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Hush() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hush() = %v, want %v", got, tt.want)
			}
		})
	}
}
