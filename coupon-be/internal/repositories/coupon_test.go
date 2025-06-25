package repositories

import (
	"context"
	"coupon-be/internal/model"
	"coupon-be/pkg/utils/errs"
	"testing"
	"time"

	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SEED_DATA = []model.Coupon{}

func InitializeCouponRepository(t *testing.T) CouponRepository {
	db, err := gorm.Open(mysql.Open("root:123123@tcp(localhost:3306)/zalopay?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	for i := range 72 {
		var couponType model.CouponType
		var couponUsage model.CouponUsage
		if i%2 == 0 {
			couponType = model.CouponTypeFixed
			couponUsage = model.CouponUsageManual
		} else {
			couponType = model.CouponTypePercentage
			couponUsage = model.CouponUsageAuto
		}
		SEED_DATA = append(SEED_DATA, model.Coupon{
			CouponCode:  "TEST" + fmt.Sprint(i),
			Title:       "Test Coupon " + fmt.Sprint(i),
			Description: "Description for Test Coupon " + fmt.Sprint(i),
			CouponType:  couponType,
			Usage:       couponUsage,
			ExpiredAt:   time.Now().AddDate(0, 0, 10),
			CouponValue: float64(i * 10),
		})
		err = db.Create(&SEED_DATA[i]).Error
		if err != nil {
			t.Fatalf("Failed to seed database: %v", err)
		}
	}
	return NewCouponRepository(db)
}

func RemoveDatabaseSeed(t *testing.T) {
	db, err := gorm.Open(mysql.Open("root:123123@tcp(localhost:3306)/zalopay?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	err = db.Exec("DELETE FROM coupons").Error
	if err != nil {
		t.Fatalf("Failed to remove seed data: %v", err)
	}
}

func TestGetCouponsWithTotal(t *testing.T) {
	repo := InitializeCouponRepository(t)
	type args struct {
		ctx    context.Context
		offset int
		limit  int
	}
	type want struct {
		coupons []model.Coupon
		total   int64
		err     error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Get all coupons with offset 0 and limit 10",
			args: args{
				ctx:    context.Background(),
				offset: 0,
				limit:  10,
			},
			want: want{
				coupons: SEED_DATA[:10],
				total:   int64(len(SEED_DATA)),
				err:     nil,
			},
		},
		{
			name: "Get all coupons if offset and limit are 0",
			args: args{
				ctx:    context.Background(),
				offset: 0,
				limit:  0,
			},
			want: want{
				coupons: SEED_DATA,
				total:   int64(len(SEED_DATA)),
				err:     nil,
			},
		},
		{
			name: "Get all coupons if offset is greater than total coupons",
			args: args{
				ctx:    context.Background(),
				offset: 100,
				limit:  10,
			},
			want: want{
				coupons: []model.Coupon{},
				total:   int64(len(SEED_DATA)),
				err:     nil,
			},
		},
		{
			name: "Get all coupons with offset 0 and limit greater than total coupons",
			args: args{
				ctx:    context.Background(),
				offset: 0,
				limit:  100,
			},
			want: want{
				coupons: SEED_DATA,
				total:   int64(len(SEED_DATA)),
				err:     nil,
			},
		},
		{
			name: "Get all coupons with offset 10 and limit 10",
			args: args{
				ctx:    context.Background(),
				offset: 10,
				limit:  10,
			},
			want: want{
				coupons: SEED_DATA[10:20],
				total:   int64(len(SEED_DATA)),
				err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coupons, total, err := repo.GetCouponsWithTotal(tt.args.ctx, tt.args.offset, tt.args.limit)
			if err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetCouponsWithTotal(), test name: %s, error = %v, wantErr %v", tt.name, err, tt.want.err)
			}
			if total != tt.want.total {
				t.Errorf("GetCouponsWithTotal(), test name: %s, total = %v, want %v", tt.name, total, tt.want.total)
			}
			if len(coupons) != len(tt.want.coupons) {
				t.Errorf("GetCouponsWithTotal(), test name: %s coupons length = %v, want %v", tt.name, len(coupons), len(tt.want.coupons))
			}
		})
	}
	RemoveDatabaseSeed(t)
}

func TestGetCouponByID(t *testing.T) {
	repo := InitializeCouponRepository(t)
	type args struct {
		ctx context.Context
		id  string
	}
	type want struct {
		coupon model.Coupon
		err    error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Get coupon by valid ID",
			args: args{
				ctx: context.Background(),
				id:  SEED_DATA[0].CouponCode,
			},
			want: want{
				coupon: SEED_DATA[0],
				err:    nil,
			},
		},
		{
			name: "Get coupon by invalid ID",
			args: args{
				ctx: context.Background(),
				id:  "INVALID_ID",
			},
			want: want{
				coupon: model.Coupon{},
				err:    errs.NotFoundError{Message: "Coupon with ID INVALID_ID not found"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coupon, err := repo.GetCouponByID(tt.args.ctx, tt.args.id)
			if err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("GetCouponByID(), test name: %s, error = %v, wantErr %v", tt.name, err, tt.want.err)
			}
			if coupon.CouponCode != tt.want.coupon.CouponCode {
				t.Errorf("GetCouponByID(), test name: %s, coupon = %v, want %v", tt.name, coupon.CouponCode, tt.want.coupon.CouponCode)
			}
		})
	}
	RemoveDatabaseSeed(t)
}

func TestCreateCoupon(t *testing.T) {
	repo := InitializeCouponRepository(t)
	type args struct {
		ctx    context.Context
		coupon model.Coupon
	}
	type want struct {
		coupon model.Coupon
		err    error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Create coupon with valid data",
			args: args{
				ctx: context.Background(),
				coupon: model.Coupon{
					CouponCode:  "TEST100",
					Title:       "Test Coupon 100",
					Description: "Description for Test Coupon 100",
					CouponType:  "fixed",
					Usage:       "single",
					ExpiredAt:   time.Now().AddDate(0, 0, 10),
					CouponValue: 100.0,
				},
			},
			want: want{
				coupon: model.Coupon{
					CouponCode:  "TEST100",
					Title:       "Test Coupon 100",
					Description: "Description for Test Coupon 100",
					CouponType:  "fixed",
					Usage:       "single",
					ExpiredAt:   time.Now().AddDate(0, 0, 10),
					CouponValue: 100.0,
				},
				err: nil,
			},
		},
		{
			name: "Create coupon with duplicate code",
			args: args{
				ctx: context.Background(),
				coupon: model.Coupon{
					CouponCode:  SEED_DATA[0].CouponCode,
					Title:       "Duplicate Coupon",
					Description: "This should fail due to duplicate code",
					CouponType:  "fixed",
					Usage:       "single",
					ExpiredAt:   time.Now().AddDate(0, 0, 10),
					CouponValue: 50.0,
				},
			},
			want: want{
				coupon: model.Coupon{},
				err:    errs.BadRequestError{Message: "Coupon with code " + SEED_DATA[0].CouponCode + " already exists"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coupon, err := repo.CreateCoupon(tt.args.ctx, tt.args.coupon)
			if err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("CreateCoupon(), test name: %s, error = %v, wantErr %v", tt.name, err, tt.want.err)
			}
			if coupon.CouponCode != tt.want.coupon.CouponCode {
				t.Errorf("CreateCoupon(), test name: %s, coupon = %v, want %v", tt.name, coupon.CouponCode, tt.want.coupon.CouponCode)
			}
		})
	}
	RemoveDatabaseSeed(t)
}

func TestUpdateCoupon(t *testing.T) {
	repo := InitializeCouponRepository(t)
	type args struct {
		ctx  context.Context
		id   string
		data map[string]any
	}
	type want struct {
		coupon model.Coupon
		err    error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Update coupon with valid ID and data",
			args: args{
				ctx: context.Background(),
				id:  SEED_DATA[0].CouponCode,
				data: map[string]any{
					"title":       "Updated Test Coupon",
					"description": "Updated Description for Test Coupon",
				},
			},
			want: want{
				coupon: model.Coupon{
					CouponCode:  SEED_DATA[0].CouponCode,
					Title:       "Updated Test Coupon",
					Description: "Updated Description for Test Coupon",
					CouponType:  SEED_DATA[0].CouponType,
					Usage:       SEED_DATA[0].Usage,
					ExpiredAt:   SEED_DATA[0].ExpiredAt,
					CouponValue: SEED_DATA[0].CouponValue,
				},
				err: nil,
			},
		},
		{
			name: "Update coupon with invalid ID",
			args: args{
				ctx: context.Background(),
				id:  "INVALID_ID",
				data: map[string]any{
					"title":       "This should fail",
					"description": "This should fail too",
				},
			},
			want: want{
				coupon: model.Coupon{},
				err:    errs.NotFoundError{Message: "Coupon with ID INVALID_ID not found"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coupon, err := repo.UpdateCoupon(tt.args.ctx, tt.args.id, tt.args.data)
			if err != nil && tt.want.err != nil && err.Error() != tt.want.err.Error() {
				t.Errorf("UpdateCoupon(), test name: %s, error = %v, wantErr %v", tt.name, err, tt.want.err)
			}
			if err != nil && tt.want.err == nil {
				t.Errorf("UpdateCoupon(), test name: %s, unexpected error = %v", tt.name, err)
			}
			if coupon.CouponCode != tt.want.coupon.CouponCode {
				t.Errorf("UpdateCoupon(), test name: %s, coupon = %v, want %v", tt.name, coupon.CouponCode, tt.want.coupon.CouponCode)
			}
			if coupon.Title != tt.want.coupon.Title {
				t.Errorf("UpdateCoupon(), test name: %s, coupon title = %v, want %v", tt.name, coupon.Title, tt.want.coupon.Title)
			}
			if coupon.Description != tt.want.coupon.Description {
				t.Errorf("UpdateCoupon(), test name: %s, coupon description = %v, want %v", tt.name, coupon.Description, tt.want.coupon.Description)
			}
		})
	}
	RemoveDatabaseSeed(t)
}

func TestDeleteCoupon(t *testing.T) {
	repo := InitializeCouponRepository(t)
	type args struct {
		ctx context.Context
		id  string
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Delete coupon with valid ID",
			args: args{
				ctx: context.Background(),
				id:  SEED_DATA[0].CouponCode,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "Delete coupon with invalid ID",
			args: args{
				ctx: context.Background(),
				id:  "INVALID_ID",
			},
			want: want{
				err: gorm.ErrRecordNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.DeleteCoupon(tt.args.ctx, tt.args.id)
			if err != nil && (tt.want.err == nil || err.Error() != tt.want.err.Error()) {
				t.Errorf("DeleteCoupon(), test name: %s, error = %v, wantErr %v", tt.name, err, tt.want.err)
			}
		})
	}
	RemoveDatabaseSeed(t)
}
