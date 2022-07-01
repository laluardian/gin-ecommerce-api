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
	ID          xid.ID    `gorm:"<-:create;primarykey;not null;unique" json:"id"`
	Name        string    `gorm:"not null;index" json:"name"`
	Description string    `gorm:"not null" json:"description"`
	Price       uint32    `gorm:"not null" json:"price"`
	Discount    uint8     `json:"discount"`
	Quantity    uint32    `gorm:"not null" json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	WishlistedBy []*User     `gorm:"many2many:user_wishlist_products" json:"wishlisted_by,omitempty"`
	Categories   []*Category `gorm:"many2many:product_categories" json:"categories,omitempty"`
}

type ProductDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       uint32 `json:"price" binding:"required"`
	Discount    uint8  `json:"discount"`
	Quantity    uint32 `json:"quantity" binding:"required"`

	// the Categories field gonna be populated from client side with  the category ids
	// and later those ids will be used in the handlers to get product records from db
	Categories []xid.ID `json:"categories"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	p.ID = xid.New()
	return nil
}
