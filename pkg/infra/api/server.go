package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"br.com.sygnux/binance/pkg/domain/binance"
	"br.com.sygnux/binance/pkg/domain/trade"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ServerCfg struct {
	Port int
}

type Server struct {
	ctx          context.Context
	httpServer   *http.Server
	Repositories *Repositories
}

type Repositories struct {
	TradeRepo   trade.Repository
	BinanceCoin binance.CoinRepository
}

func NewServer(ctx context.Context, repos *Repositories) *Server {
	return &Server{
		ctx:          ctx,
		Repositories: repos,
	}
}

func (s *Server) Run(cfg ServerCfg) {
	logCtx := logrus.WithFields(logrus.Fields{"component": "Server", "function": "Run"})
	engine := s.setupGin()
	addr := fmt.Sprintf(":%d", cfg.Port)
	s.httpServer = &http.Server{
		Addr:              addr,
		Handler:           engine,
		ReadHeaderTimeout: (5 * time.Second),
	}

	go func() {
		logCtx.Infof("starting http server on address: %s", addr)
		err := s.httpServer.ListenAndServe()
		if err != nil {
			logCtx.Fatalf("cannot start http server: %v", err)
		}
	}()
}

func (s *Server) Shutdown(ctxTimeOut context.Context) error {
	return s.httpServer.Shutdown(ctxTimeOut)
}

func (s *Server) setupGin() *gin.Engine {
	r := gin.New()

	corsCFG := createCorsConfig()
	r.Use(gin.Recovery(), cors.New(corsCFG), gin.Logger())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong", "time": time.Now().Format("02/01/2006 15:04:05")})
	})

	v1 := r.Group("/v1")
	v1.GET("/trade/:key", s.getFindTradeByKey())
	v1.POST("/trade", s.postRegisterTrade())

	return r
}

func createCorsConfig() cors.Config {
	corsConf := cors.DefaultConfig()
	corsConf.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Api-Key", "Refresh-Token"}
	corsConf.ExposeHeaders = []string{"refresh-token"}
	corsConf.AllowOrigins = []string{"*"}

	return corsConf
}
