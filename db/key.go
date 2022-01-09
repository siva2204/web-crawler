package db

type Key struct {
	Id  int    `gorm:"column:id;primary_key"`
	Key string `gorm:"column:key;unique"`
}

func (Key) TableName() string {
	return "Key"
}
