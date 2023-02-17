package api

import (
	"encoding/json"
	"fmt"

	. "gopkg.in/check.v1"

	txType "github.com/binance-chain/go-sdk/types/tx"

	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

type TxSuite struct{}

var _ = Suite(&TxSuite{})

func (s *TxSuite) TestTx(c *C) {
	bstore := store.NewInMemory()
	bstore.AddBlock([]txType.StdTx{txType.StdTx{}})

	api := API(bstore)

	uri := fmt.Sprintf("/tx_search?tx.height=%s", bstore.LastBlock().Hash)
	w := performRequest(api, "GET", uri, nil)

	var result RPCTxSearch
	err := json.Unmarshal([]byte(w.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(200, Equals, w.Code)
}
