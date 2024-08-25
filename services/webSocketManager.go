package services

import (
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"price-alert-system/models"
)

type WebSocketManager struct {
	binanceService   *BinanceService
	indicatorService *IndicatorService
	alertService     *AlertService
	conn             *websocket.Conn
	done             chan struct{}
}

func NewWebSocketManager(binanceService *BinanceService, indicatorService *IndicatorService, alertService *AlertService) *WebSocketManager {
	return &WebSocketManager{
		binanceService:   binanceService,
		indicatorService: indicatorService,
		alertService:     alertService,
		done:             make(chan struct{}),
	}
}

func (m *WebSocketManager) Start() {
	for {
		if err := m.connect(); err != nil {
			log.Printf("WebSocket connection error: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		m.readMessages()

		select {
		case <-m.done:
			return
		default:
			log.Println("WebSocket connection closed. Reconnecting...")
			time.Sleep(5 * time.Second)
		}
	}
}

func (m *WebSocketManager) Stop() {
	close(m.done)
	if m.conn != nil {
		m.conn.Close()
	}
}

func (m *WebSocketManager) connect() error {
	var err error
	m.conn, err = m.binanceService.ConnectWebSocket()
	return err
}

func (m *WebSocketManager) readMessages() {
	for {
		trade, err := m.binanceService.ReadMessage(m.conn)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		price, err := strconv.ParseFloat(trade.Price, 64)
		if err != nil {
			log.Printf("Error parsing trade price: %v", err)
			continue
		}

		kline := models.Kline{
			OpenTime:  trade.TradeTime,
			Open:      price,
			High:      price,
			Low:       price,
			Close:     price,
			Volume:    0, // We don't have volume information in a single trade
			CloseTime: trade.TradeTime,
		}

		m.indicatorService.UpdateKlines(kline)

		select {
		case <-m.done:
			return
		default:
		}
	}
}
