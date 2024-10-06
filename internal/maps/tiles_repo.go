package maps

type TileRepository map[string]any

func NewTileRepository() TileRepository {
	return make(TileRepository)
}
