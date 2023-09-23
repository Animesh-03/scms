package core

import (
	"bytes"
	"crypto/sha256"
)

type TransactionStatus uint16

const (
	Manufactured TransactionStatus = 1
	Dispatched   TransactionStatus = 2
	Received     TransactionStatus = 3
)

type Transaction struct {
	ID        []byte            `json:"id"`
	Sender    string            `json:"sender"`
	Receiver  string            `json:"receiver"`
	ProductID string            `json:"productid"`
	Status    TransactionStatus `json:"status"`
}

func (t *Transaction) Hash() []byte {
	var data []byte
	data = append(data, []byte(t.Sender)...)
	data = append(data, []byte(t.Sender)...)
	data = append(data, []byte(t.ProductID)...)
	data = append(data, ToByte(int64(t.Status))...)

	hash := sha256.Sum256(data)
	return hash[:]
}

func NewTransaction(sender, receiver, productId string, status TransactionStatus) *Transaction {
	transaction := &Transaction{
		Sender:    sender,
		Receiver:  receiver,
		ProductID: productId,
		Status:    status,
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
