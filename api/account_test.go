package api

import (
	"encoding/json"
	"fmt"

	ctypes "github.com/binance-chain/go-sdk/common/types"
	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

type AccoutSuite struct{}

var _ = Suite(&AccoutSuite{})

func (s *AccoutSuite) TestAccount(c *C) {
	bstore := store.NewInMemory()
	addr, err := ctypes.AccAddressFromBech32("tbnb1sylyjw032eajr9cyllp26n04300qzzre38qyv5")
	c.Assert(err, IsNil)

	acc := store.NewAccount(addr)
	bstore.SetAccount(acc)

	api := API(bstore)

	w := performRequest(api, "GET", fmt.Sprintf("/account/%s", addr.String()), nil)

	var result BalanceAccount
	err = json.Unmarshal(w.Body.Bytes(), &result)
	c.Assert(err, IsNil)
	c.Assert(200, Equals, w.Code)
	c.Check(result.Number, Equals, int64(1))
	c.Check(result.Sequence, Equals, int64(0))

	a1 := bstore.GetAccount(addr)
	a1.AddCoins(ctypes.Coins{
		ctypes.Coin{
			Denom:  "BNB",
			Amount: 1024,
		},
	})
	bstore.SetAccount(a1)
	a2 := bstore.GetAccount(addr)
	a2.SubCoins(ctypes.Coins{
		ctypes.Coin{
			Denom:  "BNB",
			Amount: 512,
		},
	})

	a3 := bstore.GetAccount(addr)
	c.Assert(a3.Balances[0].Amount, Equals, a1.Balances[0].Amount)
}
