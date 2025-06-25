package schema

import (
	"coupon-be/pkg/utils/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type PaginationResponse[T any] struct {
	Data    []T    `json:"data"`
	Message string `json:"message"`
	Paging  Paging `json:"paging"`
}

type Paging struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"message"`
	Code  int    `json:"code"`
}

// func NewErrorResponse(c *gin.Context, code int, msg string) {
// 	c.AbortWithStatusJSON(code, ErrorResponse{msg})
// }

type ModifyDataResponse struct {
	ID     string `json:"id"`
	Result bool   `json:"result"`
}

func NewErrorResponse(c *gin.Context, err error) {
	switch err.(type) {
	case errs.ScheduleNotfoundError:
		c.JSON(http.StatusOK, ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusNotFound,
		})
	case errs.BadRequestError:
		c.JSON(http.StatusOK, ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusBadRequest,
		})
	case errs.NotFoundError:
		c.JSON(http.StatusOK, ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusNotFound,
		})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal Server Error"})
	}
}
