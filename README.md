# TodoList API

## Описание

TodoList API - это простой RESTful API для управления задачами. Этот проект написан на языке Go и использует базу данных PostgreSQL. API позволяет создавать, получать, обновлять и удалять задачи, а также фильтровать и сортировать их.

## Возможности

- Создание задачи
- Получение задачи по ID
- Обновление задачи
- Удаление задачи
- Список задач с фильтрацией и пагинацией

## Технологии

- Go
- PostgreSQL
- Docker
- Chi Router
- Swagger для документирования API
- pressly/goose для миграций

  ## Запуск
  ```bash
  make compose-up
  make test-migration-up
  go run ./cmd/server
  ```


  ## Тесты
  Юнит тестами покрыты handler.go и service.go
  Для запуска
  ```bash
  make test
  ```

  ## Swagger
Для генерации документации Swagger:

``` bash
make swag
```
Документация будет доступна по адресу http://localhost:8080/swagger/index.html. 


