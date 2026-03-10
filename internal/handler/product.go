package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"primo_test_11/internal/repository"
	"strconv"
)

type ProductHandler struct {
	productRepo repository.ProductRepoInterface
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

func NewProductHandler(productRepo repository.ProductRepoInterface) *ProductHandler {
	return &ProductHandler{productRepo: productRepo}
}

// CreateProduct godoc
// @Summary Create product
// @Description Create a new product
// @Tags product
// @Accept json
// @Produce json
// @Param request body repository.CreateProductRequest true "Create product payload"
// @Success 200 {object} handler.CommonResponse
// @Failure 400 {object} handler.CommonResponse
// @Router /product [post]
func (p *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateProduct")

	var params repository.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		ResponseJSON(w, http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	productId, err := p.productRepo.CreateProduct(params)

	if err != nil {
		ResponseJSON(w, http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// response
	ResponseJSON(w, http.StatusOK, SuccessResponse(map[string]any{"id": productId}))
}

// UpdateProduct godoc
// @Summary Update product
// @Description Partially update a product by id
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body repository.UpdateProductRequest true "Update product payload"
// @Success 200 {object} handler.CommonResponse
// @Failure 400 {object} handler.CommonResponse
// @Router /product/{id} [patch]
func (p *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	fmt.Println("UpdateProduct", r.PathValue("id"))
	productId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		ResponseJSON(w, http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	var params repository.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		ResponseJSON(w, http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	rowAffected, err := p.productRepo.UpdateProduct(productId, params)

	if err != nil {
		ResponseJSON(w, http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if rowAffected == 0 {
		ResponseJSON(w, http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, fmt.Sprintf("Not found product id %d", productId)))
		return
	}

	// response
	ResponseJSON(w, http.StatusOK, SuccessResponse(nil))
}
