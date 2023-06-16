package main

import (
	"gateway/pkg/api"
	"github.com/gofrs/uuid"
	"log"
	"net/http"
)

// сервер
type server struct {
	api *api.API
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
		reqID := req.URL.Query().Get("request_id")
		if reqID == "" {
			rID, _ := uuid.NewV4()
			reqID = rID.String()
		}
		req.Header.Set("X-Request-ID", reqID)
		w.Header().Set("X-Request-ID", reqID)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, req)

		statusCode := lrw.statusCode
		log.Printf("<-- client ip: %s, method: %s, url: %s, status code: %d %s, trace id: %s",
			req.RemoteAddr, req.Method, req.URL.Path, statusCode, http.StatusText(statusCode), reqID)

	})
}

func main() {

	// объект сервера
	var srv server

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New()

	srv.api.Router().Use(Middle)

	log.Print("Запуск сервера http://127.0.0.1:8000/news")

	err := http.ListenAndServe(":8000", srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
	}

}
