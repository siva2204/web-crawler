#!/bin/bash

declare -i maxtry=3
while [ $maxtry -gt 0 ]; do
    nc -z crawler_mysql 3306
    isopen=$?
    if [ $isopen -eq 0 ]; then
        break
    fi
    maxtry=${maxtry}-1
    sleep 1
done

# migration
migrate -path "./migrations" -database "mysql://root:$DB_PWD@tcp(crawler_mysql)/webcrawler" up

# start server
go run main.go
