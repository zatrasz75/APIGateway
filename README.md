# APIGateway
Skillfactory

### запустить приложение:

в каждой из директорий (censors comments news gateway) выполнить команду 
* go run cmd/main.go
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