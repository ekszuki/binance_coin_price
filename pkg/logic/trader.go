package logic

import (
	"context"

	"br.com.sygnux/binance/pkg/domain"
	"br.com.sygnux/binance/pkg/domain/trade"
	"github.com/sirupsen/logrus"
)

func RegisterTrade(
	ctx context.Context,
	tradeRepo trade.Repository,
	tradePayload *domain.Trade,
) (*domain.Trade, error) {
	logCtx := logrus.WithFields(logrus.Fields{"component": "logic", "function": "RegisterTrade"})
	trade, err := tradeRepo.Save(ctx, tradePayload)
	if err != nil {
		logCtx.Errorf("cannot save trade on database: %v", err)
	}

	return trade, err
}
