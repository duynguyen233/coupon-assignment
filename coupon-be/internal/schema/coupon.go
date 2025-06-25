package schema

import (
	"coupon-be/internal/model"
	"time"
)

type CreateCouponRequest struct {
	CouponCode  *string    `json:"coupon_code" binding:"required"`
	Title       *string    `json:"title" binding:"required"`
	Description *string    `json:"description" binding:"required"`
	CouponType  *string    `json:"coupon_type" binding:"required,oneof=fixed percentage"`
	Usage       *string    `json:"usage" binding:"required"`
	ExpiredAt   *time.Time `json:"expired_at" binding:"required"`
	CouponValue *float64   `json:"coupon_value" binding:"required,gt=0"`
}

type UpdateCouponRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	CouponType  *string    `json:"coupon_type" binding:"omitempty,oneof=fixed percentage"`
	Usage       *string    `json:"usage"`
	ExpiredAt   *time.Time `json:"expired_at"`
	CouponValue *float64   `json:"coupon_value"`
}

type CouponResponse struct {
	CouponCode  string    `json:"coupon_code"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CouponType  string    `json:"coupon_type"`
	Usage       string    `json:"usage"`
	ExpiredAt   time.Time `json:"expired_at"`
	CouponValue float64   `json:"coupon_value"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToCouponResponse(c model.Coupon) CouponResponse {
	return CouponResponse{
		CouponCode:  c.CouponCode,
		Title:       c.Title,
		Description: c.Description,
		CouponType:  c.CouponType,
		Usage:       c.Usage,
		ExpiredAt:   c.ExpiredAt,
		CouponValue: c.CouponValue,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func ToCouponResponses(coupons []model.Coupon) []CouponResponse {
	responses := make([]CouponResponse, len(coupons))
	for i, c := range coupons {
		responses[i] = ToCouponResponse(c)
	}
	return responses
}
