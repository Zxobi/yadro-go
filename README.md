# yadro-go

Web-server for searching indexed comics data, aggregated from https://xkcd.com/.<br>
This is a study project from course "YADRO - Разработка микросервисных приложений на Golang".

---
## How to run

Following command will build an app and launch http-server, that will create SQLite database file and proceed migrations.

```shell
make
./xkcd-server [-port port] [-c path_to_config_file]
```
---
## Configuration options

List of available parameters for configuration file.

- source_url - address for comics fetching. Default is `"https://xkcd.com"`;
- dns - path to database file. Default is `"database.db"`;
- migrations - path to migrations directory. Default is `"migrations"`;
- port - port to start server. Default is `20202`;
- scheduler_hour/scheduler_minute - the hour and minute of the day to run automatic comics fetching. Default is `3:00`;
- parallel - maximum number of parallel comics fetch jobs. Default is `200`;
- fetch_limit - maximum number of comics to be fetched. Default is unlimited;
- scan_timeout - timeout for scanning indexed comics. Default is unlimited;
- scan_limit - maximum number of comics to return from scanning. Default is `10`;
- request_timeout - timeout for xkcd fetching requests. Default is unlimited;
- token_secret - secret for generated JWT token. Default is `"token-secret"`;
- token_max_time - JWT token ttl. Default is `1h`
- rate_limit - rps limit for search endpoint. Default is unlimited;
- concurrency_limit - concurrent requests limit for search endpoint. Default is unlimited.

---
## API Endpoints

### POST /login
Handles user login.<br>
Three users is created via migrations:
- user1:user1 (Role = User)
- user2:user2 (Role = User)
- admin:admin (Role = Admin)

#### Request Body
```json
{
    "username": "string",
    "password": "string"
}
```
#### Response
```json
{
  "token": "string"
}
```

### POST /update
Launch database update process.<br>
Available only for admin role user.

#### Headers
```Authorization: Bearer {token}```

#### Response
```json
{
  "total": 12345
}
```

### GET /pics
Search through indexed comics by query.<br>
Available for authorized users.<br>
Limited by rps and concurrent access.

#### Query Parameters
```?search="query sentence"```

#### Headers
```Authorization: Bearer {token}```

#### Response
```json
[
  "comic_url",
  "comic_url"
]
```