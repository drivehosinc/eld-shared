package pagination

// Request represents pagination parameters from client
type Request struct {
	Page  int32 `json:"page"`
	Limit int32 `json:"limit"`
}

type Response struct {
	Page       int32 `json:"page"`
	Limit      int32 `json:"limit"`
	Offset     int32 `json:"offset"`
	Total      int64 `json:"total"`
	TotalPages int32 `json:"total_pages"`
}

type Result[T any] struct {
	Data       []T      `json:"data"`
	Pagination Response `json:"pagination"`
}

// Default pagination values
const (
	DefaultLimit = 10
	MaxLimit     = 100
	MinPage      = 1
)
