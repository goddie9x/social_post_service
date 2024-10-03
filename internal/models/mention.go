package models

type Mention struct {
	Post                  Post   `json:"postId" gorm:"foreignKey:postId,index:user_and_accept_show_in_profile_and_post,priority:3"`
	UserId                string `json:"userId" gorm:"index:user_and_accept_show_in_profile_and_post,priority:1"`
	AcceptedShowInProfile bool   `json:"acceptedShowInProfile" gorm:"index:user_and_accept_show_in_profile_and_post,priority:2"`
}
