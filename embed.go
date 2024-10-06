package brpg

import (
	"embed"
	"io/fs"
)

//go:embed resources
var f embed.FS

//go:embed resources/font/AncientModernTales-a7Po.ttf
var MainFont []byte

func FS() embed.FS {
	return f
}
func init() {
	fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		//log.Debug().Msgf("path=%q, isDir=%v\n", path, d.IsDir())
		return nil
	})
}
