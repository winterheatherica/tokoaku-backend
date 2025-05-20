package models

import (
	"time"
)

type Product struct {
	ID            string    `gorm:"primaryKey;type:varchar(100)" json:"id"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	Description   string    `gorm:"type:text" json:"description"`
	SellerID      string    `gorm:"not null" json:"seller_id"`
	ProductTypeID uint      `gorm:"not null" json:"product_type_id"`
	ImageCoverURL string    `gorm:"type:text" json:"image_cover_url"`
	CloudImageID  uint      `gorm:"not null" json:"cloud_image_id"`
	ProductFormID uint      `gorm:"not null" json:"product_form_id"`
	Slug          string    `gorm:"uniqueIndex;not null" json:"slug"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Seller       User         `gorm:"foreignKey:SellerID" json:"seller"`
	ProductType  ProductType  `gorm:"foreignKey:ProductTypeID" json:"product_type"`
	CloudService CloudService `gorm:"foreignKey:CloudImageID" json:"cloud_service"`
	ProductForm  ProductForm  `gorm:"foreignKey:ProductFormID" json:"product_form"`

	Variants          []ProductVariant  `gorm:"foreignKey:ProductID" json:"variants"`
	ProductCategories []ProductCategory `gorm:"foreignKey:ProductID"`
}
