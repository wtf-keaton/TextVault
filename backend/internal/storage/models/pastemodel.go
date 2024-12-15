package models

type Paste struct {
	ID       int64
	Title    string
	Hash     string
	AuthorID int64
	Content  string
}
