package main

import (
	"gateway/pkg/api"
	"gateway/pkg/middl"
	"log"
	"net/http"
)

// сервер
type server struct {
	api *api.API
}

func main() {

	// объект сервера
	var srv server

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New()

	srv.api.Router().Use(middl.Middle)

	log.Print("Запуск сервера http://127.0.0.1:8000/news")

	err := http.ListenAndServe(":8000", srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
	}

}
