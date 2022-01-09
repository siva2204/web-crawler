package db

type Url struct {
	KeyId int    `gorm:"column:keyId"`
	Url   string `grom:"column:url"`
}

func (Url) TableName() string {
	return "Url"
}
