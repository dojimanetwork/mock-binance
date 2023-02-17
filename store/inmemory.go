package store

import (
	"github.com/binance-chain/go-sdk/common/types"
	tType "github.com/binance-chain/go-sdk/types/tx"
)

type InMemory struct {
	Store
	Blocks   Blocks   `json:"blocks"`
	Accounts Accounts `json:"accounts"`
}

func NewInMemory() *InMemory {
	return &InMemory{
		Blocks:   make(Blocks, 0),
		Accounts: make(Accounts, 0),
	}
}

func (m *InMemory) CurrentHeight() int {
	return len(m.Blocks)
}

func (m *InMemory) ListBlocks() Blocks {
	return m.Blocks
}

func (m *InMemory) GetBlock(height int) Block {
	if m.CurrentHeight() < height || height < 1 {
		return Block{}
	}
	return m.Blocks[height-1]
}

func (m *InMemory) AddBlock(txs []tType.StdTx) {
	block := NewBlock(m.CurrentHeight()+1, txs)
	m.Blocks = append(m.Blocks, block)
}

func (m *InMemory) LastBlock() Block {
	if len(m.Blocks) > 0 {
		return m.Blocks[len(m.Blocks)-1]
	}
	return Block{}
}

func (m *InMemory) ListAccounts() Accounts {
	return m.Accounts
}

func (m *InMemory) GetAccount(addr types.AccAddress) Account {
	for _, acc := range m.Accounts {
		if addr.String() == acc.Address.String() {
			// make a copy of the coins
			coins := types.Coins{}
			for _, item := range acc.Balances {
				coins = append(coins, types.Coin{
					Denom:  item.Denom,
					Amount: item.Amount,
				})
			}
			acc.Balances = coins
			return acc
		}
	}
	return Account{Address: addr}
}

func (m *InMemory) SetAccount(acc Account) {
	for i, _ := range m.Accounts {
		if acc.Address.String() == m.Accounts[i].Address.String() {
			m.Accounts[i] = acc
			return
		}
	}
	acc.AccountNumber = int64(len(m.Accounts) + 1)
	m.Accounts = append(m.Accounts, acc)
}

func (m *InMemory) Transfer(from, to types.AccAddress, coins types.Coins) {
	fromAcc := m.GetAccount(from)
	toAcc := m.GetAccount(to)

	fromAcc.SubCoins(coins)
	toAcc.AddCoins(coins)

	m.SetAccount(fromAcc)
	m.SetAccount(toAcc)
}
