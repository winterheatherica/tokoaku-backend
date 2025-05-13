package volatile

func StartVolatileCacheRefresher() {
	go refreshCategoryDiscounts()
}
