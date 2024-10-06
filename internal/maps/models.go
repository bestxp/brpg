package maps

import (
	"encoding/xml"
	"image"
)

type Tileset struct {
	Data     *TilesetData
	Text     string `xml:",chardata" json:"text,omitempty"`
	Firstgid string `xml:"firstgid,attr" json:"firstgid,omitempty"`
	Source   string `xml:"source,attr" json:"source,omitempty"`
}

type Layer struct {
	Text    string `xml:",chardata" json:"text,omitempty"`
	ID      string `xml:"id,attr" json:"id,omitempty"`
	Name    string `xml:"name,attr" json:"name,omitempty"`
	Width   string `xml:"width,attr" json:"width,omitempty"`
	Height  string `xml:"height,attr" json:"height,omitempty"`
	Visible string `xml:"visible,attr" json:"visible,omitempty"`
	Data    struct {
		Text     string `xml:",chardata" json:"text,omitempty"`
		Encoding string `xml:"encoding,attr" json:"encoding,omitempty"`
	} `xml:"data" json:"data,omitempty"`
}

type Object struct {
	Text    string `xml:",chardata" json:"text,omitempty"`
	ID      string `xml:"id,attr" json:"id,omitempty"`
	Gid     string `xml:"gid,attr" json:"gid,omitempty"`
	X       string `xml:"x,attr" json:"x,omitempty"`
	Y       string `xml:"y,attr" json:"y,omitempty"`
	Width   string `xml:"width,attr" json:"width,omitempty"`
	Height  string `xml:"height,attr" json:"height,omitempty"`
	Polygon *struct {
		Text   string `xml:",chardata" json:"text,omitempty"`
		Points string `xml:"points,attr" json:"points,omitempty"`
	} `xml:"polygon" json:"polygon,omitempty"`
}

type ObjectGroup struct {
	Text   string   `xml:",chardata" json:"text,omitempty"`
	ID     string   `xml:"id,attr" json:"id,omitempty"`
	Name   string   `xml:"name,attr" json:"name,omitempty"`
	Object []Object `xml:"object" json:"object,omitempty"`
}

type TiledMap struct {
	XMLName      xml.Name    `xml:"map" json:"map,omitempty"`
	Text         string      `xml:",chardata" json:"text,omitempty"`
	Version      string      `xml:"version,attr" json:"version,omitempty"`
	Tiledversion string      `xml:"tiledversion,attr" json:"tiledversion,omitempty"`
	Orientation  string      `xml:"orientation,attr" json:"orientation,omitempty"`
	Renderorder  string      `xml:"renderorder,attr" json:"renderorder,omitempty"`
	Width        int         `xml:"width,attr" json:"width,omitempty"`
	Height       int         `xml:"height,attr" json:"height,omitempty"`
	Tilewidth    int         `xml:"tilewidth,attr" json:"tilewidth,omitempty"`
	Tileheight   int         `xml:"tileheight,attr" json:"tileheight,omitempty"`
	Infinite     string      `xml:"infinite,attr" json:"infinite,omitempty"`
	Nextlayerid  string      `xml:"nextlayerid,attr" json:"nextlayerid,omitempty"`
	Nextobjectid string      `xml:"nextobjectid,attr" json:"nextobjectid,omitempty"`
	Tileset      []Tileset   `xml:"tileset" json:"tileset,omitempty"`
	Layer        []Layer     `xml:"layer" json:"layer,omitempty"`
	Objectgroup  ObjectGroup `xml:"objectgroup" json:"objectgroup,omitempty"`
}

type TileAnimationFrame struct {
	Text     string `xml:",chardata" json:"text,omitempty"`
	Tileid   string `xml:"tileid,attr" json:"tileid,omitempty"`
	Duration string `xml:"duration,attr" json:"duration,omitempty"`
}

type TileAnimation struct {
	Text  string               `xml:",chardata" json:"text,omitempty"`
	Frame []TileAnimationFrame `xml:"frame" json:"frame,omitempty"`
}

type Tile struct {
	Text      string        `xml:",chardata" json:"text,omitempty"`
	ID        string        `xml:"id,attr" json:"id,omitempty"`
	Animation TileAnimation `xml:"animation" json:"animation,omitempty"`
	Image     *TilesetImage `xml:"image" json:"image,omitempty"`
}

type TilesetDataImageImage struct {
	Img image.Image
}

type TilesetImage struct {
	Image  *TilesetDataImageImage
	Text   string `xml:",chardata" json:"text,omitempty"`
	Source string `xml:"source,attr" json:"source,omitempty"`
	Width  string `xml:"width,attr" json:"width,omitempty"`
	Height string `xml:"height,attr" json:"height,omitempty"`
}

type TilesetDataGrid struct {
	Text        string `xml:",chardata" json:"text,omitempty"`
	Orientation string `xml:"orientation,attr" json:"orientation,omitempty"`
	Width       string `xml:"width,attr" json:"width,omitempty"`
	Height      string `xml:"height,attr" json:"height,omitempty"`
}

type TilesetData struct {
	XMLName      xml.Name         `xml:"tileset" json:"tileset,omitempty"`
	Text         string           `xml:",chardata" json:"text,omitempty"`
	Version      string           `xml:"version,attr" json:"version,omitempty"`
	Tiledversion string           `xml:"tiledversion,attr" json:"tiledversion,omitempty"`
	Name         string           `xml:"name,attr" json:"name,omitempty"`
	Tilewidth    string           `xml:"tilewidth,attr" json:"tilewidth,omitempty"`
	Tileheight   string           `xml:"tileheight,attr" json:"tileheight,omitempty"`
	Tilecount    string           `xml:"tilecount,attr" json:"tilecount,omitempty"`
	Columns      string           `xml:"columns,attr" json:"columns,omitempty"`
	Image        *TilesetImage    `xml:"image" json:"image,omitempty"`
	Tile         []Tile           `xml:"tile" json:"tile,omitempty"`
	Grid         *TilesetDataGrid `xml:"grid" json:"grid,omitempty"`
}
