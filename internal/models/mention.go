package models

import (
	"github.com/google/uuid"
)

type Mention struct {
	PostId                uuid.UUID `json:"postId" gorm:"type:varchar2(36);index:user_and_accept_show_in_profile_and_post,priority:3;foreignKey:PostId"`
	UserId                string    `json:"userId" gorm:"type:varchar2(24);index:user_and_accept_show_in_profile_and_post,priority:1"`
	AcceptedShowInProfile bool      `json:"acceptedShowInProfile" gorm:"index:user_and_accept_show_in_profile_and_post,priority:2"`
}
