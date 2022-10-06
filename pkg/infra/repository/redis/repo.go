package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"

	"br.com.sygnux/binance/pkg/domain"
	"br.com.sygnux/binance/pkg/domain/trade"
	"github.com/AlekSi/pointer"
	"github.com/go-redis/redis/v9"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type RedisCfg struct {
	Password string
	Address  string
	Port     int
	DB       int
	UseTLS   bool
}

type TradeRepo struct {
	redis *redis.Client
}

// FindByKey implements trade.Repository
func (r *TradeRepo) FindByKey(ctx context.Context, key string) (*domain.Trade, error) {
	sTrade, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	trade := new(domain.Trade)
	err = json.Unmarshal([]byte(sTrade), trade)
	if err != nil {
		return nil, err
	}

	return trade, nil
}

// FindTradeNotSelled implements trade.Repository
func (r *TradeRepo) FindTradeNotSelled(ctx context.Context) ([]domain.Trade, error) {
	trades := make([]domain.Trade, 0)
	filter := fmt.Sprintf("%s*", domain.TradePreKey)
	strs := r.redis.Keys(ctx, filter)

	for _, k := range strs.Val() {
		trade, err := r.FindByKey(ctx, k)
		if err != nil {
			continue
		}

		if trade.SellDate != nil {
			continue
		}

		trades = append(trades, *trade)
	}

	return trades, nil
}

// Save implements trade.Repository
func (r *TradeRepo) Save(ctx context.Context, trade *domain.Trade) (*domain.Trade, error) {
	if trade.ID == nil {
		trade.ID = pointer.ToString(uuid.New().String())
	}

	key, err := trade.BuildKey()
	if err != nil {
		return nil, err
	}

	bs, err := json.Marshal(trade)
	if err != nil {
		return nil, err
	}

	status := r.redis.Set(ctx, key, string(bs), -1)
	if status.Err() != nil {
		return nil, status.Err()
	}

	return trade, nil

}

func NewTradeRepository(cfg RedisCfg) trade.Repository {
	opts := buildOptions(cfg)
	redisCon := redis.NewClient(opts)

	checkRedisConnection(redisCon)

	return &TradeRepo{redis: redisCon}
}

func checkRedisConnection(redisCon *redis.Client) {
	logCtx := log.WithFields(log.Fields{"component": "TradeRepo", "function": "checkRedisConnection"})
	logCtx.Info("testing redis database connection...")
	status := redisCon.Ping(context.Background())
	if status.Err() != nil {
		logCtx.Fatalf("cannot connect with redis database: %v", status.Err())
	}
}

func buildOptions(cfg RedisCfg) *redis.Options {
	addr := fmt.Sprintf("%s:%d", cfg.Address, cfg.Port)
	opts := &redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	if cfg.UseTLS {
		opts.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return opts
}
