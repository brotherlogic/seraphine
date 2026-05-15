package workflow

import (
	"testing"
)

func TestExtractProjectName(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://github.com/brotherlogic/seraphine.git", "brotherlogic/seraphine"},
		{"https://github.com/brotherlogic/seraphine", "brotherlogic/seraphine"},
		{"git@github.com:brotherlogic/seraphine.git", "brotherlogic/seraphine"},
		{"git@github.com:brotherlogic/seraphine", "brotherlogic/seraphine"},
	}

	for _, tt := range tests {
		actual, err := extractProjectName(tt.url)
		if err != nil {
			t.Errorf("extractProjectName(%s) returned error: %v", tt.url, err)
			continue
		}
		if actual != tt.expected {
			t.Errorf("extractProjectName(%s) = %s; want %s", tt.url, actual, tt.expected)
		}
	}
}
