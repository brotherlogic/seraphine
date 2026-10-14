package workflow

import (
	"os"
	"path/filepath"
	"strings"
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

func TestWorkflowsRemoveLabelLast(t *testing.T) {
	workflows := []struct {
		file           string
		label          string
		subIssueMarker string
	}{
		{
			file:           "seraphine-decompose.md",
			label:          "seraphine-decompose",
			subIssueMarker: "gh issue create --parent",
		},
		{
			file:           "seraphine-needs-deep-research.md",
			label:          "seraphine-needs-deep-research",
			subIssueMarker: "native GitHub child sub-issue",
		},
		{
			file:           "seraphine-needs-requirements.md",
			label:          "seraphine-needs-requirements",
			subIssueMarker: "native GitHub sub-issue",
		},
		{
			file:           "seraphine-needs-implementation-plan.md",
			label:          "seraphine-needs-implementation-plan",
			subIssueMarker: "native GitHub sub-issue",
		},
		{
			file:           "seraphine-break-down-issue.md",
			label:          "seraphine-break-down-issue",
			subIssueMarker: "native GitHub sub-issue",
		},
		{
			file:           "seraphine-bug.md",
			label:          "seraphine-bug",
			subIssueMarker: "native GitHub sub-issue",
		},
	}

	for _, wf := range workflows {
		t.Run(wf.file, func(t *testing.T) {
			path := filepath.Join("..", "..", ".agent", "workflows", wf.file)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read %s: %v", path, err)
			}
			content := string(data)

			removeMarker := "Remove the `" + wf.label + "`"
			removeIdx := strings.LastIndex(content, removeMarker)
			if removeIdx == -1 {
				t.Fatalf("%s does not contain label removal instruction %q", wf.file, removeMarker)
			}

			subIssueIdx := strings.Index(content, wf.subIssueMarker)
			if subIssueIdx == -1 {
				t.Fatalf("%s does not contain sub-issue creation instruction %q", wf.file, wf.subIssueMarker)
			}

			if removeIdx < subIssueIdx {
				t.Errorf("%s removes label %q (at index %d) before creating sub-issue %q (at index %d); label removal must be the last step",
					wf.file, wf.label, removeIdx, wf.subIssueMarker, subIssueIdx)
			}

			// Also verify Mermaid diagram ordering if both "Remove Label" and "Create Sub-" appear
			mermaidStart := strings.Index(content, "```mermaid")
			if mermaidStart != -1 {
				mermaidEnd := strings.Index(content[mermaidStart+10:], "```")
				if mermaidEnd != -1 {
					mermaid := strings.ToLower(content[mermaidStart : mermaidStart+10+mermaidEnd])
					mRemove := strings.Index(mermaid, "remove label")
					mSubIssue := strings.Index(mermaid, "create sub-issue")
					if mSubIssue == -1 {
						mSubIssue = strings.Index(mermaid, "create breakdown sub-issue")
					}
					if mRemove != -1 && mSubIssue != -1 && mRemove < mSubIssue {
						t.Errorf("%s Mermaid diagram places 'Remove Label' before 'Create Sub-Issue'", wf.file)
					}
				}
			}
		})
	}
}

