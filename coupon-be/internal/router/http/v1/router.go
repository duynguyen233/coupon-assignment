// Package v1 implements routing paths. Each services in own file.
package router

import (
	"coupon-be/internal/controller"
	"coupon-be/pkg/logger"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.RouterGroup,
	l logger.Interface,
	couponController controller.CouponController,
	orderController controller.OrderController,
) {
	// Routers
	h := handler.Group("/v1")
	{
		NewDefaultRoutes(h, l)
		NewCouponRoutes(h, l, couponController)
		NewOrderRoutes(h, l, orderController)
	}

}
