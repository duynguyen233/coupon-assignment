package router

import (
	"coupon-be/internal/controller"
	"coupon-be/internal/schema"
	"coupon-be/pkg/logger"
	"coupon-be/pkg/utils/errs"
	"coupon-be/utils"

	"github.com/gin-gonic/gin"
)

type CouponRoutes struct {
	l                logger.Interface
	couponController controller.CouponController
}

func NewCouponRoutes(handler *gin.RouterGroup, l logger.Interface, couponController controller.CouponController) {
	r := &CouponRoutes{l, couponController}
	h := handler.Group("/coupons")
	{
		h.POST("", r.CreateCoupon)
		h.GET("", r.GetCoupons)
		h.GET("/:id", r.GetCouponByID)
		h.PUT("/:id", r.UpdateCoupon)
		h.DELETE("/:id", r.DeleteCoupon)
	}
}

// @Summary     Create a new coupon
// @Description Create a new coupon
// @ID          createCoupon
// @Tags        Coupons
// @Accept      json
// @Produce     json
// @Param       coupon body schema.CreateCouponRequest true "Coupon data"
// @Success     200 {object} schema.Response[schema.CouponResponse]
// @Failure     400 {object} schema.ErrorResponse
// @Failure     500 {object} schema.ErrorResponse
// @Router      /v1/coupons [post]
func (r *CouponRoutes) CreateCoupon(c *gin.Context) {
	var req schema.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error("Failed to bind JSON for CreateCoupon", "error", err)
		schema.NewErrorResponse(c, errs.BadRequestError{Message: "Invalid request data" + err.Error()})
		return
	}

	coupon, err := r.couponController.CreateCoupon(c.Request.Context(), req)
	if err != nil {
		r.l.Error("Failed to create coupon", "error", err)
		schema.NewErrorResponse(c, err)
		return
	}

	c.JSON(200, schema.Response[schema.CouponResponse]{
		Data:    schema.ToCouponResponse(coupon),
		Message: "Coupon created successfully",
		Code:    200,
	})
}

// @Summary     Get all coupons
// @Description Get all coupons with pagination
// @ID          getCoupons
// @Tags        Coupons
// @Accept      json
// @Produce     json
// @Param       offset query int false "Offset for pagination"
// @Param       limit query int false "Limit for pagination"
// @Param	    coupon_code query string false "Filter by coupon code"
// @Success     200 {object} schema.PaginationResponse[schema.CouponResponse]
// @Failure     500 {object} schema.ErrorResponse
// @Router      /v1/coupons [get]
func (r *CouponRoutes) GetCoupons(c *gin.Context) {
	offset, limit, err := utils.GetPaginationParams(c)
	if err != nil {
		r.l.Error("Failed to parse pagination parameters", "error", err)
		schema.NewErrorResponse(c, errs.BadRequestError{Message: "Invalid pagination parameters"})
		return
	}
	couponCode := c.Query("coupon_code")
	coupons, total, err := r.couponController.GetCouponsWithTotal(c.Request.Context(), offset, limit, couponCode)
	if err != nil {
		r.l.Error("Failed to get coupons", "error", err)
		schema.NewErrorResponse(c, err)
		return
	}

	c.JSON(200, schema.PaginationResponse[schema.CouponResponse]{
		Data:    schema.ToCouponResponses(coupons),
		Message: "Coupons retrieved successfully",
		Paging: schema.Paging{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	})
}

// @Summary     Get a coupon by ID
// @Description Get a coupon by its ID
// @ID          getCouponByID
// @Tags        Coupons
// @Accept      json
// @Produce     json
// @Param       id path string true "Coupon ID"
// @Success     200 {object} schema.Response[CouponResponse]
// @Failure     404 {object} schema.ErrorResponse
// @Failure     500 {object} schema.ErrorResponse
// @Router      /v1/coupons/{id} [get]
func (r *CouponRoutes) GetCouponByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		schema.NewErrorResponse(c, errs.BadRequestError{Message: "Coupon ID is required"})
		return
	}
	coupon, err := r.couponController.GetCouponByID(c.Request.Context(), id)
	if err != nil {
		r.l.Error("Failed to get coupon by ID", "error", err)
		schema.NewErrorResponse(c, err)
		return
	}

	c.JSON(200, schema.Response[schema.CouponResponse]{
		Data:    schema.ToCouponResponse(coupon),
		Message: "Coupon retrieved successfully",
		Code:    200,
	})
}

// @Summary     Update a coupon
// @Description Update a coupon by its ID
// @ID          updateCoupon
// @Tags        Coupons
// @Accept      json
// @Produce     json
// @Param       id path string true "Coupon ID"
// @Param       coupon body schema.UpdateCouponRequest true "Updated coupon data"
// @Success     200 {object} schema.CouponResponse
// @Failure     400 {object} schema.ErrorResponse
// @Failure     404 {object} schema.ErrorResponse
// @Failure     500 {object} schema.ErrorResponse
// @Router      /v1/coupons/{id} [put]
func (r *CouponRoutes) UpdateCoupon(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		schema.NewErrorResponse(c, errs.BadRequestError{Message: "Coupon ID is required"})
		return
	}

	var req schema.UpdateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		schema.NewErrorResponse(c, errs.BadRequestError{Message: "Invalid request data" + err.Error()})
		return
	}

	coupon, err := r.couponController.UpdateCoupon(c.Request.Context(), id, req)
	if err != nil {
		r.l.Error("Failed to update coupon", "error", err)
		schema.NewErrorResponse(c, err)
		return
	}

	c.JSON(200, schema.Response[schema.CouponResponse]{
		Data:    schema.ToCouponResponse(coupon),
		Message: "Coupon updated successfully",
		Code:    200,
	})
}

// @Summary     Delete a coupon
// @Description Delete a coupon by its ID
// @ID          deleteCoupon
// @Tags        Coupons
// @Accept      json
// @Produce     json
// @Param       id path string true "Coupon ID"
// @Success     200 {object} schema.ModifyDataResponse
// @Failure     404 {object} schema.ErrorResponse
// @Failure     500 {object} schema.ErrorResponse
// @Router      /v1/coupons/{id} [delete]
func (r *CouponRoutes) DeleteCoupon(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		schema.NewErrorResponse(c, errs.BadRequestError{Message: "Coupon ID is required"})
		return
	}

	err := r.couponController.DeleteCoupon(c.Request.Context(), id)
	if err != nil {
		r.l.Error("Failed to delete coupon", "error", err)
		schema.NewErrorResponse(c, err)
		return
	}

	c.JSON(200, schema.Response[string]{
		Message: "Coupon deleted successfully",
		Data:    "Coupon with ID " + id + " has been deleted",
		Code:    200,
	})
}
