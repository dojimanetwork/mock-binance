package api

import (
	"github.com/gin-gonic/gin"
)

type Status struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		NodeInfo struct {
			Network string `json:"network"`
		} `json:"node_info"`
	} `json:"result"`
}

func status() func(*gin.Context) {
	return func(c *gin.Context) {
		var b Status
		b.Result.NodeInfo.Network = "Binance-Chain-Ganges"
		c.JSON(200, b)
	}
}
