package models

import "github.com/google/uuid"

type Tag struct {
	PostId uuid.UUID `json:"postId" gorm:"foreignKey:postId,index"`
	Name   string    `json:"name" gorm:"index"`
}
