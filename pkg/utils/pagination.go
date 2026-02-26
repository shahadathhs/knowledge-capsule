package utils

import (
	"net/http"
	"strconv"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

// ParsePagination extracts page and limit from request query params.
// Returns (page, limit) with defaults and limits applied.
func ParsePagination(r *http.Request) (page, limit int) {
	page = DefaultPage
	limit = DefaultLimit

	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
			if limit > MaxLimit {
				limit = MaxLimit
			}
		}
	}
	return page, limit
}

// SlicePage returns the paginated slice of items.
// Items must be a slice; returns (paginated slice, total count).
func SlicePage[T any](items []T, page, limit int) ([]T, int) {
	total := len(items)
	start := (page - 1) * limit
	if start >= total {
		return []T{}, total
	}
	end := start + limit
	if end > total {
		end = total
	}
	return items[start:end], total
}
