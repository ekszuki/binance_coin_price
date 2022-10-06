package domain

import (
	"fmt"
	"time"

	"github.com/AlekSi/pointer"
)

const TradePreKey = "trade:"

type Trade struct {
	ID         *string    `json:"id"`
	BuyDate    time.Time  `json:"buy_date" binding:"required"`
	SellDate   *time.Time `json:"sell_date,omitempty"`
	CoinID     string     `json:"coin_id" binding:"required"`
	Quantity   float64    `json:"quantity" binding:"required"`
	UnitValue  float64    `json:"unit_value" binding:"required"`
	WinLos     float64    `json:"win_los"`
	LastPrice  *float64   `json:"last_price,omitempty"`
	LastUpdate *time.Time `json:"last_update,omitempty"`
}

func (t *Trade) BuildKey() (string, error) {
	if t.ID == nil {
		return "", fmt.Errorf("could not create key with nil id")
	}
	key := fmt.Sprintf("%s%s", TradePreKey, pointer.GetString(t.ID))
	return key, nil
}

type BinanceCoinPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}
