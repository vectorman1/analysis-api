package common

import "net/http"

const TRADING212_INSTRUMENTS_LINK = `https://www.trading212.com/en/Trade-Equities`
const TRADING212_SHOW_ALL_BUTTON_SELECTOR = `div.conditions-table > div > div.view-more > a`
const TRADING212_ALL_INSTRUMENTS_SELECTOR = `#all-equities`
const SYMBOLS_NAMCESPACE = `53edcce7-94d4-4deb-b2ac-d1f6d8657d8e`
const MaxConcurrency = 2000

func GetAllowedHeaders() []string {
	return []string{"Authorization", "Accept", "Origin", "DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Content-Range", "Range"}
}
func GetAllowedMethods() []string {
	return []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}
}
