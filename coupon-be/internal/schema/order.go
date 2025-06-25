package schema

import "time"

type CreateMockOrderRequest struct {
	Cost       float64   `json:"cost" binding:"required"`
	CreatedAt  time.Time `json:"created_at" binding:"required"`
	CouponCode *string   `json:"coupon_code"`
}

type CreateMockOrderResponse struct {
	Cost        float64         `json:"cost"`
	CreatedAt   time.Time       `json:"created_at"`
	CouponCode  *string         `json:"coupon_code,omitempty"`
	TotalAmount float64         `json:"total_amount"`
	Coupon      *CouponResponse `json:"coupon"`
}
