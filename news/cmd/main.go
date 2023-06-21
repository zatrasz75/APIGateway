package main

import (
	config "GoNews/configs"
	api "GoNews/pkg/api"
	"GoNews/pkg/middl"
	"GoNews/pkg/rss"
	storage "GoNews/pkg/storage"
	db "GoNews/pkg/storage/db"
	"context"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

// сервер GoNews
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

const (
	configURL = "./cmd/config.json"
)

func main() {

	// объект сервера
	var srv server

	cfg := config.New()
	// Адрес базы данных
	dbURL := cfg.News.URLdb
	// Порт по умолчанию.
	port := cfg.News.AdrPort
	// Можно сменить Порт при запуске флагом < --news-port= >
	portFlag := flag.String("news-port", port, "Порт для news сервиса")
	flag.Parse()
	portNews := *portFlag

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// объект базы данных postgresql
	db, err := db.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Удаление таблицы gonews если существует
	err = db.DropGonewsTable()
	if err != nil {
		log.Println(err)
		return
	}
	// Создание таблицы gonews если не существует
	err = db.CreateGonewsTable()
	if err != nil {
		log.Println(err)
		return
	}

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	//--------------------------------------------------------

	// каналы для обработки новостей и ошибок
	chanPosts := make(chan []storage.Post)
	chanErrs := make(chan error)

	// Чтение RSS-лент из конфига с заданным интервалом
	go func() {
		err := rss.GoNews(configURL, chanPosts, chanErrs)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// запись публикаций в db
	go func() {
		for posts := range chanPosts {
			if err := srv.db.PostsCreation(posts); err != nil {
				chanErrs <- err
			}
		}
	}()

	// вывод ошибок
	go func() {
		for err := range chanErrs {
			log.Println(err)
		}
	}()

	srv.api.Router().Use(middl.Middle)

	log.Print("Запуск сервера на http://127.0.0.1" + portNews)

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(portNews, srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер. Ошибка:", err)
	}

}
