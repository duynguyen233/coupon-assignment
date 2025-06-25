package services

import (
	"context"
	"coupon-be/internal/model"
	"coupon-be/internal/schema"
	"coupon-be/pkg/logger"
	"fmt"
)

type CouponService interface {
	ValidateCoupon(ctx context.Context, coupon model.Coupon, req schema.CreateMockOrderRequest) (bool, error)
	CalculateAmount(ctx context.Context, coupon *model.Coupon, amount float64) (float64, error)
}

type couponServiceImpl struct {
	l logger.Interface
}

func NewCouponService(l logger.Interface) CouponService {
	return &couponServiceImpl{
		l: l,
	}
}

func (c *couponServiceImpl) ValidateCoupon(ctx context.Context, coupon model.Coupon, req schema.CreateMockOrderRequest) (bool, error) {
	if req.CreatedAt.After(coupon.ExpiredAt) {
		return false, fmt.Errorf("coupon %s is expired", coupon.CouponCode)
	}
	return true, nil
}

func (c *couponServiceImpl) CalculateAmount(ctx context.Context, coupon *model.Coupon, amount float64) (float64, error) {
	// Not implemented yet
	if coupon == nil {
		c.l.Error("Coupon is nil", "amount", amount)
		return amount, nil
	}
	switch (*coupon).CouponType {
	case "fixed":
		return handleFixedCoupon(*coupon, amount)
	case "percentage":
		return handlePercentageCoupon(*coupon, amount)
	default:
		c.l.Error("Invalid coupon type", "coupon_type", coupon.CouponType)
		return 0, fmt.Errorf("invalid coupon type: %s", coupon.CouponType)
	}
}

func handleFixedCoupon(coupon model.Coupon, amount float64) (float64, error) {
	discountedAmount := amount - coupon.CouponValue
	if discountedAmount < 0 {
		return 0, nil // Ensure the total does not go below zero
	}
	return discountedAmount, nil
}

func handlePercentageCoupon(coupon model.Coupon, amount float64) (float64, error) {
	discount := (coupon.CouponValue / 100) * amount
	discountedAmount := amount - discount
	if discountedAmount < 0 {
		return 0, nil // Ensure the total does not go below zero
	}
	return discountedAmount, nil
}
