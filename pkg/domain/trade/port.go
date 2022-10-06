package trade

import (
	"context"

	"br.com.sygnux/binance/pkg/domain"
)

type Repository interface {
	Save(ctx context.Context, trade *domain.Trade) (*domain.Trade, error)
	FindByKey(ctx context.Context, key string) (*domain.Trade, error)
	FindTradeNotSelled(ctx context.Context) ([]domain.Trade, error)
}
