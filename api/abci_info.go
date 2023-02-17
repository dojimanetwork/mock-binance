package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

type info struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		Response struct {
			BlockHeight string `json:"last_block_height"`
		} `json:"response"`
	} `json:"result"`
}

// Returns ABCIinfo
func ABCIinfo(bstore store.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		var b info
		b.Result.Response.BlockHeight = fmt.Sprintf("%d", bstore.CurrentHeight())
		c.JSON(200, b)
	}
}
