package maps

type Map struct {
	tiled TiledMap

	tiles TileRepository
}

func NewMap(t TiledMap) *Map {
	return &Map{tiled: t}
}

func (m *Map) Size() (int, int) {
	return m.tiled.Width * m.tiled.Tilewidth, m.tiled.Height * m.tiled.Tileheight
}
