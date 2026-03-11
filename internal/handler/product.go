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

	// TODO: verify field of attribute that must be sent which is name and price
	//	and validate type of data

	// create product to database
	productId, err := p.productRepo.CreateProduct(params)

	// response error if got an error
	if err != nil {
		ResponseJSON(w, http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// response successful
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

	// get id from path value
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

	// update product to database
	rowAffected, err := p.productRepo.UpdateProduct(productId, params)

	// response error if got an error
	if err != nil {
		ResponseJSON(w, http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// in case no row effect, means no product id of given
	if rowAffected == 0 {
		ResponseJSON(w, http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, fmt.Sprintf("Not found product id %d", productId)))
		return
	}

	// response
	ResponseJSON(w, http.StatusOK, SuccessResponse(nil))
}
