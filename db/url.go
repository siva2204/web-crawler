package db

type Url struct {
	Id   int     `gorm:"column:id"`
	Url  string  `gorm:"column:url"`
	Rank float64 `gorm:"column:rank"`
}

func (Url) TableName() string {
	return "Url"
}
