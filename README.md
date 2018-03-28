# Noah

> Useless & danger!

Migration manager for DB

* MySQL
* PostgreSQL

## Usage

Product is very raw! Please do not use it if you find this repo.

### Running in HTTP-server mode

```bash
./noah -s --port=12000
```

Migrate:

`POST localhost:12000`

```json
{
  "config": {
    "driver": "postgres",
    "password": 12345,
    "database": "wwm_media",
    "port": 5432,
    "reconnect_timeout": 5,
    "reconnect_attempts": 5
  },
  "path": "./migrations"
}
```

## Build

```bash
GOOS=linux GOARCH=amd64 go build -o noah_linux_amd64
```