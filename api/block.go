package api

import (
	"encoding/base64"
	"fmt"
	"strconv"

	tType "github.com/binance-chain/go-sdk/types/tx"
	"github.com/gin-gonic/gin"

	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

// TODO: remove me
type RPCBlock struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		Block struct {
			Header struct {
				Height string `json:"height"`
				NumTxs string `json:"num_txs"`
			} `json:"header"`
			Data struct {
				Txs []string `json:"txs"`
			} `json:"data"`
		} `json:"block"`
	} `json:"result"`
}

// Returns block
func block(bstore store.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		block := bstore.LastBlock()
		height := c.DefaultQuery("height", "")
		if height != "" {
			i, _ := strconv.Atoi(height)
			block = bstore.GetBlock(i)

			// check if our block doesn't exist yet
			if block.Height == 0 {
				c.JSON(200, gin.H{
					"jsonrpc": "2.0",
					"id":      "",
					"error": gin.H{
						"code":    -32603,
						"message": "Internal error",
						"data":    "Height must be less than or equal to the current blockchain height",
					},
				})
				return
			}
		}
		var b RPCBlock
		b.Result.Block.Header.Height = fmt.Sprintf("%d", block.Height)
		b.Result.Block.Header.NumTxs = "1" // all blocks have one transaction

		for _, tx := range block.Txs {
			bin, _ := tType.Cdc.MarshalBinaryLengthPrefixed(tx)
			encodedTx := base64.StdEncoding.EncodeToString(bin)
			b.Result.Block.Data.Txs = append(b.Result.Block.Data.Txs, encodedTx)
		}

		c.JSON(200, b)
	}
}

func list_blocks(bstore store.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, bstore.ListBlocks())
	}
}
