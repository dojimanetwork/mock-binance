package store

import (
	tType "github.com/binance-chain/go-sdk/types/tx"
	. "gopkg.in/check.v1"
)

type InMemorySuite struct{}

var _ = Suite(&InMemorySuite{})

func (s *InMemorySuite) TestInMemory(c *C) {
	m := NewInMemory()
	c.Assert(m.Blocks, HasLen, 0)
	c.Assert(m.Accounts, HasLen, 0)

	c.Check(m.CurrentHeight(), Equals, 0)
	c.Check(m.GetBlock(0).Height, Equals, 0)
	c.Check(m.LastBlock().Height, Equals, 0)

	tx := tType.StdTx{Source: 33}
	m.AddBlock([]tType.StdTx{tx})
	tx = tType.StdTx{Source: 55}
	m.AddBlock([]tType.StdTx{tx})
	c.Assert(m.Blocks, HasLen, 2)

	c.Check(m.CurrentHeight(), Equals, 2)
	c.Check(m.GetBlock(1).Height, Equals, 1)
	c.Check(m.GetBlock(1).Txs[0].Source, Equals, int64(33))
	c.Check(m.LastBlock().Height, Equals, 2)
	c.Check(m.GetBlock(2).Txs[0].Source, Equals, int64(55))
}
