package store

import (
	"github.com/binance-chain/go-sdk/common/types"
	tType "github.com/binance-chain/go-sdk/types/tx"
)

type Store interface {
	// Blocks
	CurrentHeight() int
	GetBlock(height int) Block
	ListBlocks() Blocks
	AddBlock(tx []tType.StdTx)
	LastBlock() Block

	// Accounts
	ListAccounts() Accounts
	GetAccount(addr types.AccAddress) Account
	SetAccount(acc Account)
	Transfer(from, to types.AccAddress, coins types.Coins)
}
