package models

import "github.com/rs/xid"

// the reason I make a separate dto for the Product model is to make it easier
// for the client side to add or remove categories of a product record
//
// in this model the Categories field is set to be a slice of category ids
// instead of a slice of category structs and later, those ids will be used
// as references to get the actual categories from the database...
type ProductDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Price       uint32 `json:"price" binding:"required"`
	Discount    uint8  `json:"discount"`
	Quantity    uint32 `json:"quantity" binding:"required"`

	Categories []xid.ID `json:"categories"`
}
