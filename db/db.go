package db

import (
	"fmt"

	"github.com/siva2204/web-crawler/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// mysql connection object
var DB *gorm.DB

func InitDB() {
	dbName := config.Config.DbName
	dbPwd := config.Config.DbPassword
	dbUser := config.Config.DbUser
	dbHost := config.Config.DbHost
	dbPort := config.Config.DbPort

	// db connection str
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", dbUser, dbPwd, dbHost, dbPort, dbName)

	// connecting to db
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("error connecting DB, %+v", err))
	}

	DB = db
}

// this func persist index from redis to mysql
func PersistIndex(key []string, values map[string][]string) {
	n := len(key)

	for i := 0; i < n; i++ {
		k := key[i]

		var keyS Key

		// create key if not present
		if err := DB.Where("`key` = ?", k).First(&keyS).Error; err != nil {
			fmt.Errorf("Key not found %+v", err)

			keyS.Key = k

			if err := DB.Create(&keyS).Error; err != nil {
				fmt.Errorf("Error in creating key")
			}
		}

		// saving all the urls
		for _, ll := range values[k] {
			var newUrl Url

			newUrl.KeyId = keyS.Id
			newUrl.Url = ll

			if err := DB.Create(newUrl).Error; err != nil {
				fmt.Errorf("Error creating url in db %+v", err)
			}
		}
	}
}
