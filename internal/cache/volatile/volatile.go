package volatile

func StartVolatileCacheRefresher() {
	go refreshCategoryDiscounts()
	go refreshProductTypeDiscounts()
}
