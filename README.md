# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.




## Описание проекта

Проект создан для упрощённого планирования задач с REST API и статическим веб-интерфейсом. Планировщик задач может:

- создавать задачи с датой, временем и правилом повторения;
- получать следующую дату выполнения задачи через `/api/nextdate`;
- искать задачи по тексту и дате;
- редактировать и удалять задачи;
- отмечать задачи как выполненные.





## Выполненные задачи со звёздочкой

-  реализован API `/api/nextdate` для расчёта следующей даты по правилу повторения;
-  реализован POST `/api/task` для создания задачи;
-  реализован GET `/api/tasks` для списка задач с фильтрацией и поиском;
-  реализован GET `/api/task?id=<id>` для получения задачи по идентификатору;
-  реализован PUT `/api/task` для редактирования задачи;
-  реализован POST `/api/task/done` для отметки задачи выполненной;
-  реализован DELETE `/api/task?id=<id>` для удаления задачи;
-  добавлен поиск по `search` и фильтр по дате в `/api/tasks`;
-  реализована простая аутентификация через `TODO_PASSWORD` и выдача JWT-токена;
-  добавлен Dockerfile для сборки и запуска контейнера.





## Локальный запуск

1. Задать переменные окружения.

Windows CMD:

   bash
set TODO_PASSWORD=12345
set TODO_PORT=7540
set TODO_DBFILE=scheduler.db


Windows PowerShell:

   powershell
$env:TODO_PASSWORD = "12345"
$env:TODO_PORT = "7540"
$env:TODO_DBFILE = "scheduler.db"


Linux / Git Bash:

   bash
export TODO_PASSWORD=12345
export TODO_PORT=7540
export TODO_DBFILE=scheduler.db


2. Запуск сервера:

   bash
go run main.go

3. Адрес в браузере:

http://localhost:7540/login.html


Если пароль не задан, можно сразу открыть:

http://localhost:7540/









## Запуск тестов

В `tests/settings.go` можно задать параметры для тестов:

- `Port` — порт тестового сервера;
- `DBFile` — путь до тестовой или общей базы данных;
- `FullNextDate` — `true`, если нужно проверить расширенные правила `/api/nextdate`;
- `Search` — `true`, если нужно проверить поиск по `search`;
- `Token` — JWT-токен для авторизации, если используется `TODO_PASSWORD`.

Запустить все тесты:

go test ./tests

Запустить конкретный тест:

go test ./tests -run ^TestAddTask$






## Docker


### Сборка образа
   bash
docker build -t todo-planner .


### Запуск контейнера

```bash
docker run -d --name todo-planner \
  -p 7540:7540 \
  -v /f/goooooo/Dev/FINAL/go_final_project/scheduler.db:/data/scheduler.db \
  -e TODO_PASSWORD=12345 \
  todo-planner
```

Если вы запускаете Docker из Git Bash/WSL на Windows, путь может требовать двойного слэша:

```bash
docker run -d --name todo-planner \
  -p 7540:7540 \
  -v //f/goooooo/Dev/FINAL/go_final_project/scheduler.db:/data/scheduler.db \
  -e TODO_PASSWORD=12345 \
  todo-planner
```

> Замените путь к файлу `scheduler.db` на реальный путь на вашей машине.
