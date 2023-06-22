package postgres

import (
	"censorship/pkg/storage"
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

func TestStore_AddList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dataBase, err := New(ctx, "postgres://postgres:rootroot@localhost:5432/comm")
	str := storage.Stop{
		StopList: "ups",
	}
	dataBase.AddList(str)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Создана запись.")
}

func TestStore_AllList(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	dataBase, err := New(ctx, "postgres://postgres:rootroot@localhost:5432/comm")
	if err != nil {
		t.Fatal(err)
	}

	result, err := dataBase.AllList()
	if err != nil {
		t.Fatal(err)
	}

	// Проверка непустой таблицы
	if len(result) == 0 {
		t.Errorf("Таблица \"стоп\" пуста.")
	} else {
		t.Logf("Таблица \"стоп\" содержит %d записи.", len(result))
	}

}
