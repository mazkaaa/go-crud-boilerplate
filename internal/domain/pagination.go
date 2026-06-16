package domain

import "strings"

type PaginationParams struct {
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
}

func (p *PaginationParams) Sanitize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 || p.Limit > 100 {
		p.Limit = 10
	}
	p.SortOrder = strings.ToLower(p.SortOrder)
	if p.SortOrder == "asc" {
		p.SortOrder = "ASC"
	} else {
		p.SortOrder = "DESC"
	}
}

type PaginatedResult struct {
	Items      interface{} `json:"items"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

func ComputeTotalPages(total, limit int) int {
	pages := total / limit
	if total%limit > 0 {
		pages++
	}
	return pages
}
