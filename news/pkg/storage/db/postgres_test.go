package db

import (
	db "GoNews/pkg/storage"
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := New(ctx, "postgres://postgres:rootroot@localhost:5432/aggregator")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStore_AddPost(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dataBase, err := New(ctx, "postgres://postgres:rootroot@localhost:5432/aggregator")
	post := db.Post{
		Title:   "тестирования",
		Content: "Пробный текст",
		PubTime: 5,
		Link:    "Линка",
	}
	dataBase.AddPost(post)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Создана запись.")
}
