#!/bin/bash

# migration
migrate -path "./migrations" -database "mysql://root:$DB_PWD@tcp(crawler_mysql)/webcrawler" up

# start server
go run main.go
