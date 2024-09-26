package brpg

import (
	"embed"
)

//go:embed resources/*
var f embed.FS

//go:embed resources/font/AncientModernTales-a7Po.ttf
var MainFont []byte

func FS() embed.FS {
	return f
}
