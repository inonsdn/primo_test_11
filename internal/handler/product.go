package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"primo_test_11/internal/repository"
)

type ProductHandler struct {
	productRepo *repository.ProductRepo
}

func NewProductHandler(productRepo *repository.ProductRepo) *ProductHandler {
	return &ProductHandler{productRepo: productRepo}
}

func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateProduct")

	var params repository.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"successful": false,
			"error_code": http.StatusBadRequest,
			"data":       nil,
		})
		return
	}

	productId, err := p.productRepo.CreateProduct(params)

	if err != nil {
		// response error
		return
	}

	// response
	ResponseJSON(w, http.StatusOK, map[string]any{
		"successful": true,
		"error_code": "",
		"data": map[string]any{
			"id": productId,
		},
	})
}
