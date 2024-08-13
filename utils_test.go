package hush

import "testing"

func TestBuildFieldName(t *testing.T) {
	tests := []struct {
		name      string
		prefix    string
		fieldName string
		separator string
		want      string
	}{
		{"With prefix", "parent", "child", ".", "parent.child"},
		{"Without prefix", "", "field", ".", "field"},
		{"Custom separator", "parent", "child", "_", "parent_child"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildFieldName(tt.prefix, tt.fieldName, tt.separator); got != tt.want {
				t.Errorf("buildFieldName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultMaskFunc(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"Short string", "abc", "***"},
		{"Medium string", "password", "pa****rd"},
		{"Long string", "loooooooooooooooooooong", "lo*******************ng"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultMaskFunc(tt.value); got != tt.want {
				t.Errorf("defaultMaskFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}
