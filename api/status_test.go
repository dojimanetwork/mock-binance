package api

import (
	"encoding/json"

	. "gopkg.in/check.v1"

	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

type StatusSuite struct{}

var _ = Suite(&StatusSuite{})

func (s *StatusSuite) TestStatus(c *C) {
	bstore := store.NewInMemory()

	api := API(bstore)

	w := performRequest(api, "GET", "/status", nil)

	var result Status
	err := json.Unmarshal([]byte(w.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(200, Equals, w.Code)
	c.Assert(result.Result.NodeInfo.Network, Equals, "Binance-Chain-Ganges")
}
