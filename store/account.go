package store

import (
	"fmt"

	"github.com/binance-chain/go-sdk/common/types"
)

type Account struct {
	Address       types.AccAddress `json:"address"`
	Balances      types.Coins      `json:"balances"`
	Sequence      int64            `json:"seq"`
	AccountNumber int64            `json:"acc_num"`
}

type Accounts []Account

func NewAccount(addr types.AccAddress) Account {
	return Account{
		Address: addr,
	}
}

func (a *Account) AddCoins(coins types.Coins) {
	for _, coin := range coins {
		found := false
		for i, bCoin := range a.Balances {
			if bCoin.Denom == coin.Denom {
				a.Balances[i].Amount += coin.Amount
				found = true
			}
		}
		if !found {
			a.Balances = append(a.Balances, coin)
		}
	}
}

func (a *Account) SubCoins(coins types.Coins) {
	for _, coin := range coins {
		for i, bCoin := range a.Balances {
			if bCoin.Denom == coin.Denom {
				a.Balances[i].Amount -= coin.Amount
			}
		}
	}
}

func (a *Account) Valid() error {
	for _, coin := range a.Balances {
		if coin.Amount < 0 {
			return fmt.Errorf("Negative Balance: %d %s", coin.Amount, coin.Denom)
		}
	}

	return nil
}
