package main

import (
	config "comments/configs"
	"comments/pkg/api"
	"comments/pkg/middl"
	"comments/pkg/storage"
	"comments/pkg/storage/postgres"
	"context"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

// сервер comment
type server struct {
	db  storage.Interface
	api *api.API
}

// init вызывается перед main()
func init() {
	// загружает значения из файла .env в систему
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	// объект сервера
	var srv server

	cfg := config.New()
	// Адрес базы данных
	dbURL := cfg.Comments.URLdb
	// Порт по умолчанию.
	port := cfg.Comments.AdrPort
	// Можно сменить Порт при запуске флагом < --comments-port= >
	portFlag := flag.String("comments-port", port, "Порт для comments сервиса")
	flag.Parse()
	portComments := *portFlag

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// объект базы данных postgresql
	db, err := postgres.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	srv.api.Router().Use(middl.Middle)

	log.Print("Запуск сервера на http://127.0.0.1" + portComments)

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(portComments, srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер. Ошибка:", err)
	}
}
