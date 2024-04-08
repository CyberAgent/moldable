package str

import (
	"os"
	"path"
)

func BuildPath(file string, wd *string) string {
	mockdir := os.Getenv("MOCKDIR")
	return path.Join(*wd, mockdir, file)
}
