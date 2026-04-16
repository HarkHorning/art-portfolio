package models

type ArtModel struct {
	Id          int     `db:"id" json:"id"`
	Title       string  `db:"title" json:"title"`
	Description string  `db:"description" json:"description"`
	Portrait    bool    `db:"portrait" json:"portrait"`
	DisplayURL  string  `db:"display_url" json:"url"`
	MadeYear    *int    `db:"made_year" json:"made_year,omitempty"`
	Sold        bool    `db:"sold" json:"sold"`
	Size        *string `db:"size" json:"size,omitempty"`
	PriceCents  *int    `db:"price_cents" json:"price_cents,omitempty"`
}

type CategoryModel struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Slug string `db:"slug" json:"slug"`
}

type PrintSizeModel struct {
	Id              int    `db:"id" json:"id"`
	PrintId         int    `db:"print_id" json:"print_id"`
	Size            string `db:"size" json:"size"`
	PriceCents      int    `db:"price_cents" json:"price_cents"`
	QuantityInStock int    `db:"quantity_in_stock" json:"quantity_in_stock"`
	Sold            bool   `db:"sold" json:"sold"`
}

type PrintModel struct {
	Id          int              `db:"id" json:"id"`
	ArtTileId   int              `db:"art_tile_id" json:"art_tile_id"`
	Title       string           `db:"title" json:"title"`
	Description string           `db:"description" json:"description"`
	Portrait    bool             `db:"portrait" json:"portrait"`
	DisplayURL  string           `db:"display_url" json:"url"`
	Visible     bool             `db:"visible" json:"visible"`
	Sizes       []PrintSizeModel `json:"sizes"`
}

type ArtDetailModel struct {
	Id          int             `db:"id" json:"id"`
	Title       string          `db:"title" json:"title"`
	Description string          `db:"description" json:"description"`
	Portrait    bool            `db:"portrait" json:"portrait"`
	MadeYear    *int            `db:"made_year" json:"made_year,omitempty"`
	Sold        bool            `db:"sold" json:"sold"`
	Visible     bool            `db:"visible" json:"visible"`
	Size        *string         `db:"size" json:"size,omitempty"`
	PriceCents  *int            `db:"price_cents" json:"price_cents,omitempty"`
	Images      []ImageModel    `json:"images"`
	Categories  []CategoryModel `json:"categories"`
}
