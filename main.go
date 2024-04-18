package main

import (
	"os"

	"github.com/cyberagent-oss/moldable/src/cmd"
)

var (
  pkgVersion = ""
)

func main() {
  os.Setenv("PKG_VERSION", pkgVersion)
	cmd.Execute()
}
