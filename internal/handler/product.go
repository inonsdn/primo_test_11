package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"primo_test_11/internal/repository"
	"strconv"
)

type ProductHandler struct {
	productRepo *repository.ProductRepo
}

func (p *ProductHandler) GetRouteInfo() []RoutePath {
	return []RoutePath{
		{
			Method:  http.MethodPost,
			Path:    "/product",
			Handler: p.CreateProduct,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/product/{id}",
			Handler: p.UpdateProduct,
		},
	}
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

func (p *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateProduct")
	productId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"successful": false,
			"error_code": http.StatusBadRequest,
		})
		return
	}
	var params repository.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		ResponseJSON(w, http.StatusBadRequest, map[string]any{
			"successful": false,
			"error_code": http.StatusBadRequest,
		})
		return
	}

	err = p.productRepo.UpdateProduct(productId, params)

	if err != nil {
		// response error
		return
	}

	// response
	ResponseJSON(w, http.StatusOK, map[string]any{
		"successful": true,
		"error_code": "",
	})
}
