package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"primo_test_11/internal/config"
	"primo_test_11/internal/database"
	"primo_test_11/internal/handler"
	"primo_test_11/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func SetupRouter() http.Handler {
	mux := http.NewServeMux()
	// load config of database
	dbConfig := config.LoadDatabaseConfig()
	// construct db handler with executor
	mockExecutor := database.NewPGExecutor(dbConfig)
	err := mockExecutor.Connect()
	if err != nil {
		fmt.Println("Error when connect to database", err)
		return nil
	}

	// construct repo and handler of repo
	productRepo := repository.NewProductRepo(mockExecutor)
	productHandler := handler.NewProductHandler(productRepo)
	routePaths := productHandler.GetRouteInfo()
	for _, routePath := range routePaths {
		http.Handle(routePath.Path, handler.MakeHandler(routePath))
		mux.HandleFunc(routePath.Path, handler.MakeHandler(routePath))
	}

	return mux
}
func TestCreateProductE2E(t *testing.T) {

	router := SetupRouter()

	assert.NotNil(t, router)

	body := []byte(`{
		"name":"p1",
		"description":"test",
		"price":10,
		"sale_price":5
	}`)

	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", w.Code)
	}
}
