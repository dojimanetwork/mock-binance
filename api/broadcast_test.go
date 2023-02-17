package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	. "gopkg.in/check.v1"

	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	mType "github.com/binance-chain/go-sdk/types/msg"
	txType "github.com/binance-chain/go-sdk/types/tx"

	"gitlab.com/thorchain/bepswap/mock-binance/store"
)

type BroadcastSuite struct{}

var _ = Suite(&BroadcastSuite{})

func (s *BroadcastSuite) TestBroadcast(c *C) {
	types.Network = types.TestNetwork
	from, err := types.AccAddressFromBech32("tbnb1fds7yhw7qt9rkxw9pn65jyj004x858ny4xf2dk")
	c.Assert(err, IsNil)
	to, err := types.AccAddressFromBech32("tbnb1yxfyeda8pnlxlmx0z3cwx74w9xevspwdpzdxpj")
	c.Assert(err, IsNil)
	coins := types.Coins{types.Coin{Amount: 194765912}}
	transfers := []mType.Transfer{
		mType.Transfer{ToAddr: to, Coins: coins},
	}
	msg := mType.CreateSendMsg(from, coins, transfers)

	stdTx := txType.NewStdTx(
		[]mType.Msg{msg},
		nil,
		"OUTBOUND:1718",
		0,
		nil,
	)

	keyManager, err := keys.NewPrivateKeyManager("91a2f0e5b1495cf51b0792a009b49c54ce8ae52d0dada711e73d98b22e6698ea")
	c.Assert(err, IsNil)

	hexTx, _, err := SignTx(keyManager, stdTx, 0)
	c.Assert(err, IsNil)

	bstore := store.NewInMemory()
	api := API(bstore)

	w := performRequest(api, "POST", "/broadcast", bytes.NewReader(hexTx))

	block := bstore.LastBlock()

	var result TxCommitResultSuccess
	err = json.Unmarshal([]byte(w.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(200, Equals, w.Code)
	c.Assert(result.Logs[0].Success, Equals, true)
	c.Assert(result.Hash, HasLen, 64)
	c.Assert(block.Height, Equals, 1)
	c.Assert(block.Hash, Equals, result.Hash)
	c.Assert(block.Txs[0].Memo, Equals, stdTx.Memo)
	c.Assert(block.Txs[0].Msgs[0].GetSignBytes(), DeepEquals, stdTx.Msgs[0].GetSignBytes())

	hexTx, _, err = SignTx(keyManager, stdTx, 1)
	c.Assert(err, IsNil)

	fmt.Printf("hexString: %s\n", string(hexTx))
	q := fmt.Sprintf("/broadcast?tx=\"0x%s\"", string(hexTx))
	w = performRequest(api, "POST", q, nil)

	block = bstore.LastBlock()

	err = json.Unmarshal([]byte(w.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(200, Equals, w.Code)
	c.Assert(result.Logs[0].Success, Equals, true)
	c.Assert(result.Hash, HasLen, 64)
	c.Assert(block.Height, Equals, 2)
	c.Assert(block.Hash, Equals, result.Hash)
	c.Assert(block.Txs[0].Memo, Equals, stdTx.Memo)
	c.Assert(block.Txs[0].Msgs[0].GetSignBytes(), DeepEquals, stdTx.Msgs[0].GetSignBytes())
}

func (s *BroadcastSuite) TestEasyBroadcast(c *C) {
	types.Network = types.TestNetwork
	from, err := types.AccAddressFromBech32("tbnb1fds7yhw7qt9rkxw9pn65jyj004x858ny4xf2dk")
	c.Assert(err, IsNil)
	to, err := types.AccAddressFromBech32("tbnb1yxfyeda8pnlxlmx0z3cwx74w9xevspwdpzdxpj")
	c.Assert(err, IsNil)
	coins := types.Coins{types.Coin{Amount: 194765912}}

	tx := BroadcastEasy{
		From:  from,
		To:    to,
		Coins: coins,
		Memo:  "hello123",
	}

	txBytes, _ := json.Marshal([]BroadcastEasy{tx})

	bstore := store.NewInMemory()
	api := API(bstore)

	w := performRequest(api, "POST", "/broadcast/easy", bytes.NewReader(txBytes))

	block := bstore.LastBlock()

	var result TxCommitResult
	err = json.Unmarshal([]byte(w.Body.String()), &result)
	c.Assert(err, IsNil)
	c.Assert(200, Equals, w.Code)
	c.Check(block.Txs[0].Memo, Equals, "hello123")
	seqNo := block.Txs[0].Signatures[0].Sequence

	// Ensure a second transaction has different sequence numbers
	w = performRequest(api, "POST", "/broadcast/easy", bytes.NewReader(txBytes))

	block = bstore.LastBlock()

	var result2 TxCommitResult
	err = json.Unmarshal([]byte(w.Body.String()), &result2)
	c.Assert(err, IsNil)
	c.Assert(200, Equals, w.Code)
	c.Check(block.Txs[0].Signatures[0].Sequence, Equals, seqNo+1)
}

func SignTx(keyManager keys.KeyManager, txn txType.StdTx, seq int64) ([]byte, map[string]string, error) {
	signMsg := txType.StdSignMsg{
		ChainID:       "Binance-Chain-Ganges",
		Memo:          txn.Memo,
		Msgs:          txn.Msgs,
		Source:        txn.Source,
		Sequence:      seq,
		AccountNumber: 0,
	}

	rawBz, err := keyManager.Sign(signMsg)
	if nil != err {
		return nil, nil, errors.Wrap(err, "fail to sign message")
	}
	hexTx := []byte(hex.EncodeToString(rawBz))
	param := map[string]string{
		"sync": "true",
	}
	return hexTx, param, nil
}
