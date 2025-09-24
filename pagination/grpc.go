package pagination

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// FromGRPCRequest converts gRPC pagination request to domain request
func FromGRPCRequest(page, limit int32) Request {
	return Request{
		Page:  page,
		Limit: limit,
	}
}

// ValidateGRPCRequest validates gRPC pagination request
func ValidateGRPCRequest(page, limit int32) error {
	calc := NewCalculator()
	req := FromGRPCRequest(page, limit)

	if err := calc.ValidateRequest(req); err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid pagination: %v", err)
	}

	return nil
}
