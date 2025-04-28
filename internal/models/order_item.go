package models

type OrderItem struct {
	OrderID          uint `gorm:"not null" json:"order_id"`
	ProductVariantID uint `gorm:"not null" json:"product_variant_id"`
	Quantity         uint `gorm:"not null" json:"quantity"`

	Order          Order          `gorm:"foreignKey:OrderID;references:ID" json:"order"`
	ProductVariant ProductVariant `gorm:"foreignKey:ProductVariantID;references:ID" json:"product_variant"`
}
