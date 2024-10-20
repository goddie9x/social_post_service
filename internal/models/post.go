package models

import (
	"blob_store_service/pkg/exceptions"
	"errors"
	"post_service/internal/constants"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Post struct {
	Id        string                 `json:"id" gorm:"type:varchar2(36);primaryKey"`
	OwnerId   string                 `json:"ownerId" gorm:"type:varchar2(24);index:idx_target_owner_type_approved_privacy_created_at,priority:2;index:idx_owner"`
	Type      *constants.PostType    `json:"type" gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:3;index:idx_target_type_created_at"`
	TargetId  string                 `json:"targetId" gorm:"type:varchar2(24);index:idx_target_owner_type_approved_privacy_created_at,priority:1;index:idx_target_type_created_at"`
	Content   string                 `json:"content" gorm:"type:text"`
	CreatedAt time.Time              `json:"createdAt" gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:6;autoCreateTime;sort:desc"`
	UpdatedAt time.Time              `json:"updatedAt" gorm:"autoUpdateTime"`
	Privacy   *constants.PrivacyType `json:"privacy" gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:5"`
	Approved  bool                   `json:"approved" gorm:"index:idx_target_owner_type_approved_privacy_created_at,priority:4"`
	Mentions  []*Mention             `json:"mentions" gorm:"foreignKey:PostId"`
	Tags      []*Tag                 `json:"tags" gorm:"many2many:post_tags;"`
	BlobIds   []string               `json:"blobIds" gorm:"-"`
}

func (p *Post) Validate() error {
	//TODO check targetId exist by query to corresponding service
	if p.TargetId == "" {
		return errors.New("Post must have the targetId")
	}
	if err := p.ValidateForUpdate(); err != nil {
		return err
	}
	return nil
}
func (p *Post) ValidateForUpdate() error {
	if p.Content == "" && len(p.BlobIds) < 1 {
		return errors.New("Post must have either content or attachment")
	}
	if p.Privacy == nil {
		return errors.New("You must have to provide privacy of the post")
	}
	//TODO: validate that user exist
	if *p.Privacy != constants.Public {
		return errors.New("Group just allow public post")
	}
	return nil
}

func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.Id == "" {
		p.Id = uuid.New().String()
	}
	if err := p.Validate(); err != nil {
		return err
	}
	if *p.Type == constants.GroupPost {
		//TODO: check group exist, and the default value of Accepted = true if group is public and auto accept post
		for i := range p.Mentions {
			p.Mentions[i].AcceptedShowInProfile = false
		}
	}
	return nil
}

func (p *Post) BeforeUpdate(tx *gorm.DB) error {
	if err := p.ValidateForUpdate(); err != nil {
		return err
	}
	return nil
}
func (p *Post) BeforeDelete(tx *gorm.DB) error {
	//todo delete all related blob
	if err := tx.Where(`"post_id" = ?`, p.Id).Delete(&Mention{}).Error; err != nil {
		return err
	}
	return nil
}
func (p *Post) AcceptNewDataForUpdate(post *Post) exceptions.CommonExceptionInterface {
	if post.Content != "" {
		p.Content = post.Content
	}
	if post.Privacy != nil {
		p.Privacy = post.Privacy
	}
	p.Approved = post.Approved
	if len(post.Mentions) > 0 {
		p.Mentions = post.Mentions
	}
	if len(post.Tags) > 0 {
		p.Tags = post.Tags
	}
	return nil
}
