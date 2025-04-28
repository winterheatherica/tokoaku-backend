package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"gorm.io/gorm"
)

func SeedProductTypes(db *gorm.DB) {

	productTypes := []models.ProductType{
		{ID: 1, Name: "Automotive", Slug: "automotive", ValueTypeID: 1, Value: 12, CreatedAt: time.Now()},
		{ID: 2, Name: "Baby", Slug: "baby", ValueTypeID: 1, Value: 8, CreatedAt: time.Now()},
		{ID: 3, Name: "Books", Slug: "books", ValueTypeID: 1, Value: 4, CreatedAt: time.Now()},
		{ID: 4, Name: "CDs and Vinyl", Slug: "cds-and-vinyl", ValueTypeID: 1, Value: 6, CreatedAt: time.Now()},
		{ID: 5, Name: "Camera and Photo", Slug: "camera-and-photo", ValueTypeID: 1, Value: 10, CreatedAt: time.Now()},
		{ID: 6, Name: "Cellphones and Accessories", Slug: "cellphones-and-accessories", ValueTypeID: 1, Value: 15, CreatedAt: time.Now()},
		{ID: 7, Name: "Clothing", Slug: "clothing", ValueTypeID: 1, Value: 10, CreatedAt: time.Now()},
		{ID: 8, Name: "Computers and Accessories", Slug: "computers-and-accessories", ValueTypeID: 1, Value: 8, CreatedAt: time.Now()},
		{ID: 9, Name: "Grocery and Gourmet Food", Slug: "grocery-and-gourmet-food", ValueTypeID: 1, Value: 4, CreatedAt: time.Now()},
		{ID: 10, Name: "Health and Beauty", Slug: "health-and-beauty", ValueTypeID: 1, Value: 5, CreatedAt: time.Now()},
		{ID: 11, Name: "Home and Garden", Slug: "home-and-garden", ValueTypeID: 1, Value: 7, CreatedAt: time.Now()},
		{ID: 12, Name: "Jewelry", Slug: "jewelry", ValueTypeID: 1, Value: 12, CreatedAt: time.Now()},
		{ID: 13, Name: "Luggage and Travel Gear", Slug: "luggage-and-travel-gear", ValueTypeID: 1, Value: 8, CreatedAt: time.Now()},
		{ID: 14, Name: "Movies and TV", Slug: "movies-and-tv", ValueTypeID: 1, Value: 6, CreatedAt: time.Now()},
		{ID: 15, Name: "Musical Instruments", Slug: "musical-instruments", ValueTypeID: 1, Value: 10, CreatedAt: time.Now()},
		{ID: 16, Name: "Office Products", Slug: "office-products", ValueTypeID: 1, Value: 5, CreatedAt: time.Now()},
		{ID: 17, Name: "Other Electronics", Slug: "other-electronics", ValueTypeID: 1, Value: 8, CreatedAt: time.Now()},
		{ID: 18, Name: "Others", Slug: "others", ValueTypeID: 1, Value: 4, CreatedAt: time.Now()},
		{ID: 19, Name: "Pet Supplies", Slug: "pet-supplies", ValueTypeID: 1, Value: 6, CreatedAt: time.Now()},
		{ID: 20, Name: "Shoes", Slug: "shoes", ValueTypeID: 1, Value: 8, CreatedAt: time.Now()},
		{ID: 21, Name: "Sports and Outdoors", Slug: "sports-and-outdoors", ValueTypeID: 1, Value: 6, CreatedAt: time.Now()},
		{ID: 22, Name: "Tools and Home Improvement", Slug: "tools-and-home-improvement", ValueTypeID: 1, Value: 7, CreatedAt: time.Now()},
		{ID: 23, Name: "Toys and Games", Slug: "toys-and-games", ValueTypeID: 1, Value: 5, CreatedAt: time.Now()},
		{ID: 24, Name: "Video Games", Slug: "video-games", ValueTypeID: 1, Value: 10, CreatedAt: time.Now()},
	}

	for _, productType := range productTypes {
		if err := db.FirstOrCreate(&productType, models.ProductType{ID: productType.ID}).Error; err != nil {
			log.Printf("Gagal seeding ProductType ID %d: %v\n", productType.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  product types seeded")
}
