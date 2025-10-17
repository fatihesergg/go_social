package database

import (
	"strconv"
)

type Pagination struct {
	Limit  int
	Offset int
}

type Search struct {
	Query string
}

func NewPagination(limit, offset string) Pagination {
	if limit == "" {
		limit = "20"
	}
	if offset == "" {
		offset = "0"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 20
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0
	}
	return Pagination{
		Limit:  limitInt,
		Offset: offsetInt,
	}
}

func NewSearch(query string) Search {
	if query == "" {
		query = ""
	}
	return Search{
		Query: query,
	}
}
