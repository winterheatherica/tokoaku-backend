package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
	"gorm.io/gorm"
)

func SeedProductTypes(db *gorm.DB) {
	productTypes := []models.ProductType{
		{
			ID:           1,
			Name:         "Automotive",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629704/Automotive_fxc8bm.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        12,
			CreatedAt:    time.Now(),
		},
		{
			ID:           2,
			Name:         "Baby",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629704/Baby_zf2kzy.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        8,
			CreatedAt:    time.Now(),
		},
		{
			ID:           3,
			Name:         "Books",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629703/Books_oxptux.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        4,
			CreatedAt:    time.Now(),
		},
		{
			ID:           4,
			Name:         "CDs and Vinyl",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629703/CDs_and_Vinyl_gharcl.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        6,
			CreatedAt:    time.Now(),
		},
		{
			ID:           5,
			Name:         "Camera and Photo",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629702/Camera_and_Photo_bsc3pb.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        10,
			CreatedAt:    time.Now(),
		},
		{
			ID:           6,
			Name:         "Cellphones and Accessories",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629701/Cellphones_and_Accessories_l6vub6.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        15,
			CreatedAt:    time.Now(),
		},
		{
			ID:           7,
			Name:         "Clothing",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629701/Clothing_xt7kki.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        10,
			CreatedAt:    time.Now(),
		},
		{
			ID:           8,
			Name:         "Computers and Accessories",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629701/Computers_and_Accessories_vof4xi.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        8,
			CreatedAt:    time.Now(),
		},
		{
			ID:           9,
			Name:         "Grocery and Gourmet Food",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629700/Grocery_and_Gourmet_Food_oghbki.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        4,
			CreatedAt:    time.Now(),
		},
		{
			ID:           10,
			Name:         "Health and Beauty",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629700/Health_and_Beauty_sb18ez.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        5,
			CreatedAt:    time.Now(),
		},
		{
			ID:           11,
			Name:         "Home and Garden",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629699/Home_and_Garden_etmcy7.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        7,
			CreatedAt:    time.Now(),
		},
		{
			ID:           12,
			Name:         "Jewelry",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629698/Jewelry_axpsyd.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        12,
			CreatedAt:    time.Now(),
		},
		{
			ID:           13,
			Name:         "Luggage and Travel Gear",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629698/Luggage_and_Travel_Gear_kgkkhz.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        8,
			CreatedAt:    time.Now(),
		},
		{
			ID:           14,
			Name:         "Movies and TV",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629698/Movies_and_TV_uzawjr.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        6,
			CreatedAt:    time.Now(),
		},
		{
			ID:           15,
			Name:         "Musical Instruments",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629697/Musical_Instruments_ni10sv.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        10,
			CreatedAt:    time.Now(),
		},
		{
			ID:           16,
			Name:         "Office Products",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629698/Office_Products_mci21s.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        5,
			CreatedAt:    time.Now(),
		},
		{
			ID:           17,
			Name:         "Other Electronics",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629698/Other_Electronics_yyqmmg.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        8,
			CreatedAt:    time.Now(),
		},
		{
			ID:           18,
			Name:         "Others",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629697/Others_hrs4pc.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        4,
			CreatedAt:    time.Now(),
		},
		{
			ID:           19,
			Name:         "Pet Supplies",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629697/Pet_Supplies_zrrniq.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        6,
			CreatedAt:    time.Now(),
		},
		{
			ID:           20,
			Name:         "Shoes",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629697/Shoes_zfyz6r.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        8,
			CreatedAt:    time.Now(),
		},
		{
			ID:           21,
			Name:         "Sports and Outdoors",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629696/Sports_and_Outdoors_z5zlvv.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        6,
			CreatedAt:    time.Now(),
		},
		{
			ID:           22,
			Name:         "Tools and Home Improvement",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629696/Tools_and_Home_Improvement_kyl8kh.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        7,
			CreatedAt:    time.Now(),
		},
		{
			ID:           23,
			Name:         "Toys and Games",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629551/Toys_and_Games_yvthg3.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        5,
			CreatedAt:    time.Now(),
		},
		{
			ID:           24,
			Name:         "Video Games",
			ImageURL:     "https://res.cloudinary.com/dokzc5ogk/image/upload/v1748629512/Video_Games_pkurrl.jpg",
			CloudImageID: 2,
			ValueTypeID:  1,
			Value:        10,
			CreatedAt:    time.Now(),
		},
	}

	log.Printf("[DEBUG] Menemukan %d product types untuk seeding", len(productTypes))

	for _, pt := range productTypes {
		pt.Slug = utils.SlugifyText(pt.Name)
		if err := db.FirstOrCreate(&pt, models.ProductType{ID: pt.ID}).Error; err != nil {
			log.Printf("❌ Gagal seeding product type ID %d: %v\n", pt.ID, err)
		}
	}

	log.Println("[SEEDER] ⚙️  product types seeded")
}
