package mempool

import (
	"bytes"
	"testing"

	"github.com/btcboost/copernicus/btcutil"
	"github.com/btcboost/copernicus/core"
	"github.com/btcboost/copernicus/model"
	"github.com/btcboost/copernicus/utils"
)

func fromTx(tx *model.Tx, pool *Mempool) *TxMempoolEntry {
	var inChainValue btcutil.Amount
	if pool != nil && pool.HasNoInputsOf(tx) {
		inChainValue = btcutil.Amount(tx.GetValueOut())
	}
	entry := NewTxMempoolEntry(tx, 0, 0, 0, 1, inChainValue, false, 4, nil)
	return entry
}

func TestMempoolAddUnchecked(t *testing.T) {
	txParentPtr := model.NewTx()
	txParentPtr.Ins = make([]*model.TxIn, 1)
	txParentPtr.Ins[0] = model.NewTxIn(model.NewOutPoint(&utils.HashOne, 0), []byte{model.OP_11})
	txParentPtr.Outs = make([]*model.TxOut, 3)
	for i := 0; i < 3; i++ {
		txParentPtr.Outs[i] = model.NewTxOut(33000, []byte{model.OP_11, model.OP_EQUAL})
	}
	parentBuf := bytes.NewBuffer(nil)
	txParentPtr.Serialize(parentBuf)
	parentHash := core.DoubleSha256Hash(parentBuf.Bytes())
	var txChild [3]model.Tx
	for i := 0; i < 3; i++ {
		txChild[i].Ins = make([]*model.TxIn, 1)
		txChild[i].Ins[0] = model.NewTxIn(model.NewOutPoint(&parentHash, uint32(i)), []byte{model.OP_11})
		txChild[i].Outs = make([]*model.TxOut, 1)
		txChild[i].Outs[0] = model.NewTxOut(11000, []byte{model.OP_11, model.OP_EQUAL})
	}

	var txGrandChild [3]model.Tx
	for i := 0; i < 3; i++ {
		childBuf := bytes.NewBuffer(nil)
		txChild[i].Serialize(childBuf)
		txChildID := core.DoubleSha256Hash(childBuf.Bytes())
		txGrandChild[i].Ins = make([]*model.TxIn, 1)
		txGrandChild[i].Ins[0] = model.NewTxIn(model.NewOutPoint(&txChildID, 0), []byte{model.OP_11})
		txGrandChild[i].Outs = make([]*model.TxOut, 1)
		txGrandChild[i].Outs[0] = model.NewTxOut(11000, []byte{model.OP_11, model.OP_EQUAL})
	}
	testPool := NewMemPool(utils.FeeRate{0})

	//Nothing in pool, remove should do nothing:
	poolSize := testPool.Size()
	testPool.RemoveRecursive(txParentPtr, UNKNOWN)
	if testPool.Size() != poolSize {
		t.Errorf("current poolSize : %d, except the mempoolSize : %d\n",
			testPool.Size(), poolSize)
	}

	/*
		//Just the parent:
		testPool.AddUnchecked(&txParentPtr.Hash, fromTx(txParentPtr, nil), true)
		poolSize = testPool.Size()
		fmt.Println("---------- 6 ----------- poolSize : ", poolSize)
		testPool.RemoveRecursive(txParentPtr, UNKNOWN)
		fmt.Println("----------  7-----------")
		if testPool.Size() != poolSize-1 {
			t.Errorf("current poolSize : %d, except the mempoolSize : %d\n",
				testPool.Size(), poolSize-1)
		}
	*/
}