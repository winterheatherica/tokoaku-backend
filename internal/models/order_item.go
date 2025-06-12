package models

type OrderItem struct {
	OrderShippingID  uint   `gorm:"not null" json:"order_shipping_id"`
	ProductVariantID string `gorm:"not null" json:"product_variant_id"`
	Quantity         uint   `gorm:"not null" json:"quantity"`

	OrderShipping  *Order          `gorm:"foreignKey:OrderShippingID" json:"order_shipping"`
	ProductVariant *ProductVariant `gorm:"foreignKey:ProductVariantID;references:ID" json:"product_variant"`
}
