package web_test

import (
	"io"
	"io/fs"
	"testing"

	"github.com/brotherlogic/seraphine/web"
)

func TestGetStaticFS_NotNull(t *testing.T) {
	staticFS := web.GetStaticFS()
	if staticFS == nil {
		t.Fatalf("Expected non-nil fs.FS from GetStaticFS()")
	}
}

func TestGetStaticFS_ServesIndexHTML(t *testing.T) {
	staticFS := web.GetStaticFS()
	if staticFS == nil {
		t.Fatalf("Expected non-nil fs.FS from GetStaticFS()")
	}

	f, err := staticFS.Open("index.html")
	if err != nil {
		t.Fatalf("Failed to open index.html from staticFS: %v", err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("Failed to read index.html: %v", err)
	}

	if len(content) == 0 {
		t.Errorf("Expected non-empty index.html content")
	}

	stat, err := f.Stat()
	if err != nil {
		t.Fatalf("Failed to stat index.html: %v", err)
	}
	if stat.IsDir() {
		t.Errorf("Expected index.html to be a file, not a directory")
	}
}

func TestGetStaticFS_ImplementsFS(t *testing.T) {
	var _ fs.FS = web.GetStaticFS()
}
