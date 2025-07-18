Database tables for anon_sessions
--- 
`anon_sessions`
```sql
id              UUID (PK)
session_code    TEXT UNIQUE
creator_token   TEXT (secret, not exposed in session)
display_name    TEXT
is_active       BOOLEAN
created_at      TIMESTAMP
```

`anon_participants`
```sql
id              UUID (PK)
session_id      UUID (FK → anon_sessions.id)
display_name    TEXT
joined_at       TIMESTAMP
```

`anon_snippets`
```sql 
id              UUID (PK)
session_id      UUID (FK → anon_sessions.id)
file_name       TEXT
content         TEXT
pushed_at       TIMESTAMP
```





---- 


# Our endpoints for the anon_sessions module: 
```| Endpoint                               | Method | Auth          | Purpose                 |
| -------------------------------------- | ------ | ------------- | ----------------------- |
| `/api/v1/sessions/anon`                | `POST` | None          | Create a new session    |
| `/api/v1/sessions/anon/:code/join`     | `POST` | None          | Join as a participant   |
| `/api/v1/sessions/anon/:code/snippets` | `POST` | Creator token | Push a snippet          |
| `/api/v1/sessions/anon/:code/snippets` | `GET`  | None          | Fetch all live snippets |
| `/api/v1/sessions/anon/:code/end`      | `POST` | Creator token | End the session         |
```
