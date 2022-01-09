package db

import (
	"fmt"

	"github.com/siva2204/web-crawler/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// mysql connection object
var DB *gorm.DB

func initDB() {
	dbName := config.Getenv("DB_NAME")
	dbPwd := config.Getenv("DB_PWD")
	dbUser := config.Getenv("DB_USER")
	dbHost := config.Getenv("DB_HOST")
	dbPort := config.Getenv("DB_PORT")

	// db connection str
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", dbUser, dbPwd, dbHost, dbPort, dbName)

	// connecting to db
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("error connecting DB, %+v", err))
	}

	DB = db
}
