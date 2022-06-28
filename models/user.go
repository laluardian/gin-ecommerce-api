package models

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type User struct {
	ID        xid.ID    `gorm:"<-:create;primarykey;not null" json:"id"`
	Username  string    `gorm:"not null;unique;size:24" json:"username"`
	Email     string    `gorm:"not null;unique;" json:"email"`
	Password  string    `gorm:"not null" json:"password"`
	IsAdmin   bool      `gorm:"not null" json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Addresses []Address `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"addresses,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = xid.New()
	return nil
}
