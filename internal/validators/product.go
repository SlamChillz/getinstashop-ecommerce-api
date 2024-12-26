package validators

import (
	"errors"
	"github.com/slamchillz/getinstashop-ecommerce-api/internal/types"
)

// ValidateName checks if the Name is non-empty and within length constraints
func ValidateName(name string) string {
	var msg string
	if name == "" {
		msg = "name cannot be empty"
	}
	if len(name) < 3 || len(name) > 100 {
		msg = "name must be between 3 and 100 characters"
	}
	return msg
}

// ValidateDescription checks if the Description is non-empty and within length constraints
func ValidateDescription(description string) string {
	var msg string
	if description == "" {
		msg = "description cannot be empty"
	}
	if len(description) < 3 || len(description) > 500 {
		msg = "description must be between 3 and 500 characters"
	}
	return msg
}

// ValidatePrice checks if the Price is greater than 0
func ValidatePrice(price float64) string {
	var msg string
	if price <= 0 {
		msg = "price must be greater than 0"
	}
	return msg
}

// ValidateStock checks if the Stock is a non-negative integer
func ValidateStock(stock int) string {
	var msg string
	if stock < 0 {
		msg = "stock cannot be negative"
	}
	return msg
}

// ValidateProduct validates the CreateProductInput struct
func ValidateProduct(product types.CreateProductInput) (types.ProductErrMessage, error) {
	errMessage := types.ProductErrMessage{
		Name:        ValidateName(product.Name),
		Description: ValidateDescription(product.Description),
		Price:       ValidatePrice(product.Price),
		Stock:       ValidateStock(product.Stock),
	}
	if errMessage.Name == "" && errMessage.Description == "" && errMessage.Price == "" && errMessage.Stock == "" {
		return errMessage, nil
	}
	return errMessage, errors.New("invalid create product input")
}

func ValidateProductUpdateInput(product types.ProductUpdateInput) (types.ProductErrMessage, error) {
	var errMessage types.ProductErrMessage
	if product.Name != nil {
		if msg := ValidateName(*product.Name); msg != "" {
			errMessage.Name = msg
		}
	}
	if product.Description != nil {
		if msg := ValidateDescription(*product.Description); msg != "" {
			errMessage.Description = msg
		}
	}
	if product.Price != nil {
		if msg := ValidatePrice(*product.Price); msg != "" {
			errMessage.Price = msg
		}
	}
	if product.Stock != nil {
		stock := *product.Stock
		stock32 := int(stock)
		if msg := ValidateStock(stock32); msg != "" {
			errMessage.Stock = msg
		}
	}
	if errMessage.Name == "" && errMessage.Description == "" && errMessage.Price == "" && errMessage.Stock == "" {
		return errMessage, nil
	}
	return errMessage, errors.New("invalid product input")
}
