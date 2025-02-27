package model

// Response represents a standard API response
type Response struct {
	Data interface{} `json:"data,omitempty"`
	Meta *Metadata   `json:"meta,omitempty"`
}

// Metadata contains metadata for paginated responses
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// NewResponse creates a new API response
func NewResponse(data interface{}) Response {
	return Response{
		Data: data,
	}
}

// NewPageResponse creates a new paginated API response
func NewPageResponse(data interface{}, page, pageSize, total int) Response {
	// Calculate last page
	lastPage := (total + pageSize - 1) / pageSize
	if lastPage < 1 {
		lastPage = 1
	}

	return Response{
		Data: data,
		Meta: &Metadata{
			CurrentPage:  page,
			PageSize:     pageSize,
			FirstPage:    1,
			LastPage:     lastPage,
			TotalRecords: total,
		},
	}
}
