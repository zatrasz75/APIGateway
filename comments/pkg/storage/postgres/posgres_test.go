package postgres

import (
	"comments/pkg/storage"
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	_, err := New(ctx, "postgres://postgres:rootroot@localhost:5432/comm")
	if err != nil {
		t.Fatal(err)
	}
}

func TestStore_AddComment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dataBase, err := New(ctx, "postgres://postgres:rootroot@localhost:5432/comm")
	comment := storage.Comment{
		NewsID:  2,
		Content: "Текст проверки",
	}
	dataBase.AddComment(comment)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Создана запись.")
}

func TestStore_DeleteComment(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dataBase, err := New(ctx, "postgres://postgres:rootroot@localhost:5432/comm")
	comment := storage.Comment{
		ID: 1,
	}
	dataBase.DeleteComment(comment)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Удалена запись.")
}
