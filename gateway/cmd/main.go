package main

import (
	"flag"
	config "gateway/configs"
	"gateway/pkg/api"
	"gateway/pkg/middl"
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

	// Конфигурация
	cfg := config.New()

	// Порт по умолчанию.
	port := cfg.Gateway.AdrPort
	// Порт по умолчанию.
	newsPort := cfg.News.AdrPort
	// Порт по умолчанию.
	censorPort := cfg.Censor.AdrPort
	// Порт по умолчанию.
	comment := cfg.Comments.AdrPort

	// Можно сменить Порт при запуске флагом < --gateway-port= >
	portFlag := flag.String("gateway-port", port, "Порт для gateway сервиса")

	// Можно сменить Порт при запуске флагом < --news-port= >
	portFlagNews := flag.String("news-port", newsPort, "Порт для news сервиса")

	// Можно сменить Порт при запуске флагом < --censor-port= >
	portFlagCensor := flag.String("censor-port", censorPort, "Порт для censor сервиса")

	// Можно сменить Порт при запуске флагом < --comments-port= >
	portFlagComment := flag.String("comments-port", comment, "Порт для comments сервиса")

	flag.Parse()

	portGateway := *portFlag
	portNews := *portFlagNews
	portCensor := *portFlagCensor
	portComment := *portFlagComment

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(cfg, portNews, portCensor, portComment)

	srv.api.Router().Use(middl.Middle)

	log.Print("Запуск сервера http://127.0.0.1" + portGateway + "/news")

	err := http.ListenAndServe(portGateway, srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
	}

}
