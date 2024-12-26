package models

type Paste struct {
	ID       string `db:"id"`
	Title    string `db:"title"`
	Language string `db:"language"`
	AuthorID int64  `db:"authorid"`
}
