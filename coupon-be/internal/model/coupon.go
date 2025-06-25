package model

import "time"

type CouponType string
type CouponUsage string

const (
	CouponTypeFixed      CouponType  = "fixed"
	CouponTypePercentage CouponType  = "percentage"
	CouponUsageManual    CouponUsage = "manual"
	CouponUsageAuto      CouponUsage = "auto"
)

type Coupon struct {
	CouponCode  string      `json:"coupon_code" gorm:"column:coupon_code;type:varchar(255);primaryKey"`
	Title       string      `json:"title" gorm:"column:title;type:varchar(255);not null"`
	Description string      `json:"description" gorm:"column:description;type:text;not null"`
	CouponType  CouponType  `json:"coupon_type" gorm:"column:coupon_type;type:enum('fixed','percentage');not null"`
	Usage       CouponUsage `json:"usage" gorm:"column:usage;type:enum('manual','auto');not null"`
	ExpiredAt   time.Time   `json:"expired_at" gorm:"column:expired_at;type:datetime;not null"`
	CouponValue float64     `json:"coupon_value" gorm:"column:coupon_value;type:decimal(10,2);not null"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
