package tests

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"primo_test_11/internal/config"
	"primo_test_11/internal/database"
	"primo_test_11/internal/handler"
	"primo_test_11/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	router http.Handler         // shared router
	dbx    *database.PGExecutor // shared executor
)

func SetupRouter() http.Handler {
	mux := http.NewServeMux()
	// load config of database
	dbConfig := config.LoadDatabaseConfig()
	// construct db handler with executor
	dbx = database.NewPGExecutor(dbConfig)
	err := dbx.Connect()
	if err != nil {
		fmt.Println("Error when connect to database", err)
		return nil
	}

	// construct repo and handler of repo
	productRepo := repository.NewProductRepo(dbx)
	productHandler := handler.NewProductHandler(productRepo)
	routePaths := productHandler.GetRouteInfo()
	for _, routePath := range routePaths {
		http.Handle(routePath.Path, handler.MakeHandler(routePath))
		mux.HandleFunc(routePath.Path, handler.MakeHandler(routePath))
	}

	return mux
}

func teardownTables(executor database.DBExecutor) {
	ctx := context.Background()
	executor.Execute(ctx, `DROP TABLE IF EXISTS product`)
}

func TestMain(m *testing.M) {
	// setup db
	router = SetupRouter()

	// run all tests
	code := m.Run()

	// teardown
	teardownTables(dbx)
	os.Exit(code)
}

func TestCreateProductE2E(t *testing.T) {

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

func TestUpdateProductE2E(t *testing.T) {

	assert.NotNil(t, router)

	body := []byte(`{
		"name":"p1",
		"description":"test",
		"price":10,
		"sale_price":5
	}`)

	req := httptest.NewRequest(http.MethodPatch, "/product/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", w.Code)
	}
}
