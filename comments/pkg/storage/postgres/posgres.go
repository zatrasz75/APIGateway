package postgres

import (
	"comments/pkg/storage"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Store Хранилище данных
type Store struct {
	db *pgxpool.Pool
}

// New Конструктор объекта хранилища
func New(ctx context.Context, constr string) (*Store, error) {

	for {
		_, err := pgxpool.Connect(ctx, constr)
		if err == nil {
			break
		}
	}
	db, err := pgxpool.Connect(ctx, constr)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

// AllComments выводит все коменты.
func (p *Store) AllComments(newsID int) ([]storage.Comment, error) {
	rows, err := p.db.Query(context.Background(), "SELECT * FROM comments WHERE news_id = $1;", newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []storage.Comment
	for rows.Next() {
		var c storage.Comment
		err = rows.Scan(&c.ID, &c.NewsID, &c.Content, &c.PubTime)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

// AddComment добавляет коменты.
func (p *Store) AddComment(c storage.Comment) error {
	_, err := p.db.Exec(context.Background(),
		"INSERT INTO comments (news_id,content) VALUES ($1,$2);", c.NewsID, c.Content)
	if err != nil {
		return err
	}
	return nil
}

// DeleteComment удаляет коменты.
func (p *Store) DeleteComment(c storage.Comment) error {
	_, err := p.db.Exec(context.Background(),
		"DELETE FROM comments WHERE id=$1;", c.ID)
	if err != nil {
		return err
	}
	return nil
}
