package utils

import (
	"coupon-be/pkg/utils/errs"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaginationParams(c *gin.Context) (int, int, error) {
	var offset, limit int
	if c.Query("offset") != "" {
		var err error
		offset, err = strconv.Atoi(c.Query("offset"))
		if err != nil {
			return 0, 0, err
		}
		if offset < 0 {
			return 0, 0, errs.BadRequestError{Message: "Offset cannot be negative"}
		}
	}

	if c.Query("limit") != "" {
		var err error
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			return 0, 0, err
		}
		if limit <= 0 {
			return 0, 0, errs.BadRequestError{Message: "Limit must be greater than zero"}
		}
	}
	return offset, limit, nil
}
