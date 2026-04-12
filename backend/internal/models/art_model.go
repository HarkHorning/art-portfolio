package models

type ArtModel struct {
	Id          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Portrait    bool   `db:"portrait" json:"portrait"`
	URL         string `db:"url_low" json:"url"`
}

type CategoryModel struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Slug string `db:"slug" json:"slug"`
}

type ArtDetailModel struct {
	Id          int            `db:"id" json:"id"`
	Title       string         `db:"title" json:"title"`
	Description string         `db:"description" json:"description"`
	Portrait    bool           `db:"portrait" json:"portrait"`
	URL         string         `db:"url_low" json:"url"`
	Categories  []CategoryModel `json:"categories"`
}
