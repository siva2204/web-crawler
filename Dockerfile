FROM golang:1.17

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    netcat \
    && rm -rf /var/lib/apt/lists/*

WORKDIR  /app

RUN  go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod go.sum ./

RUN go mod download 

COPY . .

CMD bash docker-entry.sh
