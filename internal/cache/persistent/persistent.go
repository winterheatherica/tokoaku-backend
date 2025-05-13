package persistent

func StartPersistentCacheRefresher() {
	go refreshBankList()
	go refreshBankTransferFees()
	go refreshCategories()
	go refreshCloudServices()
	go refreshProductForms()
	go refreshProductTypes()
	go refreshValueTypes()
}
