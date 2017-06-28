package msg

import (
	"github.com/btccom/copernicus/model"
	"io"
)

const (
	HeaderSize       = 80
	AllowedTimeDrift = 2 * 60 * 60
	MaxBlockSize     = 1 * 1000 * 1000
)

type BlockMessage struct {
	Message
	Block *model.Block
	Txs   []*TxMessage
}

func (msg *BlockMessage) AddTx(tx *TxMessage) error {
	msg.Txs = append(msg.Txs, tx)
	return nil
}

func (msg *BlockMessage) ClearTxs() {
	msg.Txs = make([]*TxMessage, 0, 2048)
}

func (blockMessage *BlockMessage) BitcoinSerialize(w io.Writer, size uint32) error {
	return nil
}

func (blockMessage *BlockMessage) BitcoinParse(reader io.Reader, size uint32) error {
	return nil
}

func (blockMessage *BlockMessage) Command() string {
	return CommandBlock
}

func (blockMessage *BlockMessage) MaxPayloadLength(size uint32) uint32 {
	return 0
}
