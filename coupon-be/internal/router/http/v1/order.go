package router

import (
	"coupon-be/internal/controller"
	"coupon-be/internal/schema"
	"coupon-be/pkg/logger"
	"coupon-be/pkg/utils/errs"

	"github.com/gin-gonic/gin"
)

type OrderRoutes struct {
	l               logger.Interface
	orderController controller.OrderController
}

func NewOrderRoutes(handler *gin.RouterGroup, l logger.Interface, orderController controller.OrderController) {
	r := &OrderRoutes{l, orderController}
	h := handler.Group("/orders")
	{
		h.POST("/mock", r.CreateMockOrder)
	}
}

// CreateMockOrder godoc
// @Summary     Create a mock order
// @Description Create a mock order with optional coupon code
// @ID          createMockOrder
// @Tags        Orders
// @Accept      json
// @Produce     json
// @Param       order body schema.CreateMockOrderRequest true "Order data"
// @Success     200 {object} schema.Response[schema.CreateMockOrderResponse]
// @Failure     400 {object} schema.ErrorResponse
// @Failure     500 {object} schema.ErrorResponse
// @Router      /v1/orders/mock [post]
func (r *OrderRoutes) CreateMockOrder(c *gin.Context) {
	var req schema.CreateMockOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error("Failed to bind JSON for CreateMockOrder", "error", err)
		schema.NewErrorResponse(c, errs.BadRequestError{Message: "Invalid request data: " + err.Error()})
		return
	}

	order, err := r.orderController.CreateMockOrder(c.Request.Context(), req)
	if err != nil {
		schema.NewErrorResponse(c, err)
		return
	}

	c.JSON(200, schema.Response[schema.CreateMockOrderResponse]{
		Data:    order,
		Message: "Mock order created successfully",
		Code:    200,
	})
}
