package pagination

import (
	"fmt"
)

// Calculator handles pagination calculations
type Calculator struct {
	defaultLimit int32
	maxLimit     int32
}

// NewCalculator creates a new pagination calculator
func NewCalculator() *Calculator {
	return &Calculator{
		defaultLimit: DefaultLimit,
		maxLimit:     MaxLimit,
	}
}

// NewCalculatorWithLimits creates a calculator with custom limits
func NewCalculatorWithLimits(defaultLimit, maxLimit int32) *Calculator {
	return &Calculator{
		defaultLimit: defaultLimit,
		maxLimit:     maxLimit,
	}
}

// Calculate processes pagination request and returns normalized values
func (c *Calculator) Calculate(req Request, total int64) Response {
	// Normalize page
	page := req.Page
	if page <= 0 {
		page = MinPage
	}

	// Normalize limit
	limit := req.Limit
	if limit <= 0 {
		limit = c.defaultLimit
	}
	if limit > c.maxLimit {
		limit = c.maxLimit
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Calculate total pages
	totalPages := int32(0)
	if total > 0 {
		totalPages = (int32(total) + limit - 1) / limit
	}

	return Response{
		Page:       page,
		Limit:      limit,
		Offset:     offset,
		Total:      total,
		TotalPages: totalPages,
	}
}

// ValidateRequest validates pagination parameters
func (c *Calculator) ValidateRequest(req Request) error {
	if req.Page < 0 {
		return fmt.Errorf("page cannot be negative")
	}

	if req.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}

	if req.Limit > c.maxLimit {
		return fmt.Errorf("limit cannot exceed %d", c.maxLimit)
	}

	return nil
}

// CreateResult creates a paginated result
func CreateResult[T any](data []T, pagination Response) *Result[T] {
	return &Result[T]{
		Data:       data,
		Pagination: pagination,
	}
}
