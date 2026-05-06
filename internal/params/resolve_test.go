package params

import (
	"testing"
)

func TestResolveParameters(t *testing.T) {
	tests := []struct {
		name      string
		paramDefs map[string]string
		cliValues map[string]string
		expected  map[string]string
	}{
		{
			name:      "defaults only",
			paramDefs: map[string]string{"age": "25", "name": "Alice"},
			cliValues: nil,
			expected:  map[string]string{"age": "25", "name": "Alice"},
		},
		{
			name:      "CLI overrides defaults",
			paramDefs: map[string]string{"age": "25"},
			cliValues: map[string]string{"age": "30"},
			expected:  map[string]string{"age": "30"},
		},
		{
			name:      "unknown CLI values ignored",
			paramDefs: map[string]string{"age": "25"},
			cliValues: map[string]string{"age": "30", "foo": "bar"},
			expected:  map[string]string{"age": "30"},
		},
		{
			name:      "empty both",
			paramDefs: nil,
			cliValues: nil,
			expected:  map[string]string{},
		},
		{
			name:      "empty defs nonempty CLI",
			paramDefs: nil,
			cliValues: map[string]string{"x": "1"},
			expected:  map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveParameters(tt.paramDefs, tt.cliValues)
			if len(got) != len(tt.expected) {
				t.Fatalf("got %d params, want %d", len(got), len(tt.expected))
			}
			for k, v := range tt.expected {
				if got[k] != v {
					t.Errorf("got[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}

func TestGetMissingRequired(t *testing.T) {
	tests := []struct {
		name         string
		paramDefs    map[string]string
		currentVals  map[string]string
		expectedLen  int
		expectedName string
	}{
		{
			name:        "all have defaults",
			paramDefs:   map[string]string{"age": "25"},
			currentVals: nil,
			expectedLen: 0,
		},
		{
			name:         "required with no value",
			paramDefs:    map[string]string{"name": ""},
			currentVals:  nil,
			expectedLen:  1,
			expectedName: "name",
		},
		{
			name:        "required with value provided",
			paramDefs:   map[string]string{"name": ""},
			currentVals: map[string]string{"name": "Alice"},
			expectedLen: 0,
		},
		{
			name:        "empty defs",
			paramDefs:   nil,
			currentVals: nil,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMissingRequired(tt.paramDefs, tt.currentVals)
			if len(got) != tt.expectedLen {
				t.Fatalf("got %d missing, want %d: %v", len(got), tt.expectedLen, got)
			}
			if tt.expectedName != "" && len(got) > 0 && got[0] != tt.expectedName {
				t.Errorf("got %q, want %q", got[0], tt.expectedName)
			}
		})
	}
}

func TestValidateCLIValues(t *testing.T) {
	tests := []struct {
		name      string
		cliValues map[string]string
		paramDefs map[string]string
		wantErr   bool
	}{
		{
			name:      "all valid",
			cliValues: map[string]string{"age": "30"},
			paramDefs: map[string]string{"age": "25"},
			wantErr:   false,
		},
		{
			name:      "unknown param",
			cliValues: map[string]string{"foo": "bar"},
			paramDefs: map[string]string{"age": "25"},
			wantErr:   true,
		},
		{
			name:      "empty CLI values",
			cliValues: nil,
			paramDefs: map[string]string{"age": "25"},
			wantErr:   false,
		},
		{
			name:      "empty defs nonempty CLI",
			cliValues: map[string]string{"x": "1"},
			paramDefs: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCLIValues(tt.cliValues, tt.paramDefs)
			if tt.wantErr && err == nil {
				t.Error("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestValidateParamNames(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]string
		wantErr bool
	}{
		{
			name:    "normal names",
			params:  map[string]string{"user_id": "1", "age": "25"},
			wantErr: false,
		},
		{
			name:    "reserved edit",
			params:  map[string]string{"edit": "val"},
			wantErr: true,
		},
		{
			name:    "reserved format",
			params:  map[string]string{"format": "csv"},
			wantErr: true,
		},
		{
			name:    "reserved f",
			params:  map[string]string{"f": "json"},
			wantErr: true,
		},
		{
			name:    "reserved last",
			params:  map[string]string{"last": "true"},
			wantErr: true,
		},
		{
			name:    "empty params",
			params:  nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateParamNames(tt.params)
			if tt.wantErr && err == nil {
				t.Error("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
