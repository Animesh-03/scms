package node

import (
	"encoding/base64"

	"github.com/Animesh-03/scms/core"
	"github.com/Animesh-03/scms/logger"
	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

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
		return
	}

	c.IndentedJSON(200, transaction)
}

func GetNodeInfo(c *gin.Context, node *Node) {
	c.IndentedJSON(200, node)
}

type ProductStatusData struct {
	ProductId string `json:"productid"`
}

func GetProductStatus(c *gin.Context, node *Node) {
	var productStatus ProductStatusData
	c.BindJSON(&productStatus)
	status, _ := node.GetStatusOfProduct(productStatus.ProductId)

	statusString := ""
	switch status {
	case 0:
		statusString = "Product Not Manufactured"
	case core.Manufactured:
		statusString = "Manufactured"
	case core.Dispatched:
		statusString = "Dispatched"
	case core.Received:
		statusString = "Delivered"
	}

	img, err := qrcode.Encode(statusString, qrcode.Medium, 512)
	if err != nil {
		c.IndentedJSON(500, gin.H{
			"error": "error generating QR code",
		})
	}
	imgBytes := base64.StdEncoding.EncodeToString(img)

	c.HTML(200, "qrcode.html", gin.H{
		"image": imgBytes,
	})
}

func Dispute(c *gin.Context, node *Node) {
	var productStatus ProductStatusData
	c.BindJSON(&productStatus)
	status, _ := node.GetStatusOfProduct(productStatus.ProductId)

	if status == 0 {
		// Product not yet made or sent to distributor
		c.IndentedJSON(200, gin.H{
			"error": "product not found",
		})
	} else if status == core.Received {
		// Consumer is wrong, product is dispatched
		node.Network.Broadcast("dispute", []byte(node.ID))
		logger.LogInfo("Consumer is wrong, stake is being deducted\n")
		c.IndentedJSON(200, gin.H{
			"error": "Consumer is wrong, stake is being deducted",
		})
	} else {
		// Distributor is wrong
		tx, _ := node.GetTransactionOfProduct(productStatus.ProductId)
		logger.LogInfo("Distributor is wrong, stake is being deducted\n")
		node.Network.Broadcast("dispute", []byte(tx.Sender))
		c.IndentedJSON(200, gin.H{
			"error": "Distributor is wrong, stake is being deducted",
		})
	}
}
