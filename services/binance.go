package services

import (
	"encoding/json"
	"net/url"

	"github.com/gorilla/websocket"

	"price-alert-system/models"
)

type BinanceService struct {
	baseURL string
	wsURL   string
}

func NewBinanceService() *BinanceService {
	return &BinanceService{
		baseURL: "https://api.binance.com",
		wsURL:   "wss://stream.binance.com:443/ws/btcusdt@trade",
	}
}

// func (s *BinanceService) GetUIKlines(symbol string, interval string, limit int, startTime, endTime int64, timeZone string) ([]models.Kline, error) {
// 	url := fmt.Sprintf("%s/api/v3/uiKlines?symbol=%s&interval=%s&limit=%d", s.baseURL, symbol, interval, limit)

// 	if startTime != 0 {
// 		url += fmt.Sprintf("&startTime=%d", startTime)
// 	}
// 	if endTime != 0 {
// 		url += fmt.Sprintf("&endTime=%d", endTime)
// 	}
// 	if timeZone != "" {
// 		url += fmt.Sprintf("&timeZone=%s", timeZone)
// 	}

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var rawKlines [][]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&rawKlines); err != nil {
// 		return nil, err
// 	}

// 	klines := make([]models.Kline, len(rawKlines))
// 	for i, raw := range rawKlines {
// 		openTime, _ := strconv.ParseInt(raw[0].(string), 10, 64)
// 		open, _ := strconv.ParseFloat(raw[1].(string), 64)
// 		high, _ := strconv.ParseFloat(raw[2].(string), 64)
// 		low, _ := strconv.ParseFloat(raw[3].(string), 64)
// 		close, _ := strconv.ParseFloat(raw[4].(string), 64)
// 		volume, _ := strconv.ParseFloat(raw[5].(string), 64)
// 		closeTime, _ := strconv.ParseInt(raw[6].(string), 10, 64)
// 		quoteAssetVolume, _ := strconv.ParseFloat(raw[7].(string), 64)
// 		numberOfTrades, _ := strconv.ParseInt(raw[8].(string), 10, 64)
// 		takerBuyBaseAssetVolume, _ := strconv.ParseFloat(raw[9].(string), 64)
// 		takerBuyQuoteAssetVolume, _ := strconv.ParseFloat(raw[10].(string), 64)

// 		klines[i] = models.Kline{
// 			OpenTime:                 openTime,
// 			Open:                     open,
// 			High:                     high,
// 			Low:                      low,
// 			Close:                    close,
// 			Volume:                   volume,
// 			CloseTime:                closeTime,
// 			QuoteAssetVolume:         quoteAssetVolume,
// 			NumberOfTrades:           numberOfTrades,
// 			TakerBuyBaseAssetVolume:  takerBuyBaseAssetVolume,
// 			TakerBuyQuoteAssetVolume: takerBuyQuoteAssetVolume,
// 		}
// 	}

// 	return klines, nil
// }

func (s *BinanceService) ConnectWebSocket() (*websocket.Conn, error) {
	u, _ := url.Parse(s.wsURL)
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *BinanceService) ReadMessage(conn *websocket.Conn) (models.Trade, error) {
	_, message, err := conn.ReadMessage()
	if err != nil {
		return models.Trade{}, err
	}

	var trade models.Trade
	if err := json.Unmarshal(message, &trade); err != nil {
		return models.Trade{}, err
	}

	// print indent json data
	// jsonData, _ := json.MarshalIndent(trade.Price, "", "  ")
	// fmt.Println(string(jsonData))

	return trade, nil
}
