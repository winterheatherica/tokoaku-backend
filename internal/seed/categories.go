package seed

import (
	"log"
	"time"

	"github.com/winterheatherica/tokoaku-backend/internal/models"
	"github.com/winterheatherica/tokoaku-backend/internal/utils"
	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB) {
	categories := []models.Category{
		// For
		{Name: "Men", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Women", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Kids", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Unisex", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Teen Boys", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Teen Girls", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Infants", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Toddlers", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Baby Boys", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Baby Girls", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Elderly Men", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Elderly Women", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Men Formal", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Women Casual", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Boys Sportswear", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Girls Party", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Unisex Travel", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Men Outdoor", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},
		{Name: "Women Work", Code: nil, CategoryLabelID: 1, CreatedAt: time.Now()},

		// Opinion
		{Name: "Elegant", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Sporty", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Classic", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Trendy", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Luxury", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Minimalist", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Bohemian", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Streetwear", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Athleisure", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Eclectic", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Preppy", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Chic", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Rugged", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Urban", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},
		{Name: "Avant-garde", Code: nil, CategoryLabelID: 2, CreatedAt: time.Now()},

		// Size
		{Name: "Extra Small", Code: ptr("XS"), CategoryLabelID: 3, CreatedAt: time.Now()},
		{Name: "Small", Code: ptr("S"), CategoryLabelID: 3, CreatedAt: time.Now()},
		{Name: "Medium", Code: ptr("M"), CategoryLabelID: 3, CreatedAt: time.Now()},
		{Name: "Large", Code: ptr("L"), CategoryLabelID: 3, CreatedAt: time.Now()},
		{Name: "Extra Large", Code: ptr("XL"), CategoryLabelID: 3, CreatedAt: time.Now()},
		{Name: "Double Extra Large", Code: ptr("XXL"), CategoryLabelID: 3, CreatedAt: time.Now()},
		{Name: "Triple Extra Large", Code: ptr("XXXL"), CategoryLabelID: 3, CreatedAt: time.Now()},
		{Name: "One Size", Code: ptr("OS"), CategoryLabelID: 3, CreatedAt: time.Now()},

		// Product Age
		{Name: "New", Code: nil, CategoryLabelID: 4, CreatedAt: time.Now()},
		{Name: "Vintage", Code: nil, CategoryLabelID: 4, CreatedAt: time.Now()},
		{Name: "Used", Code: nil, CategoryLabelID: 4, CreatedAt: time.Now()},
		{Name: "Antique", Code: nil, CategoryLabelID: 4, CreatedAt: time.Now()},
		{Name: "Brand New", Code: nil, CategoryLabelID: 4, CreatedAt: time.Now()},
		{Name: "Like New", Code: nil, CategoryLabelID: 4, CreatedAt: time.Now()},
		{Name: "Refurbished", Code: nil, CategoryLabelID: 4, CreatedAt: time.Now()},
		{Name: "Second-hand", Code: nil, CategoryLabelID: 4, CreatedAt: time.Now()},
		{Name: "Pre-Owned", Code: nil, CategoryLabelID: 4, CreatedAt: time.Now()},
		{Name: "Retro", Code: nil, CategoryLabelID: 4, CreatedAt: time.Now()},

		// Shape
		{Name: "Round", Code: nil, CategoryLabelID: 5, CreatedAt: time.Now()},
		{Name: "Square", Code: nil, CategoryLabelID: 5, CreatedAt: time.Now()},
		{Name: "Oval", Code: nil, CategoryLabelID: 5, CreatedAt: time.Now()},
		{Name: "Rectangular", Code: nil, CategoryLabelID: 5, CreatedAt: time.Now()},
		{Name: "Heart", Code: nil, CategoryLabelID: 5, CreatedAt: time.Now()},
		{Name: "Triangle", Code: nil, CategoryLabelID: 5, CreatedAt: time.Now()},
		{Name: "Hexagonal", Code: nil, CategoryLabelID: 5, CreatedAt: time.Now()},
		{Name: "Pentagon", Code: nil, CategoryLabelID: 5, CreatedAt: time.Now()},
		{Name: "Star", Code: nil, CategoryLabelID: 5, CreatedAt: time.Now()},
		{Name: "Diamond", Code: nil, CategoryLabelID: 5, CreatedAt: time.Now()},

		// Color
		{Name: "White", Code: ptr("#FFFFFF"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Black", Code: ptr("#000000"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Red", Code: ptr("#FF0000"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Blue", Code: ptr("#0000FF"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Yellow", Code: ptr("#FFFF00"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Green", Code: ptr("#00FF00"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Pink", Code: ptr("#FFC0CB"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Purple", Code: ptr("#800080"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Orange", Code: ptr("#FFA500"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Brown", Code: ptr("#A52A2A"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Gray", Code: ptr("#808080"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Cyan", Code: ptr("#00FFFF"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Magenta", Code: ptr("#FF00FF"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Maroon", Code: ptr("#800000"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Olive", Code: ptr("#808000"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Teal", Code: ptr("#008080"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Navy", Code: ptr("#000080"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Lime", Code: ptr("#00FF00"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Coral", Code: ptr("#FF7F50"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Turquoise", Code: ptr("#40E0D0"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Lavender", Code: ptr("#E6E6FA"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Beige", Code: ptr("#F5F5DC"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Mint", Code: ptr("#98FF98"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Peach", Code: ptr("#FFE5B4"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Ivory", Code: ptr("#FFFFF0"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Gold", Code: ptr("#FFD700"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Silver", Code: ptr("#C0C0C0"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Bronze", Code: ptr("#CD7F32"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Khaki", Code: ptr("#F0E68C"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Charcoal", Code: ptr("#36454F"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Salmon", Code: ptr("#FA8072"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Crimson", Code: ptr("#DC143C"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Indigo", Code: ptr("#4B0082"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Violet", Code: ptr("#8A2BE2"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Periwinkle", Code: ptr("#CCCCFF"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Aqua", Code: ptr("#00FFFF"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Chartreuse", Code: ptr("#7FFF00"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Fuchsia", Code: ptr("#FF00FF"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Emerald", Code: ptr("#50C878"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Amber", Code: ptr("#FFBF00"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Mustard", Code: ptr("#FFDB58"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Sand", Code: ptr("#C2B280"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Rust", Code: ptr("#B7410E"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Mauve", Code: ptr("#E0B0FF"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Sage", Code: ptr("#BCB88A"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Lilac", Code: ptr("#C8A2C8"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Blush", Code: ptr("#DE5D83"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Cerulean", Code: ptr("#007BA7"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Plum", Code: ptr("#DDA0DD"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Taupe", Code: ptr("#483C32"), CategoryLabelID: 6, CreatedAt: time.Now()},
		{Name: "Sepia", Code: ptr("#704214"), CategoryLabelID: 6, CreatedAt: time.Now()},

		// Origin
		{Name: "Japan", Code: ptr("JP"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "USA", Code: ptr("US"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "France", Code: ptr("FR"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Italy", Code: ptr("IT"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Indonesia", Code: ptr("ID"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Germany", Code: ptr("DE"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "United Kingdom", Code: ptr("GB"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Canada", Code: ptr("CA"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Australia", Code: ptr("AU"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "South Korea", Code: ptr("KR"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "China", Code: ptr("CN"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Singapore", Code: ptr("SG"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Malaysia", Code: ptr("MY"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Thailand", Code: ptr("TH"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Vietnam", Code: ptr("VN"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "India", Code: ptr("IN"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Brazil", Code: ptr("BR"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Mexico", Code: ptr("MX"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Russia", Code: ptr("RU"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "South Africa", Code: ptr("ZA"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "New Zealand", Code: ptr("NZ"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Netherlands", Code: ptr("NL"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Spain", Code: ptr("ES"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Sweden", Code: ptr("SE"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Norway", Code: ptr("NO"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Finland", Code: ptr("FI"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Denmark", Code: ptr("DK"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Belgium", Code: ptr("BE"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Switzerland", Code: ptr("CH"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Austria", Code: ptr("AT"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Poland", Code: ptr("PL"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Portugal", Code: ptr("PT"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Greece", Code: ptr("GR"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Turkey", Code: ptr("TR"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "Saudi Arabia", Code: ptr("SA"), CategoryLabelID: 7, CreatedAt: time.Now()},
		{Name: "United Arab Emirates", Code: ptr("AE"), CategoryLabelID: 7, CreatedAt: time.Now()},

		// Material
		{Name: "Cotton", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Wool", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Silk", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Leather", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Polyester", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Denim", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Linen", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Rayon", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Velvet", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Canvas", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Acrylic", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Spandex", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Cashmere", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Bamboo", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Hemp", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Nylon", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Modal", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Suede", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Jersey", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},
		{Name: "Microfiber", Code: nil, CategoryLabelID: 8, CreatedAt: time.Now()},

		// Purpose
		{Name: "Formal", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Casual", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Party", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Work", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Travel", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Sportswear", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Business", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Wedding", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Beachwear", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Homewear", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Loungewear", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Sleepwear", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Hiking", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Cycling", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Running", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Training", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Formal Event", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Picnic", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Outdoor", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
		{Name: "Festival", Code: nil, CategoryLabelID: 9, CreatedAt: time.Now()},
	}

	for i := range categories {
		category := &categories[i]
		category.Slug = utils.SlugifyText(category.Name)
		if err := db.FirstOrCreate(category, models.Category{
			Name:            category.Name,
			CategoryLabelID: category.CategoryLabelID,
		}).Error; err != nil {
			log.Printf("Gagal seeding category %s: %v\n", category.Name, err)
		}
	}

	log.Println("[SEEDER] ⚙️  categories seeded")
}

func ptr(s string) *string {
	return &s
}
