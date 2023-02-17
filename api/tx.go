package api

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"

	tType "github.com/binance-chain/go-sdk/types/tx"
	"github.com/gin-gonic/gin"

	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

var re = regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

type tx struct {
	Hash   string `json:"hash"`
	Height string `json:"height"`
	Tx     string `json:"tx"`
}

type RPCTxSearch struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		Txs        []tx   `json:"txs"`
		TotalCount string `json:"total_count"`
	} `json:"result"`
}

// Returns tx search
func txSearch(bstore store.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		txHeight := c.DefaultQuery("query", "tx.height=0")

		submatchall := re.FindAllString(txHeight, -1)
		if len(submatchall) == 0 {
			c.JSON(400, gin.H{"error": "could not find tx height"})
			return
		}

		height, _ := strconv.Atoi(submatchall[0])
		block := bstore.GetBlock(height)

		var txSearch RPCTxSearch
		for _, blockTx := range block.Txs {
			bin, _ := tType.Cdc.MarshalBinaryLengthPrefixed(blockTx)
			encodedTx := base64.StdEncoding.EncodeToString(bin)

			txn := tx{
				Hash:   block.Hash,
				Height: txHeight,
				Tx:     encodedTx,
			}
			txSearch.Result.Txs = append(txSearch.Result.Txs, txn)
		}
		txSearch.Result.TotalCount = fmt.Sprintf("%d", len(txSearch.Result.Txs))

		c.JSON(200, txSearch)
	}
}
