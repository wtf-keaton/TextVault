package models

type Paste struct {
	ID       int64  `db:"id"`
	Title    string `db:"title"`
	Hash     string `db:"hash"`
	AuthorID int64  `db:"authorid"`
}
