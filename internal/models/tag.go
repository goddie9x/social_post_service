package models

type Tag struct {
	Name string  `json:"name" gorm:"primaryKey"`
	Post []*Post `gorm:"many2many:post_tags;"`
}

type PostTag struct {
	PostId  string `gorm:"type:varchar2(36);index:idx_post_tag;"`
	TagName string `gorm:"type:varchar2(1024);index:idx_post_tag;"`
}
