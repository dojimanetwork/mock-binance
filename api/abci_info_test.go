package api

import (
	"encoding/json"
	"fmt"

	. "gopkg.in/check.v1"

	txType "github.com/binance-chain/go-sdk/types/tx"

	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

type ABCIinfoSuite struct{}

var _ = Suite(&ABCIinfoSuite{})

func (s *ABCIinfoSuite) TestABCI(c *C) {
	bstore := store.NewInMemory()
	bstore.AddBlock([]txType.StdTx{txType.StdTx{}})

	api := API(bstore)

	w := performRequest(api, "GET", "/abci_info", nil)

	var result info
	err := json.Unmarshal(w.Body.Bytes(), &result)
	c.Assert(err, IsNil)
	c.Assert(200, Equals, w.Code)
	height := fmt.Sprintf("%d", bstore.CurrentHeight())
	c.Assert(result.Result.Response.BlockHeight, Equals, height)
}
