package models

type ArtModel struct { // remember to add categories
	Id          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Portrait    bool   `db:"portrait" json:"portrait"`
	URL         string `db:"url_low" json:"url"`
}
