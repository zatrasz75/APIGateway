package postgres

import (
	"censorship/pkg/storage"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
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

func (p *Store) AllList() ([]storage.Stop, error) {
	rows, err := p.db.Query(context.Background(), "SELECT * FROM stop")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []storage.Stop
	for rows.Next() {
		var c storage.Stop
		err = rows.Scan(&c.ID, &c.StopList)
		if err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (p Store) AddList(c storage.Stop) error {
	_, err := p.db.Exec(context.Background(),
		"INSERT INTO stop (stop_list) VALUES ($1);", c.StopList)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// CreateStopTable Создает таблицу
func (p *Store) CreateStopTable() error {
	_, err := p.db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS stop (
			id SERIAL PRIMARY KEY,
			stop_list TEXT
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

// DropStopTable Удаляет таблицу
func (p *Store) DropStopTable() error {
	_, err := p.db.Exec(context.Background(), "DROP TABLE IF EXISTS stop;")
	if err != nil {
		return err
	}
	return nil
}
