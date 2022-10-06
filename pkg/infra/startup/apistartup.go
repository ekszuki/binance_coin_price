package startup

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"br.com.sygnux/binance/pkg/infra/api"
	"br.com.sygnux/binance/pkg/infra/repository/binance"
	"br.com.sygnux/binance/pkg/infra/repository/redis"
	"br.com.sygnux/binance/pkg/logic"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type APPConfig struct {
	RedisPassword   string        `mapstructure:"REDIS_PASSWORD"`
	RedisAddress    string        `mapstructure:"REDIS_ADDRESS"`
	RedisPort       int           `mapstructure:"REDIS_PORT"`
	RedisDB         int           `mapstructure:"REDIS_DB"`
	RedisUseTLS     bool          `mapstructure:"REDIS_USETLS"`
	BinanceWaitTime time.Duration `mapstructure:"BINANCE_WAITTIME"`
	BinancePriceURL string        `mapstructure:"BINANCE_PRICE_URL"`
	HTTPPort        int           `mapstructure:"HTTP_PORT"`
}

func (a *APPConfig) Validate() error {
	if a.BinanceWaitTime == 0 {
		return fmt.Errorf("binance waiting time not defined")
	}

	if a.HTTPPort == 0 {
		return fmt.Errorf("HTTP Port not defined")
	}

	if strings.TrimSpace(a.BinancePriceURL) == "" {
		return fmt.Errorf("binance price URL not defined")
	}

	return nil
}

type APIStartup struct {
	ctx context.Context
}

func NewAPIStartup(ctx context.Context) *APIStartup {
	return &APIStartup{ctx: ctx}
}

func (s *APIStartup) loadConfigEnvs() (*APPConfig, error) {
	appConfig := new(APPConfig)

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return appConfig, err
	}

	err = viper.Unmarshal(appConfig)

	return appConfig, err
}

func (s *APIStartup) Initialize() {
	time.Local = time.UTC
	logCtx := log.WithFields(log.Fields{"component": "APIStartup", "function": "Initialize"})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	logCtx.Info("loading environments")
	appConfig, err := s.loadConfigEnvs()
	if err != nil {
		logCtx.Fatalf("cannot load app env file: %v", err)
	}

	err = appConfig.Validate()
	if err != nil {
		logCtx.Fatalf("env validation error: %v", err)
	}

	logCtx.Info("create Redis connection")

	redisCFG := redis.RedisCfg{
		Password: appConfig.RedisPassword,
		Address:  appConfig.RedisAddress,
		Port:     appConfig.RedisPort,
		DB:       appConfig.RedisDB,
		UseTLS:   appConfig.RedisUseTLS,
	}
	tradeRepo := redis.NewTradeRepository(redisCFG)
	binanceCoin := binance.NewCoinRepository(
		s.ctx, http.DefaultClient, appConfig.BinancePriceURL,
	)

	apiRepo := &api.Repositories{
		TradeRepo:   tradeRepo,
		BinanceCoin: binanceCoin,
	}

	apiServer := api.NewServer(s.ctx, apiRepo)
	apiServer.Run(api.ServerCfg{Port: appConfig.HTTPPort})

	// start real time monitor
	logic.RealTimeBinanceCheck(s.ctx, tradeRepo, binanceCoin, appConfig.BinanceWaitTime)

	<-c

	timeOutCtx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	err = apiServer.Shutdown(timeOutCtx)
	if err != nil {
		logCtx.Warnf("cannot close api server gracefully: %v", err)
	}

	logCtx.Info("api server shutdown.... bye")
}
