package services

import (
	"price-alert-system/models"
	"sync"
	"time"
)

type IndicatorService struct {
	mutex        sync.RWMutex
	klines       []models.Kline
	rsi          float64
	macd         float64
	lastCalc     int64
	calcInterval int64
}

func NewIndicatorService() *IndicatorService {
	return &IndicatorService{
		calcInterval: 60, // Calculate indicators every 60 seconds
		lastCalc:     time.Now().Unix(),
	}
}

func (s *IndicatorService) UpdateKlines(kline models.Kline) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Update the last kline if it's within the same minute, otherwise append
	if len(s.klines) > 0 && s.klines[len(s.klines)-1].CloseTime/60 == kline.CloseTime/60 {
		s.klines[len(s.klines)-1] = kline
	} else {
		s.klines = append(s.klines, kline)
		if len(s.klines) > 100 {
			s.klines = s.klines[1:]
		}
	}

	// Only calculate indicators if the calc interval has passed
	if kline.CloseTime-s.lastCalc >= s.calcInterval {
		s.calculateIndicators()
		s.lastCalc = kline.CloseTime
	}
}

func (s *IndicatorService) calculateIndicators() {
	if len(s.klines) < 26 {
		return // Ensure there are enough klines to calculate indicators
	}

	closes := make([]float64, len(s.klines))
	for i, k := range s.klines {
		closes[i] = k.Close
	}

	s.rsi = calculateRSI(closes, 14)
	s.macd, _, _ = calculateMACD(closes, 12, 26, 9)
}

func (s *IndicatorService) GetIndicators() (float64, float64) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.rsi, s.macd
}

// Calculate RSI for a set of prices
func calculateRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		return 50 // Default to neutral when not enough data
	}

	var gains, losses float64
	// Calculate the initial gains and losses over the period
	for i := 1; i <= period; i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change // Accumulate the loss (note: losses is positive here)
		}
	}

	// Calculate initial average gain and average loss
	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	// Smooth the gains and losses over subsequent prices
	for i := period + 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			avgGain = (avgGain*(float64(period)-1) + change) / float64(period)
			avgLoss = (avgLoss * (float64(period) - 1)) / float64(period)
		} else {
			avgGain = (avgGain * (float64(period) - 1)) / float64(period)
			avgLoss = (avgLoss*(float64(period)-1) - change) / float64(period)
		}
	}

	if avgLoss == 0 {
		return 100 // If there is no loss, RSI is 100 (indicating extremely overbought conditions)
	}

	// Calculate Relative Strength (RS) and RSI
	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}

// Calculate MACD for a set of prices
func calculateMACD(prices []float64, fastPeriod, slowPeriod, signalPeriod int) (float64, float64, float64) {
	if len(prices) < slowPeriod {
		return 0, 0, 0 // Not enough data
	}

	fastEMA := ema(prices, fastPeriod)
	slowEMA := ema(prices, slowPeriod)

	macdLine := fastEMA - slowEMA

	macdHistory := make([]float64, len(prices)-slowPeriod+1)
	for i := range macdHistory {
		fast := ema(prices[:len(prices)-len(macdHistory)+i+1], fastPeriod)
		slow := ema(prices[:len(prices)-len(macdHistory)+i+1], slowPeriod)
		macdHistory[i] = fast - slow
	}
	signalLine := ema(macdHistory, signalPeriod)

	histogram := macdLine - signalLine

	return macdLine, signalLine, histogram
}

// Calculate the Exponential Moving Average (EMA) for a set of prices
func ema(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}

	multiplier := 2.0 / float64(period+1)
	ema := average(prices[:period]) // Initialize EMA with the average of the first period

	for i := period; i < len(prices); i++ {
		ema = ((prices[i] - ema) * multiplier) + ema
	}

	return ema
}

// Calculate the average of a set of numbers
func average(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}
