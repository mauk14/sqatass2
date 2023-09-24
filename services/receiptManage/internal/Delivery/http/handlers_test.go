package httpDelivery

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"messanger/pkg/postgres"
	"messanger/services/receiptManage/internal/Domain"
	"messanger/services/receiptManage/internal/Repository"
	"messanger/services/receiptManage/internal/Use_Case"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCreateReceipt(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	db, err := postgres.OpenDb(os.Getenv("RECEIPTDB_URI"))
	if err != nil {
		t.Fatal(err)
	}
	app := &App{receiptUseCase: Use_Case.NewReceiptUseCase(Repository.NewReceiptRepository(db))}
	router.POST("/create-receipt", app.createReceipt)

	type Cases struct {
		name         string
		receipt      *Domain.Receipt
		expectedCode int
	}

	cases := []Cases{
		{
			name: "All good",
			receipt: &Domain.Receipt{
				Title:       "Kruto",
				Author:      "Chelovek",
				Description: "Delai chto hochesh"},
			expectedCode: 201,
		},
		{
			name:         "receipt is clean",
			receipt:      &Domain.Receipt{},
			expectedCode: 400,
		},
		{
			name:         "title and author clean",
			receipt:      &Domain.Receipt{Description: "kuku"},
			expectedCode: 400},
	}

	for _, c := range cases {
		requestBody, err := json.Marshal(c.receipt)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v, in test %s", err, c.name)
		}
		req, err := http.NewRequest("POST", "/create-receipt", bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v, in test %s ", err, c.name)
		}
		req.Header.Set("Content-Type", "application/json")
		respRecorder := httptest.NewRecorder()
		router.ServeHTTP(respRecorder, req)
		if status := respRecorder.Code; status != c.expectedCode {
			t.Errorf("Expected status code %d, but got %d, in case %s", c.expectedCode, status, c.name)
		}
	}

}

func TestDeleteReceipt(t *testing.T) {
	router := gin.New()
	db, err := postgres.OpenDb(os.Getenv("RECEIPTDB_URI"))
	if err != nil {
		t.Fatal(err)
	}
	app := &App{receiptUseCase: Use_Case.NewReceiptUseCase(Repository.NewReceiptRepository(db))}
	router.DELETE("/delete-receipt/:id", app.deleteReceipt)

	type Cases struct {
		name         string
		id           string
		expectedCode int
	}
	cases := []Cases{
		{
			name:         "All good",
			id:           "22",
			expectedCode: 200,
		},

		{
			name:         "id less than 1",
			id:           "-1",
			expectedCode: 400,
		},

		{
			name:         "more than need",
			id:           "500",
			expectedCode: 400,
		},
		{
			name:         "id with symbols",
			id:           "wda",
			expectedCode: 400,
		},
	}

	for _, c := range cases {
		req, err := http.NewRequest("DELETE", "/delete-receipt/"+(c.id), nil)
		if err != nil {
			t.Fatalf("Failed to delete request: %v, in test %s ", err, c.name)
		}
		respRecorder := httptest.NewRecorder()
		router.ServeHTTP(respRecorder, req)
		if status := respRecorder.Code; status != c.expectedCode {
			t.Errorf("Expected status code %d, but got %d, in case %s", c.expectedCode, status, c.name)
		}
	}

}

func TestGetReceipt(t *testing.T) {
	router := gin.New()
	db, err := postgres.OpenDb(os.Getenv("RECEIPTDB_URI"))
	if err != nil {
		t.Fatal(err)
	}
	app := &App{receiptUseCase: Use_Case.NewReceiptUseCase(Repository.NewReceiptRepository(db))}
	router.GET("/get-receipt/:id", app.getReceipt)

	type Cases struct {
		name         string
		id           string
		expectedCode int
	}
	cases := []Cases{
		{
			name:         "All good",
			id:           "2",
			expectedCode: 200,
		},

		{
			name:         "id less than 1",
			id:           "-1",
			expectedCode: 400,
		},

		{
			name:         "more than need",
			id:           "500",
			expectedCode: 400,
		},
		{
			name:         "id with symbols",
			id:           "wda",
			expectedCode: 400,
		},
	}

	for _, c := range cases {
		req, err := http.NewRequest("GET", "/get-receipt/"+(c.id), nil)
		if err != nil {
			t.Fatalf("Failed to get request: %v, in test %s ", err, c.name)
		}
		respRecorder := httptest.NewRecorder()
		router.ServeHTTP(respRecorder, req)
		if status := respRecorder.Code; status != c.expectedCode {
			t.Errorf("Expected status code %d, but got %d, in case %s", c.expectedCode, status, c.name)
		}
	}

}

func TestGetAllReceipt(t *testing.T) {
	router := gin.New()
	db, err := postgres.OpenDb(os.Getenv("RECEIPTDB_URI"))
	if err != nil {
		t.Fatal(err)
	}
	app := &App{receiptUseCase: Use_Case.NewReceiptUseCase(Repository.NewReceiptRepository(db))}
	router.GET("/get-receipt", app.getAllReceipt)

	req, err := http.NewRequest("GET", "/get-receipt", nil)
	if err != nil {
		t.Fatalf("Failed to get request: %v", err)
	}
	respRecorder := httptest.NewRecorder()
	router.ServeHTTP(respRecorder, req)
	if status := respRecorder.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, status)
	}

}

func TestUpdateReceipt(t *testing.T) {
	router := gin.New()
	db, err := postgres.OpenDb(os.Getenv("RECEIPTDB_URI"))
	if err != nil {
		t.Fatal(err)
	}
	app := &App{receiptUseCase: Use_Case.NewReceiptUseCase(Repository.NewReceiptRepository(db))}
	router.PATCH("/update-receipt/:id", app.updateReceipt)

	type Cases struct {
		name         string
		id           string
		receipt      *Domain.Receipt
		expectedCode int
	}

	cases := []Cases{
		{
			name: "All good",
			id:   "5",
			receipt: &Domain.Receipt{
				Title: "Cinema",
			},
			expectedCode: 200,
		},
		{
			name: "id less than 1",
			id:   "-1",
			receipt: &Domain.Receipt{
				Title: "Cinema",
			},
			expectedCode: 400,
		},
		{
			name: "id more than need",
			id:   "500",
			receipt: &Domain.Receipt{
				Title: "Cinema",
			},
			expectedCode: 400,
		},

		{
			name: "bad id",
			id:   "wadsadwa",
			receipt: &Domain.Receipt{
				Title: "Cinema",
			},
			expectedCode: 400,
		},
	}

	for _, c := range cases {
		requestBody, err := json.Marshal(c.receipt)
		if err != nil {
			t.Fatalf("Failed to marshal JSON: %v, in test %s", err, c.name)
		}
		req, err := http.NewRequest("PATCH", "/update-receipt/"+(c.id), bytes.NewBuffer(requestBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v, in test %s ", err, c.name)
		}
		req.Header.Set("Content-Type", "application/json")
		respRecorder := httptest.NewRecorder()
		router.ServeHTTP(respRecorder, req)
		if status := respRecorder.Code; status != c.expectedCode {
			t.Errorf("Expected status code %d, but got %d, in case %s", c.expectedCode, status, c.name)
		}
	}

}
