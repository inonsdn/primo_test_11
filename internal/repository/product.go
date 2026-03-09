package repository

import (
	"context"
	"fmt"
	"log/slog"
	"primo_test_11/internal/database"
	"strings"
	"time"
)

const (
	PRODUCT_TABLE_NAME = "product"
)

type ProductRepo struct {
	dbHandler *database.DBHandler
}

func NewProductRepo(dbHandler *database.DBHandler) *ProductRepo {
	return &ProductRepo{
		dbHandler: dbHandler,
	}
}

type CreateProductRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price"`
}

type CreateProductResponse struct {
	Id int `json:"id"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       *float64 `json:"price"`
}

type UpdateProductResponse struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	UserId    int    `json:"user_id"`
	NoteGroup int    `json:"note_group"`
}

func (p *ProductRepo) CreateProduct(c CreateProductRequest) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	colNames, values := sqlStructExtraction(c)
	colStr := strings.Join(colNames, ",")
	placeholdersStr := buildPlaceholders(len(values), 0)

	statement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", PRODUCT_TABLE_NAME, colStr, placeholdersStr)
	var productId int
	err := p.dbHandler.Dbx.QueryRow(ctx, []any{&productId}, statement, values...)
	if err != nil {
		fmt.Println(err.Error())
		return productId, err
	}
	return productId, nil
}

func (p *ProductRepo) UpdateProduct(productId int, u UpdateProductRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	colNames, values := sqlStructExtraction(u)
	// no value update
	if len(colNames) == 0 {
		return nil
	}
	colStr := strings.Join(colNames, ",")
	placeholdersStr := buildPlaceholders(len(values), 1)

	args := []any{}
	args = append(args, productId)
	args = append(args, values...)

	updateColStr := colStr
	updatePlaceholdersStr := placeholdersStr
	if len(colNames) > 1 {
		updateColStr = fmt.Sprintf("( %s )", updateColStr)
		updatePlaceholdersStr = fmt.Sprintf("( %s )", updatePlaceholdersStr)
	}

	statement := fmt.Sprintf("UPDATE %s SET %s = %s WHERE id = $1", PRODUCT_TABLE_NAME, updateColStr, updatePlaceholdersStr)

	rowAffect, err := p.dbHandler.Dbx.Execute(ctx, statement, args...)
	if err != nil {
		fmt.Println(err.Error())
		slog.Error(err.Error())
		return err
	}
	fmt.Println("Row affect", rowAffect)
	return nil
}
