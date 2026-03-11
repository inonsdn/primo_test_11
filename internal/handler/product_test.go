package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"primo_test_11/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepoInterface struct {
	mock.Mock
}

func (m *MockProductRepoInterface) CreateProduct(c repository.CreateProductRequest) (int, error) {
	called := m.Called(c)
	return called.Int(0), called.Error(1)
}
func (m *MockProductRepoInterface) UpdateProduct(productId int, u repository.UpdateProductRequest) (int, error) {
	called := m.Called(productId, u)
	return called.Int(0), called.Error(1)
}

func TestCreateProduct(t *testing.T) {
	mockProductRepo := new(MockProductRepoInterface)
	productHandler := NewProductHandler(mockProductRepo)

	body := map[string]any{
		"name":        "shoes",
		"description": "nice shoes",
		"price":       19.99,
	}
	bodyBytes, _ := json.Marshal(body)

	mockProductRepo.On("CreateProduct", mock.Anything).Return(1, nil)

	// create request for test
	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// call test function
	productHandler.CreateProduct(w, req)

	// get result from response
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]any
	json.NewDecoder(res.Body).Decode(&response)
	assert.Equal(t, true, response["successful"])
	assert.Equal(t, float64(1), response["data"].(map[string]any)["id"])
}

func TestUpdateProduct(t *testing.T) {
	mockProductRepo := new(MockProductRepoInterface)
	productHandler := NewProductHandler(mockProductRepo)

	body := map[string]any{
		"name": "shoes",
	}
	bodyBytes, _ := json.Marshal(body)

	mockProductRepo.On("UpdateProduct", mock.Anything, mock.Anything).Return(1, nil)

	// create request for test
	req := httptest.NewRequest(http.MethodPatch, "/product/1", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "1")
	w := httptest.NewRecorder()

	// call test function
	productHandler.UpdateProduct(w, req)

	// get result from response
	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	var response map[string]any
	json.NewDecoder(res.Body).Decode(&response)
	fmt.Println(response)
	assert.Equal(t, true, response["successful"])
}

// TODO: test case error response from update
