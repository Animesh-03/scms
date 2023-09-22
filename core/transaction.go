package core

import (
	"bytes"
	"crypto/sha256"
)

type Transaction struct {
	ID        []byte `json:"id"`
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	ProductID string `json:"productid"`
}

func (t *Transaction) Hash() []byte {
	var data []byte
	data = append(data, []byte(t.Sender)...)
	data = append(data, []byte(t.Sender)...)
	data = append(data, []byte(t.ProductID)...)

	hash := sha256.Sum256(data)
	return hash[:]
}

func NewTransaction(sender, receiver, productId string) *Transaction {
	transaction := &Transaction{
		Sender:    sender,
		Receiver:  receiver,
		ProductID: productId,
	}

	transaction.ID = transaction.Hash()

	return transaction
}

func (t *Transaction) Verify() bool {
	if len(t.ID) == 0 || !bytes.Equal(t.Hash(), t.ID) {
		return false
	}

	return true
}
