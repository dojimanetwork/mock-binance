package api

import (
	"encoding/json"
	"fmt"

	. "gopkg.in/check.v1"

	txType "github.com/binance-chain/go-sdk/types/tx"

	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

type BlockSuite struct{}

var _ = Suite(&BlockSuite{})

type ErrorBlock struct {
	Error struct {
		Code    int64  `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error"`
}

func (s *BlockSuite) TestBlock(c *C) {
	bstore := store.NewInMemory()
	bstore.AddBlock([]txType.StdTx{txType.StdTx{}})

	api := API(bstore)

	w := performRequest(api, "GET", "/block", nil)

	var result RPCBlock
	err := json.Unmarshal([]byte(w.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(200, Equals, w.Code)
	height := fmt.Sprintf("%d", bstore.CurrentHeight())
	c.Assert(result.Result.Block.Header.Height, Equals, height)

	w = performRequest(api, "GET", "/block?height=1", nil)
	err = json.Unmarshal([]byte(w.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(200, Equals, w.Code)
	c.Assert(result.Result.Block.Header.Height, Equals, "1")

	var errorBlock ErrorBlock
	w = performRequest(api, "GET", "/block?height=8", nil)
	err = json.Unmarshal([]byte(w.Body.String()), &errorBlock)
	c.Assert(err, IsNil)
	c.Assert(200, Equals, w.Code)
	c.Assert(errorBlock.Error.Code, Equals, int64(-32603))
	c.Assert(errorBlock.Error.Message, Equals, "Internal error")
}
