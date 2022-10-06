package logic

import (
	"context"
	"strconv"
	"time"

	"br.com.sygnux/binance/pkg/domain/binance"
	"br.com.sygnux/binance/pkg/domain/trade"
	"github.com/AlekSi/pointer"
	"github.com/sirupsen/logrus"
)

func RealTimeBinanceCheck(
	ctx context.Context,
	tradeRepo trade.Repository,
	binanceCoinRepo binance.CoinRepository,
	timeInterval time.Duration,
) {

	go func() {
		logCtx := logrus.WithFields(logrus.Fields{"component": "logic", "function": "RealTimeBinanceCheck"})

		for {
			logCtx.Infof("waiting %v to start new check", timeInterval)
			time.Sleep(timeInterval)

			now := time.Now()
			logCtx.Infof("starting check at %s", now.Format(time.RFC3339))

			trades, err := tradeRepo.FindTradeNotSelled(ctx)
			if err != nil {
				logCtx.Errorf("could not find trades not selled: %v", err)
				continue
			}

			for _, t := range trades {
				logCtx.Infof("getting price on coin %s on the binance api", t.CoinID)
				binanceCoinPrice, err := binanceCoinRepo.GetCoinValue(ctx, t.CoinID)
				if err != nil {
					logCtx.Errorf("could not get price of coin %s on the binance api: %v", t.CoinID, err)
					continue
				}

				currentPrice, err := strconv.ParseFloat(binanceCoinPrice.Price, 64)
				if err != nil {
					logCtx.Errorf("could not convert current price of coin %s: %v", t.CoinID, err)
					continue
				}

				if t.UnitValue == 0 {
					logCtx.Errorf("could not calculate win/los of coin %s: unit value is zero", t.CoinID)
					continue
				}

				t.WinLos = ((currentPrice * 100) / t.UnitValue) - 100
				t.LastUpdate = pointer.ToTime(now)
				t.LastPrice = pointer.ToFloat64(currentPrice)

				_, err = tradeRepo.Save(ctx, &t)
				if err != nil {
					logCtx.Errorf("could not update coin %s on the redis: %v", t.CoinID, err)
				}

				logCtx.Infof("Coin %s - Last price %s - Win/Los %f porcent", t.CoinID, binanceCoinPrice.Price, t.WinLos)
			}

		}

	}()
}
