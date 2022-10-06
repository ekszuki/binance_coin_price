package binance

import (
	"context"

	"br.com.sygnux/binance/pkg/domain"
)

type CoinRepository interface {
	GetCoinValue(ctx context.Context, coinID string) (*domain.BinanceCoinPrice, error)
}
