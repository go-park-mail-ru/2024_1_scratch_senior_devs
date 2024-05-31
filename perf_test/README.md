# Домашнее задание №3

Основная сущность нашего приложения - заметка. Мы провели нагрузочное тестирование: создали 100000 заметок, затем прочитали 100000 заметок, потом оценили результаты. Для проведения нагрузочного тестирования мы выбрали https://github.com/wg/wrk.

## Создание заметок

Для генерации заметок мы написали скрипт [create.lua](./create.lua):
```lua
wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjUzMTcwOTg5NjYsImlkIjoiYzM2YjQzODAtNmNjMS00MmU4LWFkZmMtY2RmNTQ1ZWUxZmUzIiwidXNyIjoiZWxhc3RpYyJ9.1XcLZh_O_WfqEKrdudqQCbWM7XholfkYIYb4Sz2-b_E"
wrk.headers["X-Csrf-Token"] = "2e495587-8347-4c0a-8041-d092df5a6ffa"
wrk.headers["Cookie"] = "YouNoteJWT=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjUzMTcwOTg5NjYsImlkIjoiYzM2YjQzODAtNmNjMS00MmU4LWFkZmMtY2RmNTQ1ZWUxZmUzIiwidXNyIjoiZWxhc3RpYyJ9.1XcLZh_O_WfqEKrdudqQCbWM7XholfkYIYb4Sz2-b_E; Path=/; Secure; HttpOnly; Expires=Mon, 29 May 2124 19:01:36 GMT;"

local counter = 0

request = function()
    counter = counter + 1

    local title = "title of my note " .. counter
    local content = "my text of my note " .. counter
    local body = string.format('{"data": {"title": "%s", "content": "%s"}}', title, content)
    return wrk.format(nil, "/api/note/add", nil, body)
end
```

Запускали вот такой командой:
```bash
wrk -t4 -c100 -d10000s --script=create.lua --latency http://127.0.0.1:8080
```

Вывод команды:
![img.png](images/img.png)

## Чтение заметок

Для чтения заметок мы написали скрипт [read.lua](./read.lua):

```lua
wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjUzMTcwOTg5NjYsImlkIjoiYzM2YjQzODAtNmNjMS00MmU4LWFkZmMtY2RmNTQ1ZWUxZmUzIiwidXNyIjoiZWxhc3RpYyJ9.1XcLZh_O_WfqEKrdudqQCbWM7XholfkYIYb4Sz2-b_E"
wrk.headers["X-Csrf-Token"] = "2e495587-8347-4c0a-8041-d092df5a6ffa"
wrk.headers["Cookie"] = "YouNoteJWT=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjUzMTcwOTg5NjYsImlkIjoiYzM2YjQzODAtNmNjMS00MmU4LWFkZmMtY2RmNTQ1ZWUxZmUzIiwidXNyIjoiZWxhc3RpYyJ9.1XcLZh_O_WfqEKrdudqQCbWM7XholfkYIYb4Sz2-b_E; Path=/; Secure; HttpOnly; Expires=Mon, 29 May 2124 19:01:36 GMT;"

request = function()
    return wrk.format(nil, "/api/note/756d7250-1460-4314-9850-82da0f09cfff", nil, nil)
end

```

Запускали вот такой командой:
```bash
wrk -t4 -c100 -d10000s --script=read.lua --latency http://127.0.0.1:8080
```

Вывод команды:
![img_1.png](images/img_1.png)

## Анализ результатов

При создании заметок получили около 150 rps - это небольшой показатель, но и не плохой.

При чтении заметок получили также около 150 rps.

## Оптимизации

Можно оптимизтировать запрос на получение заметки: можно денормализовать схему данных, избавившись от таблицы избранного. Вместо таблицы избранного можно будет сохранять эти данные внутри каждой заметки. Это поможет избавиться от JOIN в запросе.
