package yaml

type Texture struct {
	entities   map[string]*Entity
	animations map[string]*Animation
	images     map[string]*Image

	Entities   []*Entity    `yaml:"entity"`
	Animations []*Animation `yaml:"animations"`
	Images     []*Image     `yaml:"images"`
}

type Entity struct {
	Name    string         `yaml:"name"`
	Size    Size           `yaml:"size"`
	Actions []EntityAction `yaml:"actions"`
}

type Size struct {
	W int `yaml:"w"`
	H int `yaml:"h"`
}

type EntityAction struct {
	Name   string `yaml:"name"`
	Sprite string `yaml:"sprite"`
}

type Animation struct {
	Name     string   `yaml:"name"`
	Frames   []string `yaml:"frames"` // array of path to image or image
	Duration int      `yaml:"duration"`
	Loop     bool     `yaml:"loop"`
	Flip     string   `yaml:"flip,omitempty"`
}

type Image struct {
	Name   string      `yaml:"name"`
	Source string      `yaml:"source"`
	Tiles  *ImageTiles `yaml:"tiles,omitempty"`
}

type ImageTiles struct {
	Tile Size `yaml:"tile"`
	Cols int  `yaml:"cols"`
	Rows int  `yaml:"rows"`
}
