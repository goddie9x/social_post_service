package models

import (
	"errors"
	"post_service/internal/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Mention struct {
	Post   Post   `json:"postId" gorm:"foreignKey:postId,index"`
	UserId string `json:"userId" gorm:"index"`
}
type Tag struct {
	PostId Post   `json:"postId" gorm:"foreignKey:postId,index"`
	Name   string `json:"name" gorm:"index"`
}
type Post struct {
	Id        uuid.UUID             `json:id gorm:"type:uuid,default:uuid_generate_v4(),primary_key"`
	OwnerId   string                `json:ownerId binding:"required" gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:2;index:idx_owner"`
	Type      constants.PostType    `json:type binding:"required" gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:3";index:idx_target_type_created_at`
	TargetId  string                `json:targetId binding:"required" gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:1";index:idx_target_type_created_at`
	Content   string                `json:content`
	CreatedAt time.Time             `json:createAt gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:6,autoCreateTime"`
	UpdateAt  time.Time             `json:updateAt gorm:"autoUpdateTime"`
	Privacy   constants.PrivacyType `json:privacy gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:5"`
	Approved  bool                  `json:approved gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:4"`
	Mentions  []Mention             `json:"mentions"`
	Tags      []Tag                 `json:"tags" gorm:"many2many:tag"`
}

func (p *Post) Validate() (err error) {
	if p.Type == constants.GroupPost && p.Privacy != constants.Public {
		return errors.New("Group just allow public post")
	}
	return
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if err := p.Validate(); err != nil {
		return err
	}
	return
}

func (p *Post) BeforeUpdate(tx *gorm.DB) (err error) {
	if err := p.Validate(); err != nil {
		return err
	}
	return
}
func (p *Post) BeforeDelete(tx *gorm.DB) (err error) {
	if err := tx.Where("postId = ?", p.Id).Delete(Mention{}).Error; err != nil {
		return err
	}
	return
}
func (p *Post) AcceptNewData(post *Post) {
	if post.OwnerId != "" {
		p.OwnerId = post.OwnerId
	}
	p.Type = post.Type
	if post.TargetId != "" {
		p.TargetId = post.TargetId
	}
	if post.Content != "" {
		p.Content = post.Content
	}
	if !post.CreatedAt.IsZero() {
		p.CreatedAt = post.CreatedAt
	}
	if !post.UpdateAt.IsZero() {
		p.UpdateAt = post.UpdateAt
	}
	p.Privacy = post.Privacy
	p.Approved = post.Approved
	if len(post.Mentions) > 0 {
		p.Mentions = post.Mentions
	}
	if len(post.Tags) > 0 {
		p.Tags = post.Tags
	}
}
