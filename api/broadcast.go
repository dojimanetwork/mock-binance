package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	ctypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types/msg"
	txType "github.com/binance-chain/go-sdk/types/tx"
	"github.com/gin-gonic/gin"

	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

var SEED = "tbnb1ht7v08hv2lhtmk8y7szl2hjexqryc3hcldlztl"

type Log struct {
	Success bool   `json:"success"`
	Log     string `json:"log"`
}

type TxCommitResultSuccess struct {
	Hash   string `json:"hash"`
	Height string `json:"height"`
	Logs   []Log  `json:"logs"`
}

type TxCommitResult struct {
	Ok   bool   `json:"ok"`
	Log  string `json:"log"`
	Hash string `json:"hash"`
	Code int32  `json:"code"`
	Data string `json:"data"`
}

func broadcast(bstore store.Store, mutex *sync.Mutex) func(*gin.Context) {
	return func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()
		ctypes.Network = ctypes.TestNetwork

		txString := c.DefaultQuery("tx", "")
		if txString == "" {
			body, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
				return
			}
			txString = string(body)
		}

		// Remove quotes around our string (if present)
		if len(txString) > 0 && txString[0] == '"' {
			txString = txString[1:]
		}
		if len(txString) > 0 && txString[len(txString)-1] == '"' {
			txString = txString[:len(txString)-1]
		}

		// Remove 0x prefix if exists
		if strings.HasPrefix(txString, "0x") {
			txString = txString[2:]
		}

		buf, err := hex.DecodeString(txString)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		var t txType.StdTx
		if err := txType.Cdc.UnmarshalBinaryLengthPrefixed(buf, &t); nil != err {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		msg := t.Msgs[0].(msg.SendMsg)

		input := msg.Inputs[0]
		output := msg.Outputs[0]
		coins := output.Coins

		from := bstore.GetAccount(input.Address)
		if from.Sequence != t.Signatures[0].Sequence {
			fmt.Printf("Bad Seq: %d %d \n", from.Sequence, t.Signatures[0].Sequence)
			fmt.Printf("Message: %+v\n", t)
			result := TxCommitResult{
				Ok:   true,
				Code: 4,
			}
			c.JSON(200, result)
			return
		}
		from.Sequence++
		from.SubCoins(coins)

		// minus gas...
		gas := 37500
		if len(coins) > 1 {
			gas = 30000 * len(coins)
		}
		gasCoin := ctypes.Coin{Denom: "BNB", Amount: int64(gas)}
		from.SubCoins(ctypes.Coins{gasCoin})

		// if the from address isn't the SEED address, check we have enough funds
		if from.Address.String() != SEED {
			if err := from.Valid(); err != nil {
				fmt.Printf("Invalid Account: %s: %s\n", from.Address, err.Error())
				result := TxCommitResult{
					Ok:   true,
					Code: 5,
					Log:  err.Error(),
				}
				c.JSON(200, result)
				return
			}
		}

		bstore.SetAccount(from)

		to := bstore.GetAccount(output.Address)
		to.AddCoins(coins)
		bstore.SetAccount(to)

		fmt.Printf("Transaction: (Acct: %d | Seq: %d) %s ==> %s (%s) | %+v\n", t.Signatures[0].AccountNumber, t.Signatures[0].Sequence, input.Address.String(), output.Address.String(), t.Memo, coins)

		bstore.AddBlock([]txType.StdTx{t})
		block := bstore.LastBlock()

		result := TxCommitResultSuccess{
			Hash:   block.Hash,
			Height: fmt.Sprintf("%d", block.Height),
			Logs: []Log{
				{Success: true, Log: "OK"},
			},
		}

		c.JSON(http.StatusOK, result)
	}
}

type BroadcastEasy struct {
	From  ctypes.AccAddress `json:"from"`
	To    ctypes.AccAddress `json:"to"`
	Coins ctypes.Coins      `json:"coins"`
	Memo  string            `json:"memo"`
}

func broadcastEasy(bstore store.Store, mutex *sync.Mutex) func(*gin.Context) {
	return func(c *gin.Context) {
		mutex.Lock()
		defer mutex.Unlock()
		ctypes.Network = ctypes.TestNetwork
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		var txs []BroadcastEasy
		if err := json.Unmarshal(body, &txs); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		stdTxs := make([]txType.StdTx, 0)

		for _, tx := range txs {

			from := bstore.GetAccount(tx.From)
			// from.Sequence = from.Sequence + 1
			from.SubCoins(tx.Coins)

			// minus gas...
			gas := 37500
			if len(tx.Coins) > 1 {
				gas = 30000 * len(tx.Coins)
			}
			gasCoin := ctypes.Coin{Denom: "BNB", Amount: int64(gas)}
			from.Sequence++
			from.SubCoins(ctypes.Coins{gasCoin})

			bstore.SetAccount(from)

			to := bstore.GetAccount(tx.To)
			to.AddCoins(tx.Coins)
			bstore.SetAccount(to)

			fmt.Printf("Transaction: %s ==> %s (%s) | %+v\n", tx.From.String(), tx.To.String(), tx.Memo, tx.Coins)

			m := msg.NewMsgSend(
				[]msg.Input{msg.NewInput(tx.From, tx.Coins)},
				[]msg.Output{msg.NewOutput(tx.To, tx.Coins)},
			)
			sig := txType.StdSignature{AccountNumber: from.AccountNumber, Sequence: from.Sequence}

			t := txType.NewStdTx([]msg.Msg{m}, []txType.StdSignature{sig}, tx.Memo, 0, nil)
			stdTxs = append(stdTxs, t)
		}
		bstore.AddBlock(stdTxs)
		block := bstore.LastBlock()

		c.JSON(200, block)
	}
}
