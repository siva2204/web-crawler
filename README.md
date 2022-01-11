## web-crawler

### [Presentation Link](https://docs.google.com/presentation/d/1AcLnB_yMLxfJI1PYAtjSAdepspRR0RSYdPitVqSltl8/edit#slide=id.p)

### Prerequisites

- Go 1.16+
- MySQL
- redis [download link](https://redis.io/download) v5.0.7
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
    go run main.go --threads <no of threads>
    ```

### Docker Setup
 - Install docker and docker-compose
 - Run `cp .docker.sample.env .docker.env`, fill the env varaiables
 - Run `docker-compose up` to start the services
