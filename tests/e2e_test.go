package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"primo_test_11/internal/config"
	"primo_test_11/internal/database"
	"primo_test_11/internal/handler"
	"primo_test_11/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	serviceHandler http.Handler
	dbx            *database.PGExecutor
)

func queryProductById(productId int) ([]map[string]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	statement := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", repository.PRODUCT_TABLE_NAME)
	colToVals, err := dbx.Query(ctx, statement, productId)
	if err != nil {
		return nil, err
	}

	return colToVals, nil
}

func SetupService() http.Handler {
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
	serviceHandler = SetupService()

	// run all tests
	code := m.Run()

	// teardown
	teardownTables(dbx)
	os.Exit(code)
}

func TestCreateProductE2E(t *testing.T) {

	assert.NotNil(t, serviceHandler)

	// body data for create product via request
	body := []byte(`{
		"name":"p1",
		"description":"test",
		"price":10,
		"sale_price":5
	}`)
	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	// serve http path
	serviceHandler.ServeHTTP(w, req)

	// check status code should be OK
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", w.Code)
	}

	// get response and decoded it
	res := w.Result()
	var response map[string]any
	json.NewDecoder(res.Body).Decode(&response)

	// get product id from response data
	productId := response["data"].(map[string]any)["id"]

	// query product value from database to validate data in database
	productInfoList, err := queryProductById(int(productId.(float64)))

	assert.NoError(t, err)

	// expected to have only 1
	assert.Equal(t, 1, len(productInfoList))

	// check each field of data
	productInfo := productInfoList[0]
	assert.Equal(t, "p1", productInfo["name"])
	assert.Equal(t, "test", productInfo["description"])
	assert.Equal(t, 10, int(productInfo["price"].(float64)))
	assert.Equal(t, 5, int(productInfo["sale_price"].(float64)))
}

func TestUpdateProductE2E(t *testing.T) {

	assert.NotNil(t, serviceHandler)

	// body data for update product via request
	// this for update name of product
	body := []byte(`{
		"name":"p1 update"
	}`)
	req := httptest.NewRequest(http.MethodPatch, "/product/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	// set id path to be 1 for update product id 1
	req.SetPathValue("id", "1")

	w := httptest.NewRecorder()

	serviceHandler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", w.Code)
	}

	// get response and decoded it
	res := w.Result()
	var response map[string]any
	json.NewDecoder(res.Body).Decode(&response)

	// query product value id 1 from database to validate data in database
	productInfoList, err := queryProductById(1)

	assert.NoError(t, err)

	// expected to have only 1
	assert.Equal(t, 1, len(productInfoList))

	productInfo := productInfoList[0]

	// check each field of data, for name attribute must updated
	//	otherwise should be the same of previous
	assert.Equal(t, "p1 update", productInfo["name"])
	assert.Equal(t, "test", productInfo["description"])
	assert.Equal(t, 10, int(productInfo["price"].(float64)))
	assert.Equal(t, 5, int(productInfo["sale_price"].(float64)))
}

// TODO: test case error response from update
