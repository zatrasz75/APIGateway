package db

// Post Публикация, получаемая из RSS.
type Post struct {
	ID      int    `json:"ID,omitempty"`      // Номер записи
	Title   string `json:"title,omitempty"`   // Заголовок публикации
	Content string `json:"content,omitempty"` // Содержание публикации
	PubTime int64  `json:"pubTime,omitempty"` // Время публикации
	Link    string `json:"link,omitempty"`    // Ссылка на источник
}

type Pagination struct {
	NumOfPages int `json:"numOfPages,omitempty"`
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
}

// Interface задаёт контракт на работу с БД
type Interface interface {
	Posts(limit, offset int) ([]Post, error)                                       // Получение n-ого кол-ва публикаций
	AddPost(p Post) error                                                          // Добавление новой публикации в базу
	PostSearchILIKE(keyWord string, limit, offset int) ([]Post, Pagination, error) // Поиск по заголовку
	PostsCreation([]Post) error                                                    // Создание n-ого кол-ва публикаций
	PostDetal(id int) (Post, error)                                                // Детальный вывод
	CreateGonewsTable() error
	DropGonewsTable() error
}
