package models

type ImageModel struct {
	Id         int    `db:"id" json:"id"`
	ArtTileId  int    `db:"art_tile_id" json:"art_tile_id"`
	Variant    string `db:"variant" json:"variant"`
	URL        string `db:"url" json:"url"`
	Filename   string `db:"filename" json:"filename"`
	SortOrder  int    `db:"sort_order" json:"sort_order"`
}
