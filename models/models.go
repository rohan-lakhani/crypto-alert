package models

type AlertRequest struct {
	UserID    int     `json:"user_id"`
	Email     string  `json:"email"`
	Value     float64 `json:"value"`
	Direction string  `json:"direction"`
	Indicator string  `json:"indicator"`
}

type Alert struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Email     string  `json:"email"`
	Value     float64 `json:"value"`
	Direction string  `json:"direction"`
	Indicator string  `json:"indicator"`
	Status    string  `json:"status"`
}

type Kline struct {
	OpenTime                 int64   `json:"openTime"`
	Open                     float64 `json:"open"`
	High                     float64 `json:"high"`
	Low                      float64 `json:"low"`
	Close                    float64 `json:"close"`
	Volume                   float64 `json:"volume"`
	CloseTime                int64   `json:"closeTime"`
	QuoteAssetVolume         float64 `json:"quoteAssetVolume"`
	NumberOfTrades           int64   `json:"numberOfTrades"`
	TakerBuyBaseAssetVolume  float64 `json:"takerBuyBaseAssetVolume"`
	TakerBuyQuoteAssetVolume float64 `json:"takerBuyQuoteAssetVolume"`
}

type Trade struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	TradeID   int64  `json:"a"`
	Price     string `json:"p"`
	Quantity  string `json:"q"`
	FirstID   int64  `json:"f"`
	LastID    int64  `json:"l"`
	TradeTime int64  `json:"T"`
	IsBuyerMM bool   `json:"m"`
	Ignore    bool   `json:"M"`
}
