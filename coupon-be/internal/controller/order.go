package controller

import (
	"context"
	"coupon-be/internal/model"
	"coupon-be/internal/repositories"
	"coupon-be/internal/schema"
	"coupon-be/internal/services"
	"coupon-be/pkg/logger"
	"coupon-be/pkg/utils/errs"
)

type OrderController interface {
	CreateMockOrder(ctx context.Context, req schema.CreateMockOrderRequest) (schema.CreateMockOrderResponse, error)
}

type orderController struct {
	l  logger.Interface
	cr repositories.CouponRepository
	cs services.CouponService
}

func NewOrderController(l logger.Interface, cr repositories.CouponRepository, cs services.CouponService) OrderController {
	return &orderController{
		l:  l,
		cr: cr,
		cs: cs,
	}
}

func (c *orderController) CreateMockOrder(ctx context.Context, req schema.CreateMockOrderRequest) (schema.CreateMockOrderResponse, error) {
	var coupon model.Coupon
	var totalCost float64 = req.Cost
	if req.CouponCode != nil {
		var err error
		coupon, err = c.getAndValidateCoupon(ctx, req)
		if err != nil {
			return schema.CreateMockOrderResponse{}, err
		}
		totalCost, err = c.cs.CalculateAmount(ctx, &coupon, req.Cost)
		if err != nil {
			c.l.Error("Failed to calculate total amount", "error", err)
			return schema.CreateMockOrderResponse{}, errs.BadRequestError{
				Message: "Failed to calculate total amount",
			}
		}
		return schema.CreateMockOrderResponse{
			Cost:        req.Cost,
			CreatedAt:   req.CreatedAt,
			CouponCode:  req.CouponCode,
			TotalAmount: totalCost,
			Coupon: &schema.CouponResponse{
				CouponCode:  coupon.CouponCode,
				Title:       coupon.Title,
				Description: coupon.Description,
				CouponType:  coupon.CouponType,
				CouponValue: coupon.CouponValue,
				ExpiredAt:   coupon.ExpiredAt,
				Usage:       coupon.Usage,
				CreatedAt:   coupon.CreatedAt,
				UpdatedAt:   coupon.UpdatedAt,
			},
		}, nil
	}
	return schema.CreateMockOrderResponse{
		Cost:        req.Cost,
		CreatedAt:   req.CreatedAt,
		CouponCode:  nil,
		TotalAmount: totalCost,
	}, nil
}

func (c *orderController) getAndValidateCoupon(ctx context.Context, req schema.CreateMockOrderRequest) (model.Coupon, error) {
	coupon, err := c.cr.GetCouponByID(ctx, *req.CouponCode)
	if err != nil {
		c.l.Error("Failed to get coupon by code", "coupon_code", *req.CouponCode, "error", err)
		return model.Coupon{}, errs.BadRequestError{
			Message: "Coupon not found",
		}
	}
	isValid, err := c.cs.ValidateCoupon(ctx, coupon, req)
	if err != nil || !isValid {
		c.l.Error("Coupon validation failed", "error", err)
		return model.Coupon{}, errs.BadRequestError{
			Message: err.Error(),
		}
	}
	return coupon, nil
}
