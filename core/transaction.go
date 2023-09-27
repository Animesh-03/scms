package core

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
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
	Signature []byte            `json:"signature"`
}

func (t *Transaction) Bytes() []byte {
	return bytes.Join([][]byte{
		[]byte(t.Sender),
		[]byte(t.Receiver),
		[]byte(t.ProductID),
		ToByte(int64(t.Status)),
	}, []byte{})
}

func (t *Transaction) Stringify() string {
	txJson, _ := json.MarshalIndent(t, "", "	")

	return string(txJson) + "\n"
}

func (t *Transaction) Hash() []byte {
	data := t.Bytes()

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

func (t *Transaction) Verify(pubKey ecdsa.PublicKey) bool {
	if len(t.ID) == 0 || !bytes.Equal(t.Hash(), t.ID) {
		return false
	}

	if v := ecdsa.VerifyASN1(&pubKey, t.Bytes(), t.Signature); !v {
		return false
	}

	return true
}
