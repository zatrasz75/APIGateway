package main

import (
	config "censorship/configs"
	"censorship/pkg/api"
	"censorship/pkg/middl"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

// сервер
type server struct {
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

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New()

	srv.api.Router().Use(middl.Middle)

	cfg := config.New()
	// Порт по умолчанию.
	port := cfg.Censor.AdrPort
	// Можно сменить Порт при запуске флагом < --censor-port= >
	portFlag := flag.String("censor-port", port, "Порт для censor сервиса")
	flag.Parse()
	portCensor := *portFlag

	log.Print("Запуск сервера на http://127.0.0.1" + portCensor)

	err := http.ListenAndServe(portCensor, srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
	}

}
