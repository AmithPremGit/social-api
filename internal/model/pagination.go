package model

import (
	"net/http"
	"strconv"
)

// Pagination contains pagination parameters
type Pagination struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Sort     string `json:"sort"`
	SortBy   string `json:"sort_by"`
}

// DefaultPage is the default page number
const DefaultPage = 1

// DefaultPageSize is the default page size
const DefaultPageSize = 20

// MaxPageSize is the maximum page size allowed
const MaxPageSize = 100

// DefaultSort is the default sort order
const DefaultSort = "desc"

// DefaultSortBy is the default sort field
const DefaultSortBy = "created_at"

// GetPagination extracts pagination parameters from a request
func GetPagination(r *http.Request) Pagination {
	// Initialize with defaults
	pagination := Pagination{
		Page:     DefaultPage,
		PageSize: DefaultPageSize,
		Sort:     DefaultSort,
		SortBy:   DefaultSortBy,
	}

	// Get query parameters
	query := r.URL.Query()

	// Extract page
	if page := query.Get("page"); page != "" {
		if val, err := strconv.Atoi(page); err == nil && val > 0 {
			pagination.Page = val
		}
	}

	// Extract page size
	if pageSize := query.Get("page_size"); pageSize != "" {
		if val, err := strconv.Atoi(pageSize); err == nil && val > 0 {
			// Limit to max page size
			if val > MaxPageSize {
				val = MaxPageSize
			}
			pagination.PageSize = val
		}
	}

	// Extract sort
	if sort := query.Get("sort"); sort != "" {
		if sort == "asc" || sort == "desc" {
			pagination.Sort = sort
		}
	}

	// Extract sort by
	if sortBy := query.Get("sort_by"); sortBy != "" {
		// You might want to validate allowed sort fields here
		pagination.SortBy = sortBy
	}

	return pagination
}

// GetOffset calculates the offset for database queries
func (p Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}
