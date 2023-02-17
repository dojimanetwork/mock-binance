package store

import (
	. "gopkg.in/check.v1"

	tType "github.com/binance-chain/go-sdk/types/tx"
)

type BlockSuite struct{}

var _ = Suite(&BlockSuite{})

func (s *BlockSuite) TestBlock(c *C) {
	tx := tType.StdTx{Source: 33}
	b := NewBlock(12, []tType.StdTx{tx})
	c.Check(b.Height, Equals, 12)
	c.Check(b.Txs[0].Source, Equals, int64(33))
	c.Check(b.Hash, HasLen, 64)
}
