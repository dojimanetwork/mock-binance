package api

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	ctypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types"
	"github.com/gin-gonic/gin"

	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

type queryResult struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  struct {
		Response struct {
			Key         string `json:"key"`
			Value       string `json:"value"`
			BlockHeight string `json:"height"`
		} `json:"response"`
	} `json:"result"`
}

// Returns ABCIquery
func ABCIquery(bstore store.Store) func(*gin.Context) {
	return func(c *gin.Context) {
		path := c.Query("path")

		// remove quotes from around path
		if len(path) > 0 && path[0] == '"' {
			path = path[1:]
		}
		if len(path) > 0 && path[len(path)-1] == '"' {
			path = path[0 : len(path)-1]
		}

		data := c.Query("data")
		// remove 0x from beginning of data
		if len(data) > 0 && strings.HasPrefix(data, "0x") {
			data = data[2:]
		}
		bz, err := hex.DecodeString(data)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if path == "/store/acc/key" {
			parts := strings.Split(string(bz), ":")
			if len(parts) != 2 {
				c.JSON(400, gin.H{"error": "invalid data string"})
				return
			}

			// TODO: don't hardcode test network
			ctypes.Network = ctypes.TestNetwork
			address := ctypes.AccAddress([]byte(parts[1]))
			acc := bstore.GetAccount(address)

			var account ctypes.AppAccount
			account.BaseAccount.Address = address
			account.BaseAccount.Coins = acc.Balances
			account.BaseAccount.AccountNumber = acc.AccountNumber
			account.BaseAccount.Sequence = acc.Sequence

			cdc := types.NewCodec()
			bz, err := cdc.MarshalBinaryBare(account)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			encoded := base64.StdEncoding.EncodeToString(bz)

			var result queryResult
			result.Result.Response.Key = data
			result.Result.Response.Value = encoded
			result.Result.Response.BlockHeight = fmt.Sprintf("%d", bstore.CurrentHeight())

			c.JSON(200, result)
		} else if strings.HasPrefix(path, "/account/") {
			parts := strings.Split(path, "/")
			addr := parts[len(parts)-1]

			ctypes.Network = ctypes.TestNetwork
			address, err := ctypes.AccAddressFromBech32(addr)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			acc := bstore.GetAccount(address)

			var account ctypes.AppAccount
			var base ctypes.BaseAccount
			base.Address = address
			base.Coins = acc.Balances
			base.AccountNumber = acc.AccountNumber
			base.Sequence = acc.Sequence
			account.BaseAccount = base

			cdc := types.NewCodec()
			bz, err := cdc.MarshalBinaryBare(account)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			encoded := base64.StdEncoding.EncodeToString(bz)

			var result queryResult
			result.Result.Response.Value = encoded
			c.JSON(200, result)
		} else if path == "/param/fees" {
			var result queryResult
			result.Result.Response.Value = "ngUKHcKpb6MKD3N1Ym1pdF9wcm9wb3NhbBCAyrXuARgBChPCqW+jCgdkZXBvc2l0EKToAxgBCgzCqW+jCgR2b3RlGAMKHsKpb6MKEGNyZWF0ZV92YWxpZGF0b3IQgJTr3AMYAQodwqlvowoQcmVtb3ZlX3ZhbGlkYXRvchCAwtcvGAEKFsKpb6MKB2RleExpc3QQgNDbw/QCGAIKEMKpb6MKCG9yZGVyTmV3GAMKE8Kpb6MKC29yZGVyQ2FuY2VsGAMKF8Kpb6MKCGlzc3VlTXNnEIDo7aG6ARgCChXCqW+jCgdtaW50TXNnEIDKte4BGAIKF8Kpb6MKCnRva2Vuc0J1cm4QgOHrFxgBChjCqW+jCgx0b2tlbnNGcmVlemUQoMIeGAEKGJo9J2kKDAoEc2VuZBD8pAIYARCw6gEYAgqgAUlaUEQKDwoJRXhwaXJlRmVlEKjDAQoUCg9FeHBpcmVGZWVOYXRpdmUQiCcKDwoJQ2FuY2VsRmVlEKjDAQoUCg9DYW5jZWxGZWVOYXRpdmUQiCcKDAoHRmVlUmF0ZRDoBwoSCg1GZWVSYXRlTmF0aXZlEJADChEKDElPQ0V4cGlyZUZlZRCQTgoXChJJT0NFeHBpcmVGZWVOYXRpdmUQxBMKFMKpb6MKCHRpbWVMb2NrEMCEPRgBChbCqW+jCgp0aW1lVW5sb2NrEMCEPRgBChbCqW+jCgp0aW1lUmVsb2NrEMCEPRgBChzCqW+jCg9zZXRBY2NvdW50RmxhZ3MQgMLXLxgBChDCqW+jCgRIVExUEPykAhgBChfCqW+jCgtkZXBvc2l0SFRMVBD8pAIYAQoVwqlvowoJY2xhaW1IVExUEPykAhgBChbCqW+jCgpyZWZ1bmRIVExUEPykAhgB"
			c.JSON(200, result)
		} else {
			c.JSON(500, gin.H{"error": fmt.Sprintf("Invalid path: %s", path)})
		}
	}
}
