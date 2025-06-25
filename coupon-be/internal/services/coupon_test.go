package services

import (
	"context"
	"coupon-be/internal/model"
	"coupon-be/internal/schema"
	"coupon-be/pkg/logger"
	"testing"
	"time"
)

func TestValidateCoupon(t *testing.T) {
	logger := logger.New("test")
	cs := NewCouponService(logger)
	testString := "TEST123"
	tests := []struct {
		name   string
		coupon model.Coupon
		req    schema.CreateMockOrderRequest
		want   bool
	}{
		{
			name: "TC1.1: Valid Coupon",
			coupon: model.Coupon{
				CouponCode:  testString,
				Title:       "Test Coupon",
				Description: "This is a test coupon",
				CouponType:  "fixed",
				Usage:       "lorem ipsum",
				ExpiredAt:   time.Now().Add(24 * time.Hour),
				CouponValue: 15000,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			req: schema.CreateMockOrderRequest{
				CouponCode: &testString,
				Cost:       100000,
				CreatedAt:  time.Now(),
			},
			want: true,
		},
		{
			name: "TC1.2: Expired Coupon",
			coupon: model.Coupon{
				CouponCode:  testString,
				Title:       "Expired Coupon",
				Description: "This coupon is expired",
				CouponType:  "fixed",
				Usage:       "lorem ipsum",
				ExpiredAt:   time.Now().Add(-24 * time.Hour), // Expired
				CouponValue: 15000,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			req: schema.CreateMockOrderRequest{
				CouponCode: &testString,
				Cost:       100000,
				CreatedAt:  time.Now(),
			},
			want: false,
		},
		{
			name: "TC1.3: Valid Coupon with Different Cost",
			coupon: model.Coupon{
				CouponCode:  testString,
				Title:       "Test Coupon",
				Description: "This is a test coupon",
				CouponType:  "percentage",
				Usage:       "lorem ipsum",
				ExpiredAt:   time.Now().Add(24 * time.Hour),
				CouponValue: 20, // 20% discount
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			req: schema.CreateMockOrderRequest{
				CouponCode: &testString,
				Cost:       50000,
				CreatedAt:  time.Now(),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cs.ValidateCoupon(context.Background(), tt.coupon, tt.req)
			if err != nil {
				t.Errorf("ValidateCoupon() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateCoupon() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalculateAmount(t *testing.T) {
	logger := logger.New("test")
	cs := NewCouponService(logger)
	testString := "TEST123"
	tests := []struct {
		name    string
		coupon  *model.Coupon
		amount  float64
		want    float64
		wantErr bool
	}{
		{
			name: "TC2.1: Fixed Coupon",
			coupon: &model.Coupon{
				CouponCode:  testString,
				Title:       "Test Fixed Coupon",
				Description: "This is a test fixed coupon",
				Usage:       "lorem ipsum",
				ExpiredAt:   time.Now().Add(24 * time.Hour),
				CouponType:  "fixed",
				CouponValue: 15000,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			amount: 100000,
			want:   85000,
		},
		{
			name: "TC2.2: Fixed Coupon with Zero ",
			coupon: &model.Coupon{
				CouponCode:  testString,
				Title:       "Test Fixed Coupon with Zero",
				Description: "This is a test fixed coupon with zero value",
				Usage:       "lorem ipsum",
				ExpiredAt:   time.Now().Add(24 * time.Hour),
				CouponType:  "fixed",
				CouponValue: 100000,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			amount: 100000,
			want:   0, // Total should not go below zero
		},
		{
			name: "TC2.3: Fixed Coupon with Negative Amount",
			coupon: &model.Coupon{
				CouponCode:  testString,
				Title:       "Test Fixed Coupon with Negative Amount",
				Description: "This is a test fixed coupon with negative amount",
				Usage:       "lorem ipsum",
				ExpiredAt:   time.Now().Add(24 * time.Hour),
				CouponType:  "fixed",
				CouponValue: 15000,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			amount: 10000, // Amount is less than coupon value
			want:   0,     // Total should not go below zero
		},
		{
			name: "TC2.4: Percentage Coupon",
			coupon: &model.Coupon{
				CouponCode:  testString,
				Title:       "Test Percentage Coupon",
				Description: "This is a test percentage coupon",
				Usage:       "lorem ipsum",
				ExpiredAt:   time.Now().Add(24 * time.Hour),
				CouponType:  "percentage",
				CouponValue: 20, // 20% discount
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			amount: 100000,
			want:   80000, // 20% off of 100000
		},
		{
			name: "TC2.5: Percentage Coupon with Zero",
			coupon: &model.Coupon{
				CouponCode:  testString,
				Title:       "Test Percentage Coupon with Zero",
				Description: "This is a test percentage coupon with zero value",
				Usage:       "lorem ipsum",
				ExpiredAt:   time.Now().Add(24 * time.Hour),
				CouponType:  "percentage",
				CouponValue: 100, // 100% discount
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			amount: 100000,
			want:   0, // Total should not go below zero
		},
		{
			name: "TC2.6: Percentage Coupon with Negative Amount",
			coupon: &model.Coupon{
				CouponCode:  testString,
				Title:       "Test Percentage Coupon with Negative Amount",
				Description: "This is a test percentage coupon with negative amount",
				Usage:       "lorem ipsum",
				ExpiredAt:   time.Now().Add(24 * time.Hour),
				CouponType:  "percentage",
				CouponValue: 20, // 20% discount
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			amount: 10000, // Amount is less than coupon value
			want:   8000,  // 20% off of 10000
		},
		{
			name:    "TC2.7: No Coupon",
			coupon:  nil,
			amount:  100000,
			want:    100000, // No coupon applied, total should be the same
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cs.CalculateAmount(context.Background(), tt.coupon, tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateAmount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CalculateAmount() got = %v, want %v", got, tt.want)
			}
		})
	}
}
