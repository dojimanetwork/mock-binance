package api

import (
	ctypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/gin-gonic/gin"
	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

// Balance Account definition
// https://github.com/binance-chain/go-sdk/blob/master/common/types/account.go#L178
type BalanceAccount struct {
	Number    int64   `json:"account_number"`
	Address   string  `json:"address"`
	PublicKey []uint8 `json:"public_key"`
	Sequence  int64   `json:"sequence"`
	Flags     uint64  `json:"flags"`
	// Balances  []TokenBalance `json:"balances"`
}

func acc_info(bstore store.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		name := c.Param("acc")
		var b BalanceAccount

		address, err := ctypes.AccAddressFromBech32(name)
		if err != nil {
			c.JSON(400, gin.H{"Error": err.Error()})
		}

		acc := bstore.GetAccount(address)

		b.Sequence = acc.Sequence
		b.Number = acc.AccountNumber
		c.JSON(200, b)
	}
}

func acc_list(bstore store.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, bstore.ListAccounts())
	}
}
