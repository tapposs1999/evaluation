package pricemodels

type PriceDataDB struct {
	Pair  string   `json:"pair"`
	Data  []DataDB `json:"data"`
}

type DataDB struct {
	Timestamp int64   `json:"timestamp"`
	Price     float64 `json:"price"`
}

type BitkubTradingViewRequestBody struct {
	Symbol string   `json:"symbol"`
	Resolution     string `json:"resolution"`
	FromTime int64 `json:"fromTime"`
}

type BitkubTradingViewResponse struct {
	S string    `json:"s"`
	T []int64   `json:"t"`
	C []float64 `json:"c"`
}