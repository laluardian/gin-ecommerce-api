package models

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

// an address belongs to a user, a user can have many address
type Address struct {
	ID                  xid.ID    `gorm:"<-:create;primarykey;not null" json:"id"`
	AddressName         string    `gorm:"not null;size:32" json:"address_name"`
	ReceiverName        string    `gorm:"not null;size:32" json:"receiver_name"`
	ReceiverPhoneNumber string    `gorm:"not null" json:"receiver_phone_number"`
	StreetAddress       string    `gorm:"not null;size:64" json:"street_address"`
	City                string    `gorm:"not null" json:"city"`
	Province            string    `gorm:"not null" json:"province"`
	Country             string    `gorm:"not null" json:"country"`
	ZipCode             string    `gorm:"not null" json:"zip_code"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	UserID xid.ID `gorm:"not null" json:"-"`
}

func (a *Address) BeforeCreate(tx *gorm.DB) error {
	a.ID = xid.New()
	return nil
}
