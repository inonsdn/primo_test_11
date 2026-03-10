package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) QueryRow(ctx context.Context, dest []any, statement string, args ...any) error {
	called := m.Called(ctx, dest, statement, args)

	err := called.Error(0)
	if err != nil {
		return err
	}
	// has dest value
	if len(dest) > 0 {
		// assign value from mock Return to pointer of dest
		if ptr, ok := dest[0].(*int); ok {
			*ptr = called.Int(1)
		}
	}
	return called.Error(0)
}

func (m *MockDB) Execute(ctx context.Context, statement string, args ...any) (int64, error) {
	called := m.Called(ctx, statement, args)
	return int64(called.Int(0)), called.Error(1)
}

func (m *MockDB) Query(ctx context.Context, statement string, args ...any) ([]map[string]any, error) {
	called := m.Called(ctx, statement, args)
	return called.Get(0).([]map[string]any), called.Error(1)
}

func (m *MockDB) Connect() error {
	called := m.Called()
	return called.Error(0)
}

func (m *MockDB) Disconnect() {
	called := m.Called()
	called.Error(0)
}

func TestCreateProduct(t *testing.T) {
	mockDB := new(MockDB)
	repo := NewProductRepo(mockDB)

	name := "product1"
	description := "product of living"
	salePrice := float64(10.0)
	price := float64(5.0)

	createProductRequest := CreateProductRequest{
		Name:        name,
		Description: &description,
		SalePrice:   &salePrice,
		Price:       price,
	}

	mockDB.On(
		"QueryRow",
		mock.Anything, // ctx
		mock.Anything, // dest
		mock.Anything, // query
		mock.Anything, // args
	).Return(nil, 1)

	productId, err := repo.CreateProduct(createProductRequest)
	assert.NoError(t, err)
	assert.Equal(t, 1, productId)
	mockDB.AssertExpectations(t)
}

func TestCreateProduct_WithErrorQueryRow(t *testing.T) {
	mockDB := new(MockDB)
	repo := NewProductRepo(mockDB)

	name := "product1"
	description := "product of living"
	salePrice := float64(10.0)
	price := float64(5.0)

	createProductRequest := CreateProductRequest{
		Name:        name,
		Description: &description,
		SalePrice:   &salePrice,
		Price:       price,
	}

	mockDB.On(
		"QueryRow",
		mock.Anything, // ctx
		mock.Anything, // dest
		mock.Anything, // query
		mock.Anything, // args
	).Return(fmt.Errorf("Error execute QueryRow"), nil)

	productId, err := repo.CreateProduct(createProductRequest)
	assert.Error(t, err, "Error execute QueryRow")
	assert.Equal(t, 0, productId)
	mockDB.AssertExpectations(t)
}

func TestUpdateProduct(t *testing.T) {
	mockDB := new(MockDB)
	repo := NewProductRepo(mockDB)

	name := "product1"
	description := "product of living"
	salePrice := float64(10.0)
	price := float64(5.0)

	updateProductRequest := UpdateProductRequest{
		Name:        &name,
		Description: &description,
		SalePrice:   &salePrice,
		Price:       &price,
	}

	mockDB.On(
		"Execute",
		mock.Anything, // ctx
		mock.Anything, // query
		mock.Anything, // args
	).Return(1, nil)

	rowAffected, err := repo.UpdateProduct(1, updateProductRequest)
	assert.Equal(t, 1, rowAffected)
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestUpdateProduct_WithErrorExecute(t *testing.T) {
	mockDB := new(MockDB)
	repo := NewProductRepo(mockDB)

	name := "product1"
	description := "product of living"
	salePrice := float64(10.0)
	price := float64(5.0)

	updateProductRequest := UpdateProductRequest{
		Name:        &name,
		Description: &description,
		SalePrice:   &salePrice,
		Price:       &price,
	}

	mockDB.On(
		"Execute",
		mock.Anything, // ctx
		mock.Anything, // query
		mock.Anything, // args
	).Return(0, fmt.Errorf("Error execute Execute"))

	rowAffected, err := repo.UpdateProduct(1, updateProductRequest)
	assert.Equal(t, 0, rowAffected)
	assert.Error(t, err, "Error execute Execute")
	mockDB.AssertExpectations(t)
}

func TestUpdateProduct_NoUpdate(t *testing.T) {
	mockDB := new(MockDB)
	repo := NewProductRepo(mockDB)

	updateProductRequest := UpdateProductRequest{}

	mockDB.On(
		"Execute",
		mock.Anything, // ctx
		mock.Anything, // query
		mock.Anything, // args
	).Return(0, fmt.Errorf("Error execute Execute"))

	rowAffected, err := repo.UpdateProduct(1, updateProductRequest)
	assert.Equal(t, 0, rowAffected)
	assert.NoError(t, err)
	mockDB.AssertNotCalled(t, "Execute")
}
