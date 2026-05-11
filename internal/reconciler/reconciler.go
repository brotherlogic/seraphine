package reconciler

import (
	"os"
	"path/filepath"

	pb "github.com/brotherlogic/seraphine/proto"
)

func Reconcile(files []*pb.File) error {
	for _, file := range files {
		dir := filepath.Dir(file.Path)
		if dir != "." {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
		}

		err := os.WriteFile(file.Path, file.Content, 0644)
		if err != nil {
			return err
		}
	}
	return nil
}
