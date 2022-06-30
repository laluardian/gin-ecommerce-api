package models

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

// note that this model is under the presumption that all products are belong
// to a single entity which is the ecommerce itself and any user assigned
// the role "admin" can create, update, or delete those products
type Product struct {
	ID          xid.ID    `gorm:"primarykey;not null;unique" json:"id"`
	Name        string    `gorm:"not null;index" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	Price       uint32    `gorm:"not null" json:"price"`
	Discount    uint8     `json:"discount"`
	Quantity    uint32    `gorm:"not null" json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	WishlistedBy []*User `gorm:"many2many:user_wishlist_products" json:"wishlisted_by,omitempty"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.ID = xid.New()
	return nil
}
