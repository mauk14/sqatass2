package main

import (
	"log"
	"messanger/pkg/postgres"
	httpDelivery "messanger/services/receiptManage/internal/Delivery/http"
	"messanger/services/receiptManage/internal/Repository"
	"messanger/services/receiptManage/internal/Use_Case"
	"os"
)

func main() {
	db, err := postgres.OpenDb(os.Getenv("RECEIPTDB_URI"))
	if err != nil {
		log.Fatal(err)
		return
	}
	app := httpDelivery.NewApp(Use_Case.NewReceiptUseCase(Repository.NewReceiptRepository(db)))

	err = app.Route().Run(":4000")
	if err != nil {
		log.Fatal(err)
		return
	}
}
