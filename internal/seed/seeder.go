package seed

import (
	"log"

	"gorm.io/gorm"
)

func RunAllSeeders(db *gorm.DB) {
	SeedValueTypes(db)
	SeedEventTypes(db)
	SeedBankLists(db)
	SeedDefaultFees(db)
	SeedBankTransferFees(db)
	SeedCategories(db)
	SeedNotificationTypes(db)
	SeedPaymentMethods(db)
	SeedProductForms(db)
	SeedProductTypes(db)
	SeedPromos(db)
	SeedProviders(db)
	SeedCloudServices(db)
	SeedRoles(db)
	SeedDiscountSponsors(db)
	SeedDiscounts(db)
	SeedSentiments(db)
	SeedShippingOptions(db)
	SeedStatuses(db)
	SeedCurrentEvent(db)

	log.Println("[SEEDER] ✅ All seeders executed")
}
