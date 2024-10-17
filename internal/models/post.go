package models

import (
	"errors"
	"post_service/internal/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	Id        string                `json:"id" gorm:"type:varchar2(16);primaryKey"`
	OwnerId   string                `json:"ownerId" gorm:"type:varchar2(24);index:idx_target_owner_type_approved_privacy_created_at,priority:2;index:idx_owner"`
	Type      constants.PostType    `json:"type" binding:"required" gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:3;index:idx_target_type_created_at"`
	TargetId  string                `json:"targetId" gorm:"type:varchar2(24);index:idx_target_owner_type_approved_privacy_created_at,priority:1;index:idx_target_type_created_at"`
	Content   string                `json:"content" gorm:"type:text"`
	CreatedAt time.Time             `json:"createdAt" gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:6;autoCreateTime;sort:desc"`
	UpdatedAt time.Time             `json:"updatedAt" gorm:"autoUpdateTime"`
	Privacy   constants.PrivacyType `json:"privacy" gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:5"`
	Approved  bool                  `json:"approved" gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:4"`
	Mentions  []Mention             `json:"mentions" gorm:"foreignKey:PostId"`
	Tags      []*Tag                `json:"tags" gorm:"many2many:post_tags;"`
	BlobIds   []string              `json:"blobIds" gorm:"-"`
}

func (p *Post) Validate() (err error) {
	if p.Content == "" && len(p.BlobIds) < 1 {
		return errors.New("Post must have either content nor attachment")
	}
	if p.Type != constants.GroupPost {
		return
	}
	if p.Privacy != constants.Public {
		return errors.New("Group just allow public post")
	}
	return
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if p.Id == "" {
		p.Id = uuid.New().String()
	}
	if err := p.Validate(); err != nil {
		return err
	}
	if p.Type == constants.GroupPost {
		//TODO: check group exist, and the default value of Accepted = true if group is public and auto accept post
		for i := range p.Mentions {
			p.Mentions[i].AcceptedShowInProfile = false
		}
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
	//todo delete all related blob
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
	if !post.UpdatedAt.IsZero() {
		p.UpdatedAt = post.UpdatedAt
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
