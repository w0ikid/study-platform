uml-диаграмма {login | register} (https://www.websequencediagrams.com/cgi-bin/cdraw?lz=c2VxdWVuY2VEaWFncmFtCiAgICBwYXJ0aWNpcGFudCDQmtC70LjQtdC90YIACxKh0LXRgNCy0LXRgCBhcwACDihHaW4pAD4SkdCUACYFkdCw0LfQsCDQtNCw0L3QvdGL0YUgKFBvc3RncmVTUUwANQYAgQsFJSUg0KDQtdCz0LjRgdGC0YDQsNGG0LjRjyDQv9C-0LvRjNC30L7QstCw0YLQtdC70Y8AgUIFAIEvDC0-PgCBIQw6IFBPU1QgL2FwaS9hdXRoL3JlZ2lzdGVyIHt1c2VybmFtZSwgZW1haWwsIHBhc3N3b3JkLCByb2xlfQCCHQVOb3RlIG92ZXIAgXYNOiDQktCwAIIrBbQAgRkKAIFTDVNob3VsZEJpbmRKU09OAIFZC2FsdCDQndC10LrQvtGA0YDQtdC60YLQvdGL0LUAghkL0LUAgx0FAIFaBQCCcgstLT4-AIMeDDogNDAwIEJhZCBSZXF1ZXN0AINTBWVsc2Ug0JQAQgogAFwSAEkWPj7QkdCUOiDQn9GAAIJlBbXRgNC40YLRjCDRg9C90LjQutCwAIMEBb3QvgCDHwWMAII5BiDQuCAAgkoIAIExCQCBOglhbHQgRQAgB9C7ACALINGD0LbQtSDRgdGD0YnQtdGB0YLQstGD0LUAhH8HAIF8CZHQlACBeAUAgnkOntGI0LjQsdC60LAAgQoX0LgAOg4Agi4fOSBDb25mbGljAIJEBgCCRQkAgTYIAIEtDQCCAw0AgkgLAIQdHaXQtdGI0LgAglkHsACCUAW1INC_0LAAgm4FAIVNCACDAR6h0L7RhdGAAIZNBbUAOQoAhgoaAIIpH5LQvtC30LIAhm0FidCw0LXRgiBJRAAsJwCFAx0yMDEgQ3JlYXRlZCB7aWQsAIQdCQCGbwlyb2xlLCBsZXZlbCwgeHAAhnMGAIVCBW5kAIVJBgABBwCIGwmQ0LLRgtC-0YDQuNC3AIdSUmxvZ2luIHsAiAgPAIYRggDQvtC40YHQugCKcxkAiyAFAIpHBgCHbBYAPwUAizMSjCDQvdC1INC90LDQudC00LXQvQCHVyQAMB0AiC4iAIo2HzEgVW5hdXRob3JpemUAhQEGAIpQCgCBLRsAgRYtAIsVDQCHGiQAjHYZAIsjDLrQsACINRoAikgNAI0DCQCLYga90YvQuQCIfgyMAIp2DQCBfDsAgi8MAFQJIABrDABZEQCOYRmTAIlsBbUAkBILSldUINGC0L7QugCKBwWwICh1c2VyLklEAIhnBi5Sb2xlAJEbBgCJCywwIE9LIHt0b2tlbjogSldUAIkOCgCJEAwAiRgQ&s=default)
```
sequenceDiagram
    participant Клиент
    participant Сервер as Сервер (Gin)
    participant БД as База данных (PostgreSQL)
    
    %% Регистрация пользователя
    Клиент->>Сервер: POST /api/auth/register {username, email, password, role}
    Note over Сервер: Валидация данных (ShouldBindJSON)
    
    alt Некорректные данные
        Сервер-->>Клиент: 400 Bad Request
    else Данные корректны
        Сервер->>БД: Проверить уникальность email и username
        
        alt Email или username уже существует
            БД-->>Сервер: Ошибка уникальности
            Сервер-->>Клиент: 409 Conflict
        else Email и username уникальны
            Note over Сервер: Хеширование пароля
            Сервер->>БД: Сохранение пользователя
            БД-->>Сервер: Возвращает ID пользователя
            Сервер-->>Клиент: 201 Created {id, username, email, role, level, xp}
        end
    end
    
    %% Авторизация пользователя
    Клиент->>Сервер: POST /api/auth/login {email, password}
    Note over Сервер: Валидация данных (ShouldBindJSON)
    
    alt Некорректные данные
        Сервер-->>Клиент: 400 Bad Request
    else Данные корректны
        Сервер->>БД: Поиск пользователя по email
        
        alt Пользователь не найден
            БД-->>Сервер: Пользователь не существует
            Сервер-->>Клиент: 401 Unauthorized
        else Пользователь найден
            БД-->>Сервер: Данные пользователя
            Note over Сервер: Проверка пароля
            
            alt Неверный пароль
                Сервер-->>Клиент: 401 Unauthorized
            else Пароль верный
                Note over Сервер: Генерация JWT токена (user.ID, user.Role)
                Сервер-->>Клиент: 200 OK {token: JWT}
            end
        end
    end
```    
erd-диаграмма link (https://dbdiagram.io/d/67f7ce354f7afba1841dcc66)
```
// Educational Platform ERD Diagram

Table users {
  id integer [primary key]
  username text [unique, not null]
  email text [unique, not null]
  password text [not null]
  role text [not null]
  xp integer [default: 0]
  level integer [default: 1]
  created_at timestamp [default: `CURRENT_TIMESTAMP`]
  updated_at timestamp [default: `CURRENT_TIMESTAMP`]
}

Table courses {
  id integer [primary key]
  name text [not null]
  description text
  teacher_id integer [not null, ref: > users.id]
  created_at timestamp [default: `CURRENT_TIMESTAMP`]
  updated_at timestamp [default: `CURRENT_TIMESTAMP`]
  status text [default: 'active', note: 'Can be active, inactive, completed']
}

Table lessons {
  id integer [primary key]
  course_id integer [not null, ref: > courses.id]
  title text [not null]
  content text [not null]
  video_url text
  created_at timestamp [default: `CURRENT_TIMESTAMP`]
  updated_at timestamp [default: `CURRENT_TIMESTAMP`]
}

Table enrollments {
  id integer [primary key]
  user_id integer [not null, ref: > users.id]
  course_id integer [not null, ref: > courses.id]
  created_at timestamp [default: `CURRENT_TIMESTAMP`]
  updated_at timestamp [default: `CURRENT_TIMESTAMP`]
  
  indexes {
    (user_id, course_id) [unique]
  }
}

Table certificates {
  id integer [primary key]
  user_id integer [not null, ref: > users.id]
  course_id integer [not null, ref: > courses.id]
  issued_at timestamp [default: `CURRENT_TIMESTAMP`]
  
  indexes {
    (user_id, course_id) [unique, note: 'Ensures one certificate per course']
  }
}
```