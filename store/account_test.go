package store

import (
	. "gopkg.in/check.v1"

	"github.com/binance-chain/go-sdk/common/types"
)

type AccountSuite struct{}

var _ = Suite(&AccountSuite{})

func (s *AccountSuite) TestAccount(c *C) {
	bstore := NewInMemory()
	from, err := types.AccAddressFromBech32("bnb136ns6lfw4zs5hg4n85vdthaad7hq5m4gtkgf23")
	c.Assert(err, IsNil)
	to, err := types.AccAddressFromBech32("bnb1jxfh2g85q3v0tdq56fnevx6xcxtcnhtsmcu64m")
	c.Assert(err, IsNil)
	coins := types.Coins{
		types.Coin{Denom: "BNB", Amount: 23},
		types.Coin{Denom: "LOKI", Amount: 44},
	}

	// send coins (ok that from address has no coins already)
	bstore.Transfer(from, to, coins)
	fromAcc := bstore.GetAccount(from)
	toAcc := bstore.GetAccount(to)
	c.Check(fromAcc.Balances, HasLen, 0)
	c.Check(toAcc.Balances, HasLen, 2)
	c.Check(toAcc.Balances[0].Amount, Equals, int64(23))
	c.Check(toAcc.Balances[1].Amount, Equals, int64(44))

	// send coins back
	bstore.Transfer(to, from, coins)
	fromAcc = bstore.GetAccount(from)
	toAcc = bstore.GetAccount(to)
	c.Check(fromAcc.Balances, HasLen, 2)
	c.Check(fromAcc.Balances[0].Amount, Equals, int64(23))
	c.Check(fromAcc.Balances[1].Amount, Equals, int64(44))
	c.Check(toAcc.Balances, HasLen, 2)
	c.Check(toAcc.Balances[0].Amount, Equals, int64(0))
	c.Check(toAcc.Balances[1].Amount, Equals, int64(0))
}
