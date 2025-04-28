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
	ImageCover    string    `gorm:"type:text" json:"image_cover"`
	ImageURL      string    `gorm:"type:text" json:"image_url"`
	CloudImageID  uint      `gorm:"not null" json:"cloud_image_id"`
	ProductFormID uint      `gorm:"not null" json:"product_form_id"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Seller       User         `gorm:"foreignKey:SellerID" json:"seller"`
	ProductType  ProductType  `gorm:"foreignKey:ProductTypeID" json:"product_type"`
	CloudService CloudService `gorm:"foreignKey:CloudImageID" json:"cloud_service"`
	ProductForm  ProductForm  `gorm:"foreignKey:ProductFormID" json:"product_form"`
}
