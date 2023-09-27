package node

import (
	"encoding/base64"

	"github.com/Animesh-03/scms/core"
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
	status, err := node.GetStatusOfProduct(productStatus.ProductId)

	statusString := "Product Not Found"
	if err != nil {
		switch status {
		case core.Manufactured:
			statusString = "Manufactured"
		case core.Dispatched:
			statusString = "Dispatched"
		case core.Received:
			statusString = "Delivered"
		}
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
