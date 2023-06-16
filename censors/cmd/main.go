package main

import (
	"censorship/pkg/api"
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

	srv.api.Router().Use(Middle)

	log.Print("Запуск сервера на http://127.0.0.1:8083")

	err := http.ListenAndServe(censorAddr, srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
	}

}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func Middle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqID := req.Header.Get("X-Request-ID")

		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, req)

		statusCode := lrw.statusCode
		log.Printf("<-- client ip: %s, method: %s, url: %s, status code: %d %s, trace id: %s",
			req.RemoteAddr, req.Method, req.URL.Path, statusCode, http.StatusText(statusCode), reqID)

	})
}
