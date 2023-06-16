package rss

import (
	storage "GoNews/pkg/storage"
	"encoding/json"
	"encoding/xml"
	strip "github.com/grokify/html-strip-tags-go"
	"log"
	"net/http"
	"os"
	"time"
)

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Content string `xml:"description"`
	PubDate string `xml:"pubDate"`
}

type Channel struct {
	Items []Item `xml:"channel>item"`
}

type config struct {
	Rss           []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}

// GoNews Чтение RSS-лент из конфига
func GoNews(configURL string, chanPosts chan<- []storage.Post, chanErrs chan<- error) error {
	//чтение конфига
	file, err := os.Open(configURL)
	if err != nil {
		return err
	}
	var conf config
	err = json.NewDecoder(file).Decode(&conf)
	if err != nil {
		return err
	}

	log.Println("начинаю смотреть rss-каналы")
	//запуск горутины для каждой rss-ленты
	for i, r := range conf.Rss {
		go func(r string, i int, chanPosts chan<- []storage.Post, chanErrs chan<- error) {
			for {
				log.Println("запустил  goroutine", i, "по ссылке", r)
				p, err := GetRss(r)
				if err != nil {
					chanErrs <- err
					continue
				}
				chanPosts <- p
				log.Println("insert posts from goroutine", i, "по ссылке", r)
				log.Println("Goroutine ", i, ": ожидание следующей итерации")
				time.Sleep(time.Duration(conf.RequestPeriod) * time.Second * 10)
			}
		}(r, i, chanPosts, chanErrs)
	}
	return nil
}

// GetRss Выгрузка RSS-ленты по заданному URL
func GetRss(url string) ([]storage.Post, error) {
	var c Channel
	//запрос к rss-ленте
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	err = xml.NewDecoder(res.Body).Decode(&c)
	if err != nil {
		return nil, err
	}

	//преобразование данных из rss в список публикаций
	var posts []storage.Post
	for _, i := range c.Items {
		var p storage.Post
		p.Title = i.Title
		p.Content = i.Content
		p.Content = strip.StripTags(p.Content)
		p.Link = i.Link

		t, err := time.Parse(time.RFC1123, i.PubDate)
		if err != nil {
			t, err = time.Parse(time.RFC1123Z, i.PubDate)
		}
		if err != nil {
			t, err = time.Parse("Mon, _2 Jan 2006 15:04:05 -0700", i.PubDate)
		}
		if err == nil {
			p.PubTime = t.Unix()
		}

		posts = append(posts, p)
	}
	return posts, nil
}
