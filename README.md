## web-crawler

### Prerequisites

- Go 1.16+
- MySQL
- redis [download link](https://redis.io/download)
- [air](https://github.com/cosmtrek/air) (for live reload) [installation link](https://github.com/cosmtrek/air#prefer-installsh)

### Migrations

- Install

```bash
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

- Create Migrations

```bash
migrate create -ext sql -dir ./migrations <MIGRATION_NAME>
```

- Up Migrations

```bash
migrate -path "./migrations" -database "mysql://root:YOUR_MYSQL_PASSWORD@/webcrawler" up
```

- Down Migrations

```bash
migarte -path "./migrations"  -database "mysql://root:YOUR_MYSQL_PASSWORD@/webcrawler" down
```

- start server

  - development

    ```bash
    air
    ```

  - production
    ```bash
    go run main.go
    ```