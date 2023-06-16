package db

import (
	storage "GoNews/pkg/storage"
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

// PostsCreation Создание n-ого кол-ва публикаций
func (p *Store) PostsCreation(posts []storage.Post) error {
	for _, post := range posts {
		err := p.AddPost(post)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddPost запись новых данных в БД
func (s *Store) AddPost(p storage.Post) error {

	err := s.db.QueryRow(context.Background(), `
		INSERT INTO gonews (title, content, pubtime, link)
		VALUES ($1, $2, $3, $4);
		`,
		p.Title,
		p.Content,
		p.PubTime,
		p.Link,
	).Scan()
	return err
}

// PostSearchILIKE Поиск по заголовку
func (p *Store) PostSearchILIKE(pattern string, limit, offset int) ([]storage.Post, storage.Pagination, error) {
	pattern = "%" + pattern + "%"

	pagination := storage.Pagination{
		Page:  offset/limit + 1,
		Limit: limit,
	}
	row := p.db.QueryRow(context.Background(), "SELECT count(*) FROM gonews WHERE title ILIKE $1;", pattern)
	err := row.Scan(&pagination.NumOfPages)

	if pagination.NumOfPages%limit > 0 {
		pagination.NumOfPages = pagination.NumOfPages/limit + 1
	} else {
		pagination.NumOfPages /= limit
	}

	if err != nil {
		return nil, storage.Pagination{}, err
	}

	rows, err := p.db.Query(context.Background(), "SELECT * FROM gonews WHERE title ILIKE $1 ORDER BY pubtime DESC LIMIT $2 OFFSET $3;", pattern, limit, offset)
	if err != nil {
		return nil, storage.Pagination{}, err
	}
	defer rows.Close()
	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(&p.ID, &p.Title, &p.Content, &p.PubTime, &p.Link)
		if err != nil {
			return nil, storage.Pagination{}, err
		}
		posts = append(posts, p)
	}
	return posts, pagination, rows.Err()
}

// Posts Получение странице с определенным номером
func (s *Store) Posts(limit, offset int) ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT * FROM gonews
	ORDER BY pubtime DESC LIMIT $1 OFFSET $2
	`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.Post
	// итерированное по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, p)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return posts, rows.Err()
}

// PostDetal Получение публикаций по id
func (p *Store) PostDetal(id int) (storage.Post, error) {
	row := p.db.QueryRow(context.Background(), `
	SELECT * FROM gonews 
    WHERE id =$1;
	`, id)
	var post storage.Post
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.PubTime,
		&post.Link)
	if err != nil {
		return storage.Post{}, err
	}
	return post, nil
}
