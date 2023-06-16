package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
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

type ResponseDetailed struct {
	NewsDetailed struct {
		ID      int    `json:"ID"`
		Title   string `json:"Title"`
		Content string `json:"Content"`
		PubTime int    `json:"PubTime"`
		Link    string `json:"Link"`
	} `json:"NewsDetailed"`
	Comments []struct {
		ID              int    `json:"ID"`
		NewsID          int    `json:"newsID"`
		ParentCommentID int    `json:"parentCommentID"`
		Content         string `json:"content"`
		PubTime         int    `json:"pubTime"`
	} `json:"Comments"`
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	api.r.HandleFunc("/news", api.newsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/latest", api.newsLatestHandler).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/news/search", api.newsDetailedHandler).Methods(http.MethodGet, http.MethodOptions)
	api.r.HandleFunc("/comments/add", api.addCommentHandler).Methods(http.MethodPost, http.MethodOptions)
	api.r.HandleFunc("/comments/del", api.delCommentHandler).Methods(http.MethodDelete, http.MethodOptions)

}

func (api *API) newsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/news" {
		http.NotFound(w, r)
	}

	// Создаем прокси-сервер для первого микросервиса
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:8081", // адрес первого микросервиса
	})

	// Определяем, к какому микросервису нужно проксировать запрос
	if r.URL.Path == "/news" {
		// Проксируем запрос к первому микросервису
		proxy.ServeHTTP(w, r)
	}

}

func (api *API) newsLatestHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/news/latest" {
		http.NotFound(w, r)
	}

	// Создаем прокси-сервер для первого микросервиса
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:8081", // адрес первого микросервиса
	})

	// Определяем, к какому микросервису нужно проксировать запрос
	if r.URL.Path == "/news/latest" {
		// Проксируем запрос к первому микросервису
		proxy.ServeHTTP(w, r)
	}

}

func (api *API) newsDetailedHandler(w http.ResponseWriter, r *http.Request) {

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		http.Error(w, "параметры поиска обязательны", http.StatusBadRequest)
		return
	}

	chNews := make(chan *http.Response, 2)
	chComments := make(chan *http.Response, 2)
	chErr := make(chan error, 2)
	var response ResponseDetailed
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// Отправляем запрос на порт 8081
		resp1, err := http.Get("http://localhost:8081/news/search" + "?id=" + idParam)
		chErr <- err
		chNews <- resp1
	}()

	go func() {
		defer wg.Done()
		// Отправляем запрос на порт 8082
		resp2, err := http.Get("http://localhost:8082/comments" + "?news_id=" + idParam)
		chErr <- err
		chComments <- resp2
	}()

	wg.Wait()
	close(chErr)

	for err := range chErr {
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

block:
	for {
		select {
		case respNews := <-chNews:
			body, err := ioutil.ReadAll(respNews.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = json.Unmarshal(body, &response.NewsDetailed)
		case respComment := <-chComments:
			body, err := ioutil.ReadAll(respComment.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			_ = json.Unmarshal(body, &response.Comments)
		default:
			break block
		}

	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (api *API) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comments/add" {
		http.NotFound(w, r)
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	Body1 := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	Body := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	respCensor, err := MakeRequest(r, http.MethodPost, "http://localhost:8083/comments/add", Body1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if respCensor.StatusCode != 200 {
		http.Error(w, "неправильное содержание комментария", respCensor.StatusCode)
		return
	}

	resp, err := MakeRequest(r, http.MethodPost, "http://localhost:8082/comments/add", Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for name, values := range resp.Header {
		w.Header()[name] = values
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (api *API) delCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comments/del" {
		http.NotFound(w, r)
	}

	// Создаем прокси-сервер для первого микросервиса
	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:8082", // адрес микросервиса
	})

	// Определяем, к какому микросервису нужно проксировать запрос
	if r.URL.Path == "/comments/del" {
		// Проксируем запрос к первому микросервису
		proxy.ServeHTTP(w, r)
	}
}

func MakeRequest(req *http.Request, method, url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	r.Header = req.Header
	return client.Do(r)
}
