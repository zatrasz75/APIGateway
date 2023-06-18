package main

import (
	"censorship/pkg/api"
	"censorship/pkg/middl"
	"log"
	"net/http"
)

// сервер
type server struct {
	api *api.API
}

const (
	censorAddr = ":8083"
)

func main() {

	// объект сервера
	var srv server

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New()

	srv.api.Router().Use(middl.Middle)

	log.Print("Запуск сервера на http://127.0.0.1:8083")

	err := http.ListenAndServe(censorAddr, srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
	}

}
