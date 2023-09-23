package node

import "github.com/gin-gonic/gin"

type SendTransactionData struct {
	Reciever  string `json:"receiver"`
	ProductId string `json:"productid"`
}

func SendTransaction(c *gin.Context, node *Node) {
	var transactionData SendTransactionData
	c.BindJSON(&transactionData)
	transaction, err := node.MakeTransaction(transactionData.Reciever, transactionData.ProductId)
	if err != nil {
		c.IndentedJSON(500, gin.H{
			"error": err.Error(),
		})
	}

	c.IndentedJSON(200, transaction)
}

func GetNodeInfo(c *gin.Context, node *Node) {
	c.IndentedJSON(200, node)
}
