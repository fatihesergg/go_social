package database

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Limit  int
	Offset int
}

func NewPagination(c *gin.Context) Pagination {
	limit := c.Query("limit")
	offset := c.Query("offset")
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
