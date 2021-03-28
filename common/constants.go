package common

import "net/http"

// urls and selectors for external fetch of symbols
const Trading212InstrumentsLink = `https://www.trading212.com/en/Trade-Equities`
const Trading212ShowAllButtonSelector = `div.conditions-table > div > div.view-more > a`
const Trading212AllInstrumentsSelector = `#all-equities`

// namespace of symbol uuids in order to reproduce them later
const SymbolsNamespace = `53edcce7-94d4-4deb-b2ac-d1f6d8657d8e`

// mongodb related constants
const MongoDbDatabase = `analysis`
const OverviewsCollection = `overviews`
const HistoriesCollection = `histories`

// market names, in order to extract only the relevant ones
const MarketNYSE = `NYSE`
const MarketNASDAQ = `NASDAQ`
const MarketNonISANYSE = `NON-ISA NYSE`
const MarketNonISAOTCMarkets = `NON-ISA OTC Markets`

// constants for random string generation
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// "constant" slice of allowed headers and methdfor CORS config
func GetAllowedHeaders() []string {
	return []string{"Authorization", "Accept", "Origin", "DNT", "X-CustomHeader", "Keep-Alive", "User-Agent", "X-Requested-With", "If-Modified-Since", "Cache-Control", "Content-Type", "Content-Range", "Range"}
}
func GetAllowedMethods() []string {
	return []string{http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}
}
