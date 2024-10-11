package front

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var front embed.FS

func Assets() (fs.FS, error) {
	return fs.Sub(front, "dist")
}
