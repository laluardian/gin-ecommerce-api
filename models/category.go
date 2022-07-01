package models

import (
	"time"

	"github.com/gosimple/slug"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type Category struct {
	ID          xid.ID    `gorm:"<-:create;primarykey;not null;unique" json:"id"`
	Name        string    `gorm:"not null;unique" json:"name"`
	Description string    `json:"description,omitempty"`
	Slug        string    `gorm:"not null;unique" json:"slug"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Products []*Product `gorm:"many2many:product_categories" json:"products,omitempty"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	c.ID = xid.New()
	c.Slug = slug.Make(c.Name)
	return nil
}

func (c *Category) BeforeUpdate(tx *gorm.DB) error {
	c.Slug = slug.Make(c.Name)
	return nil
}
