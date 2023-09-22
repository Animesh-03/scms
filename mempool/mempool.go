package mempool

import (
	"encoding/hex"

	"github.com/Animesh-03/scms/core"
)

type MemPool struct {
	Pool map[string]*core.Transaction
}

func (mp *MemPool) AddToPool(tx *core.Transaction) {
	mp.Pool[hex.EncodeToString(tx.ID)] = tx
}

func (mp *MemPool) AddAllToPool(txs []*core.Transaction) {
	for _, tx := range txs {
		mp.Pool[hex.EncodeToString(tx.ID)] = tx
	}
}

func (mp *MemPool) GetTransactions(count int) (txs []*core.Transaction) {
	i := 0
	for _, tx := range mp.Pool {
		txs = append(txs, tx)
		i++
		if i == count {
			break
		}
	}

	return txs
}

func (mp *MemPool) Remove(tx *core.Transaction) {
	delete(mp.Pool, hex.EncodeToString(tx.ID))
}

func (mp *MemPool) RemoveAll(txs []*core.Transaction) {
	for _, tx := range txs {
		delete(mp.Pool, hex.EncodeToString(tx.ID))
	}
}
