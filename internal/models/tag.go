package models

type Tag struct {
	Name string  `json:"name" gorm:"primaryKey"`
	Post []*Post `gorm:"many2many:post_tags;"`
}
