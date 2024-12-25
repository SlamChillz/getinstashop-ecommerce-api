package validators

import (
	"errors"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
)

// ValidateName checks if the Name is non-empty and within length constraints
func ValidateName(name string) string {
	if name == "" {
		return "name cannot be empty"
	}
	if len(name) < 3 || len(name) > 100 {
		return "name must be between 3 and 100 characters"
	}
	return ""
}

// ValidateDescription checks if the Description is non-empty and within length constraints
func ValidateDescription(description string) string {
	if description == "" {
		return "description cannot be empty"
	}
	if len(description) < 10 || len(description) > 500 {
		return "description must be between 10 and 500 characters"
	}
	return ""
}

// ValidatePrice checks if the Price is greater than 0
func ValidatePrice(price float64) string {
	if price <= 0 {
		return "price must be greater than 0"
	}
	return ""
}

// ValidateStock checks if the Stock is a non-negative integer
func ValidateStock(stock int) string {
	if stock < 0 {
		return "stock cannot be negative"
	}
	return ""
}

// ValidateProduct validates the CreateProductInput struct
func ValidateProduct(product types.CreateProductInput) (types.CreateProductErrMessage, error) {
	errMessage := types.CreateProductErrMessage{
		Name:        ValidateName(product.Name),
		Description: ValidateDescription(product.Description),
		Price:       ValidatePrice(product.Price),
		Stock:       ValidateStock(product.Stock),
	}
	if errMessage.Name == "" && errMessage.Description == "" && errMessage.Price == "" && errMessage.Stock == "" {
		return errMessage, nil
	}
	return errMessage, errors.New("invalid product input")
}
