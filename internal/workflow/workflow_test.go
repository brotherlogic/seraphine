package workflow

import (
	"fmt"
	"os"
	"os/exec"
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

func mockExecCommand(output string, err error) func(command string, args ...string) *exec.Cmd {
	return func(command string, args ...string) *exec.Cmd {
		// Mock implementation using a dummy script or something...
		// Actually, standard Go way is to use a helper process
		// But simpler here: we just can't easily return output with an exec.Cmd unless we run a helper
		cs := []string{"-test.run=TestHelperProcess", "--", command}
		cs = append(cs, args...)
		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1", "MOCK_OUTPUT="+output)
		if err != nil {
			cmd.Env = append(cmd.Env, "MOCK_ERROR=1")
		}
		return cmd
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	if os.Getenv("MOCK_ERROR") == "1" {
		os.Exit(1)
	}
	fmt.Fprint(os.Stdout, os.Getenv("MOCK_OUTPUT"))
	os.Exit(0)
}

func TestValidateRuleset(t *testing.T) {
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	tests := []struct {
		name      string
		output    string
		mockErr   error
		wantValid bool
		wantErr   bool
	}{
		{"empty list", "[]\n", nil, false, false},
		{"has rulesets", "[{\"id\": 1}]\n", nil, true, false},
		{"error running gh", "", fmt.Errorf("error"), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execCommand = mockExecCommand(tt.output, tt.mockErr)
			valid, err := validateRuleset("brotherlogic/seraphine")
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRuleset() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if valid != tt.wantValid {
				t.Errorf("validateRuleset() valid = %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

func TestFileRulesetIssue(t *testing.T) {
	origExecCommand := execCommand
	defer func() { execCommand = origExecCommand }()

	tests := []struct {
		name    string
		mockErr error
		wantErr bool
	}{
		{"success", nil, false},
		{"failure", fmt.Errorf("error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execCommand = mockExecCommand("", tt.mockErr)
			err := fileRulesetIssue("brotherlogic/seraphine")
			if (err != nil) != tt.wantErr {
				t.Errorf("fileRulesetIssue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
