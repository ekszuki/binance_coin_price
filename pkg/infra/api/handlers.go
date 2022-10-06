package api

import (
	"fmt"
	"net/http"

	"br.com.sygnux/binance/pkg/domain"
	"br.com.sygnux/binance/pkg/logic"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (s *Server) getFindTradeByKey() gin.HandlerFunc {
	logCtx := logrus.WithFields(logrus.Fields{"component": "handlers", "function": "getFindTradeByKey"})
	return func(c *gin.Context) {
		key := c.Param("key")
		if key == "" {
			logCtx.Warnf("invalid key")
			c.JSON(
				http.StatusBadRequest,
				Response{
					IsError: true,
					Message: "invalid request",
				},
			)
			return
		}

		key = fmt.Sprintf("%s%s", domain.TradePreKey, key)

		trade, err := s.Repositories.TradeRepo.FindByKey(s.ctx, key)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				Response{
					IsError: true,
					Message: "something went wrong  :(",
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			Response{
				IsError: false,
				Resp:    trade,
			},
		)
	}
}

func (s *Server) postRegisterTrade() gin.HandlerFunc {
	logCtx := logrus.WithFields(logrus.Fields{"component": "handlers", "function": "registerTrade"})
	return func(c *gin.Context) {
		payload := new(domain.Trade)

		if err := c.ShouldBind(&payload); err != nil {
			logCtx.Warnf("invalid payload: %v", err)
			c.JSON(
				http.StatusBadRequest,
				Response{
					IsError: true,
					Message: "invalid payload",
				},
			)
			return
		}

		resp, err := logic.RegisterTrade(
			s.ctx,
			s.Repositories.TradeRepo,
			payload,
		)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				Response{
					IsError: true,
					Message: "something went wrong  :(",
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			Response{
				IsError: false,
				Resp:    resp,
			},
		)
	}
}
