package router

import (
	"coupon-be/internal/schema"
	"coupon-be/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DefaultRoutes struct {
	l logger.Interface
}

func NewDefaultRoutes(handler *gin.RouterGroup, l logger.Interface) *DefaultRoutes {
	r := &DefaultRoutes{l}
	handler.GET("/ping", r.ping)
	return r
}

// @Summary     Ping default
// @Description Ping default
// @ID          ping
// @Tags  	    Default
// @Accept      json
// @Produce     json
// @Success     200 {object} schema.Response[string]
// @Failure     500 {object} schema.ErrorResponse
// @Router      /ping [get]
func (r *DefaultRoutes) ping(c *gin.Context) {
	c.JSON(http.StatusOK, schema.Response[string]{
		Message: "Ping successfully",
		Data:    "pong",
	})
}
