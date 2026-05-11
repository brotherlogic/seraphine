package reconciler

import (
	"os"
	"testing"

	pb "github.com/brotherlogic/seraphine/proto"
)

func TestReconcile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "seraphine-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to temp dir
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(oldWd)

	files := []*pb.File{
		{Path: "test.txt", Content: []byte("hello")},
		{Path: "subdir/subtest.txt", Content: []byte("world")},
	}

	err = Reconcile(files)
	if err != nil {
		t.Errorf("Reconcile failed: %v", err)
	}

	// Verify files
	content, err := os.ReadFile("test.txt")
	if err != nil || string(content) != "hello" {
		t.Errorf("test.txt content mismatch: %v, %s", err, string(content))
	}

	content, err = os.ReadFile("subdir/subtest.txt")
	if err != nil || string(content) != "world" {
		t.Errorf("subdir/subtest.txt content mismatch: %v, %s", err, string(content))
	}
}
