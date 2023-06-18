package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"regexp"
)

// API приложения.
type API struct {
	r *mux.Router // Маршрутизатор запросов
}

// New Конструктор API.
func New() *API {
	api := API{
		r: mux.NewRouter(),
	}
	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	api.r.HandleFunc("/comments/check", api.addCommentHandler).Methods(http.MethodPost, http.MethodOptions)
}

func (api *API) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	text := struct {
		Content string
	}{}
	err := json.NewDecoder(r.Body).Decode(&text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	stoplist := []string{
		"qwerty",
		"йцукен",
		"zxvbnm",
	}

	for _, stopWord := range stoplist {
		matched, err := regexp.MatchString(stopWord, text.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if matched {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
