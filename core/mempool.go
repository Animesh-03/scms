package core

import (
	"encoding/hex"
)

type MemPool struct {
	Pool map[string]*Transaction
}

func (mp *MemPool) AddToPool(tx *Transaction) {
	mp.Pool[hex.EncodeToString(tx.ID)] = tx
}

func (mp *MemPool) AddAllToPool(txs []*Transaction) {
	for _, tx := range txs {
		mp.Pool[hex.EncodeToString(tx.ID)] = tx
	}
}

func (mp *MemPool) GetTransactions(count int) (txs []*Transaction) {
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

func (mp *MemPool) Remove(tx *Transaction) {
	delete(mp.Pool, hex.EncodeToString(tx.ID))
}

func (mp *MemPool) RemoveAll(txs []*Transaction) {
	for _, tx := range txs {
		delete(mp.Pool, hex.EncodeToString(tx.ID))
	}
}
