package db

type Url struct {
	Id  int    `gorm:"column:id"`
	Url string `gorm:"column:url"`
}

func (Url) TableName() string {
	return "Url"
}
