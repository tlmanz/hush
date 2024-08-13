package hush

import (
	"testing"
)

func TestNewHush(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{"Struct", struct{ Name string }{"John"}, false},
		{"String", "test", false},
		{"Integer", 42, true},
		{"Slice", []int{1, 2, 3}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHush(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHush() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("NewHush() returned nil, want non-nil")
			}
		})
	}
}
