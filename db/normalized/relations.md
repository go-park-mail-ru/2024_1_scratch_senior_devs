# Структура базы данных

## Схема данных
```mermaid
erDiagram
    profile {
        uuid id PK
        text description
        text username
        text password_hash
        timestamp create_time
        text image_path
        text secret
    }
    note {
        uuid id PK
        text data
        timestamp create_time
        timestamp update_time
        uuid owner_id FK
    }
    attach {
        uuid id PK
        text path
        uuid note_id FK
    }
    profile ||--o{ note : "1:M"
    note ||--o{ attach : "1:M"
```

## Описание таблиц
### profile
Таблица profile содержит данные пользователей:<br/>
id - идентификатор пользователя<br/>
description - описание в профиле пользователя<br/>
username - логин пользователя, а также его имя в сервисе<br/>
password_hash - хэш пароля пользователя<br/>
create_time - дата и время регистрации<br/>
image_path - путь до файла аватарки пользователя<br/>
secret - секрет для генерации QR-кода для двухфакторной аутентификации<br/>

### note
Таблица note отвечает за хранение заметок:<br/>
id - идентификатор заметки<br/>
data - содержимое заметки<br/>
create_time - дата и время создания заметки<br/>
update_time - дата и время последнего изменения заметки<br/>
owner_id - идентификатор пользователя, который является создателем заметки<br/>

### attach
В таблице attach хранятся сведения о вложениях заметок:<br/>
id - идентификатор вложения<br/>
path - путь до файла на сервере<br/>
note_id - идентификатор заметки, к которой это вложение прикреплено<br/>

## Нормализация
### Функциональные зависимости:
profile:<br/>
{id} -> description, username, password_hash, create_time, image_path, secret<br/>
{username} -> id, description, password_hash, create_time, image_path, secret<br/>
{password_hash} -> id, description, username, create_time, image_path, secret<br/>
{create_time} -> id, description, username, password_hash, image_path, secret<br/>
{image_path} -> id, description, username, password_hash, create_time, secret<br/>
{secret} -> id, description, username, password_hash, create_time, image_path<br/>

note:<br/>
{id} -> data, create_time, update_time, owner_id<br/>
{owner_id} -> id, data, create_time, update_time<br/>
{create_time} -> id, data, update_time, owner_id<br/>
{update_time} -> id, data, create_time, owner_id<br/>

attach:
{id} -> path, note_id<br/>
{note_id} -> id, path<br/>

### Проверка нормальных форм:
Первая нормальная форма (1NF):<br/>
В схеме каждый атрибут является атомарным, так что она соответствует 1NF.<br/>

Вторая нормальная форма (2NF):<br/>
В схеме каждый неключевой атрибут зависит от всего первичного ключа, поэтому она соответствует 2NF.<br/>

Третья нормальная форма (3NF):<br/>
В схеме нет транзитивных зависимостей, так как каждый атрибут функционально зависит только от ключа и не от других атрибутов.<br/>

Нормальная форма Бойса-Кодда (BCNF):<br/>
В схеме все функциональные зависимости либо тривиальны, либо ключевые, поэтому она соответствует BCNF.<br/>
