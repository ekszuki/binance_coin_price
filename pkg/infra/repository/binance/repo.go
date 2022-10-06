package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"br.com.sygnux/binance/pkg/domain"
	"br.com.sygnux/binance/pkg/domain/binance"
)

type CoinRepository struct {
	ctx      context.Context
	http     http.Client
	priceURL string
}

// GetCoinValue implements binance.CoinRepository
func (r *CoinRepository) GetCoinValue(ctx context.Context, coinID string) (*domain.BinanceCoinPrice, error) {
	url := fmt.Sprintf(r.priceURL, coinID)

	resp, err := r.http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	coinPrice := new(domain.BinanceCoinPrice)
	err = json.Unmarshal(bs, coinPrice)
	if err != nil {
		return nil, err
	}

	return coinPrice, nil
}

func NewCoinRepository(
	ctx context.Context,
	http *http.Client,
	priceURL string,
) binance.CoinRepository {
	return &CoinRepository{
		ctx:      ctx,
		http:     *http,
		priceURL: priceURL,
	}
}
