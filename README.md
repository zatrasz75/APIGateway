# APIGateway
Skillfactory

### запустить приложение:
#### Соединение со свей базой данных PostgreSql и порты можно отредактировать в файлах .env 

в каждой из директорий (censors comments news gateway) выполнить команду 
* go run cmd/main.go

а так же можно изменить порты при запуске флагами
* go run cmd/main.go --gateway-port= --censor-port= --news-port= --comments-port=

или с помощью Makefile с портами по умолчанию
* make all
#### gateway запускается на localhost:8000

### Доступные API , примеры:

постраничная навигация
* http://localhost:8000/news?page=2&s=

вывод последних новостей
* http://localhost:8000/news/latest

вывод по номеру страннице
* http://localhost:8000/news/latest?page=2

поиск по заголовкам
* http://localhost:8000/news?s=gRPC

детальная информация о посте с комментарием
* http://localhost:8000/news/search?id=1

добавление комментария методом post в формате JSON , 
с проверкой на слова из стоп листа (qwerty , йцукен , zxvbnm)
* http://localhost:8000/comments/add

Удаляет комментарий по id методом delete в формате JSON
* http://localhost:8000/comments/del