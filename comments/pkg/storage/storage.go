package storage

type Comment struct {
	ID      int    `json:"ID,omitempty"`
	NewsID  int    `json:"newsID,omitempty"`
	Content string `json:"content,omitempty"`
	PubTime int64  `json:"pubTime,omitempty"`
}

type Interface interface {
	AllComments(newsID int) ([]Comment, error)
	AddComment(Comment) error
	DeleteComment(Comment) error
	CreateCommentTable() error
	DropCommentTable() error
}
