package volatile

func StartVolatileCacheRefresher() {
	go refreshBankList()
	go refreshBankTransferFees()
	go refreshCategories()
	go refreshCategoryDiscounts()
	go refreshCloudServices()
}
