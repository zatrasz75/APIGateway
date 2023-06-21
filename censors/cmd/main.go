package main

import (
	config "censorship/configs"
	"censorship/pkg/api"
	"censorship/pkg/middl"
	"censorship/pkg/storage"
	"censorship/pkg/storage/postgres"
	"censorship/pkg/supply"
	"context"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

// сервер
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
	port := cfg.Censor.AdrPort
	// Можно сменить Порт при запуске флагом < --censor-port= >
	portFlag := flag.String("censor-port", port, "Порт для censor сервиса")
	flag.Parse()
	portCensor := *portFlag

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// объект базы данных postgresql
	db, err := postgres.New(ctx, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// Удаление таблицы stop если существует
	err = db.DropStopTable()
	if err != nil {
		log.Println(err)
		return
	}
	// Создание таблицы stop если не существует
	err = db.CreateStopTable()
	if err != nil {
		log.Println(err)
		return
	}
	// Получение списка для стоп листа из файла words.txt
	stop, err := supply.StopList()
	if err != nil {
		log.Println(err)
	}
	// Добавление в таблицу stop полученного списка
	for _, v := range stop {
		err := db.AddList(v)
		if err != nil {
			log.Println(err)
		}
	}

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	srv.api.Router().Use(middl.Middle)

	log.Print("Запуск сервера на http://127.0.0.1" + portCensor)

	err = http.ListenAndServe(portCensor, srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
	}

}
