package db

type IndexRelation struct {
	KeyId int `gorm:"column:keyId"`
	UrlId int `gorm:"column:urlId"`
}

func (IndexRelation) TableName() string {
	return "IndexRelation"
}
