package store

import (
	tType "github.com/binance-chain/go-sdk/types/tx"

	"gitlab.com/thorchain/bepswap/mock-binance/util"
)

type Block struct {
	Height int           `json:"height"`
	Hash   string        `json:"hash"`
	Txs    []tType.StdTx `json:"std_tx"`
}

type Blocks []Block

func NewHash() string {
	return util.String(64)
}

func NewBlock(height int, txs []tType.StdTx) Block {
	return Block{
		Height: height,
		Txs:    txs,
		Hash:   NewHash(),
	}
}
