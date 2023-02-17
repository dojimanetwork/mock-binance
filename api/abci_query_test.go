package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	. "gopkg.in/check.v1"

	ctypes "github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/types"
	"github.com/tendermint/go-amino"
	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

type ABCIquerySuite struct{}

var _ = Suite(&ABCIquerySuite{})

func (s *ABCIquerySuite) SetUpSuite(c *C) {
	ctypes.Network = ctypes.TestNetwork
}

func (s *ABCIquerySuite) TestABCIFees(c *C) {
	var err error
	bstore := store.NewInMemory()
	api := API(bstore)

	path := fmt.Sprintf("/abci_query?path=\"/param/fees\"")
	w := performRequest(api, "GET", path, nil)

	var result queryResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	c.Assert(err, IsNil)
	c.Assert(200, Equals, w.Code)

	data, err := base64.StdEncoding.DecodeString(result.Result.Response.Value)
	c.Assert(err, IsNil)

	var fees []ctypes.FeeParam
	cdc := amino.NewCodec()
	ctypes.RegisterWire(cdc)
	err = cdc.UnmarshalBinaryLengthPrefixed(data, &fees)
	c.Assert(err, IsNil)
	c.Assert(fees, HasLen, 22, Commentf("%d", len(fees)))
}

func (s *ABCIquerySuite) TestABCI(c *C) {
	var err error
	bstore := store.NewInMemory()
	api := API(bstore)

	var newAcct store.Account
	newAcct.Address, err = ctypes.AccAddressFromBech32("tbnb1sylyjw032eajr9cyllp26n04300qzzre38qyv5")
	c.Assert(err, IsNil)
	newAcct.Balances = ctypes.Coins{
		ctypes.Coin{Denom: "BNB", Amount: 248},
		ctypes.Coin{Denom: "RUNE", Amount: 3048},
	}
	bstore.SetAccount(newAcct)

	key := append([]byte("account:"), newAcct.Address.Bytes()...)
	path := fmt.Sprintf("/abci_query?path=\"/store/acc/key\"&data=0x%x", key)
	w := performRequest(api, "GET", path, nil)

	var result queryResult
	err = json.Unmarshal(w.Body.Bytes(), &result)
	c.Assert(err, IsNil)
	c.Assert(200, Equals, w.Code)
	c.Check(result.Result.Response.Key, Equals, fmt.Sprintf("%x", key))

	data, err := base64.StdEncoding.DecodeString(result.Result.Response.Value)
	c.Assert(err, IsNil)

	cdc := types.NewCodec()
	var acc ctypes.AppAccount
	err = cdc.UnmarshalBinaryBare(data, &acc)
	c.Assert(err, IsNil)

	c.Check(acc.BaseAccount.Address.String(), Equals, newAcct.Address.String())
	c.Assert(acc.BaseAccount.Coins, HasLen, 2)
	for i := range acc.BaseAccount.Coins {
		c.Check(acc.BaseAccount.Coins[i].Denom, Equals, newAcct.Balances[i].Denom)
		c.Check(acc.BaseAccount.Coins[i].Amount, Equals, newAcct.Balances[i].Amount)
	}
}
