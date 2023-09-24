package httpDelivery

import (
	"github.com/gin-gonic/gin"
	"messanger/services/receiptManage/internal/Use_Case"
)

type App struct {
	receiptUseCase Use_Case.ReceiptUseCase
	router         *gin.Engine
}

func setUpRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}

func NewApp(useCase Use_Case.ReceiptUseCase) *App {
	return &App{
		receiptUseCase: useCase,
		router:         setUpRouter(),
	}
}

func (a *App) Route() *gin.Engine {
	a.router.POST("/receipts/create", a.createReceipt)
	a.router.DELETE("/receipts/delete/:id", a.deleteReceipt)
	a.router.GET("/receipts/get/:id", a.getReceipt)
	a.router.GET("/receipts/get", a.getAllReceipt)
	a.router.PATCH("/receipts/update/:id", a.updateReceipt)
	return a.router
}
