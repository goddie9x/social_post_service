package models

type Mention struct {
	PostId                string `json:"postId" gorm:"type:varchar2(36);index:user_and_accept_show_in_profile_and_post,priority:3"`
	UserId                string `json:"userId" gorm:"type:varchar2(24);index:user_and_accept_show_in_profile_and_post,priority:1"`
	AcceptedShowInProfile bool   `json:"acceptedShowInProfile" gorm:"index:user_and_accept_show_in_profile_and_post,priority:2"`
}
