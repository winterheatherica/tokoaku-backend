package fetcher

import (
	"context"
	"log"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/services/database"
)

type SellerShippingOptionResponse struct {
	ShippingOptionID   uint   `json:"shipping_option_id"`
	CourierName        string `json:"courier_name"`
	CourierServiceName string `json:"courier_service_name"`
	Fee                uint   `json:"fee"`
	EstimatedTime      string `json:"estimated_time"`
	ServiceType        string `json:"service_type"`
}

func GetSellerShippingOptions(ctx context.Context, sellerID string) ([]SellerShippingOptionResponse, error) {
	var options []models.SellerShippingOption

	err := database.DB.WithContext(ctx).
		Preload("ShippingOption").
		Where("seller_id = ?", sellerID).
		Find(&options).Error
	if err != nil {
		log.Printf("[DB] ❌ Gagal ambil opsi pengiriman seller %s: %v", sellerID, err)
		return nil, err
	}

	var results []SellerShippingOptionResponse
	for _, o := range options {
		results = append(results, SellerShippingOptionResponse{
			ShippingOptionID:   o.ShippingOptionID,
			CourierName:        o.ShippingOption.CourierName,
			CourierServiceName: o.ShippingOption.CourierServiceName,
			Fee:                o.ShippingOption.Fee,
			EstimatedTime:      o.ShippingOption.EstimatedTime,
			ServiceType:        o.ShippingOption.ServiceType,
		})
	}

	log.Printf("[DB] ✅ %d opsi pengiriman ditemukan untuk seller %s", len(results), sellerID)
	return results, nil
}
