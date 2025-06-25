package repositories

import (
	"context"
	"coupon-be/internal/model"
	"coupon-be/pkg/utils/errs"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type CouponRepository interface {
	GetCouponsWithTotal(ctx context.Context, offset, limit int) ([]model.Coupon, int64, error)
	SearchCouponsWithTotal(ctx context.Context, offset, limit int, where string) ([]model.Coupon, int64, error)
	GetCouponByID(ctx context.Context, id string) (model.Coupon, error)
	CreateCoupon(ctx context.Context, coupon model.Coupon) (model.Coupon, error)
	UpdateCoupon(ctx context.Context, id string, data map[string]any) (model.Coupon, error)
	DeleteCoupon(ctx context.Context, id string) error
}

type couponRepositoryImpl struct {
	db *gorm.DB
}

func NewCouponRepository(db *gorm.DB) CouponRepository {
	return &couponRepositoryImpl{db: db}
}

func (r *couponRepositoryImpl) GetCouponsWithTotal(ctx context.Context, offset, limit int) ([]model.Coupon, int64, error) {
	var coupons []model.Coupon
	var total int64
	tx := r.db.WithContext(ctx).Model(&model.Coupon{}).Count(&total)
	if offset != 0 || limit != 0 {
		tx = tx.Offset(offset).Limit(limit)
	}
	err := tx.Find(&coupons).Error
	if err != nil {
		return nil, 0, err
	}
	return coupons, total, nil
}

func (r *couponRepositoryImpl) GetCouponByID(ctx context.Context, id string) (model.Coupon, error) {
	var coupon model.Coupon
	if err := r.db.WithContext(ctx).Debug().First(&coupon, "coupon_code = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.Coupon{}, errs.NotFoundError{Message: "Coupon with ID " + id + " not found"}
		}
		return model.Coupon{}, err
	}
	return coupon, nil
}

func (r *couponRepositoryImpl) CreateCoupon(ctx context.Context, coupon model.Coupon) (model.Coupon, error) {
	if err := r.db.WithContext(ctx).Create(&coupon).Error; err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 { // Duplicate entry error code
				return model.Coupon{}, errs.BadRequestError{Message: "Coupon with code " + coupon.CouponCode + " already exists"}
			}
		}
		return model.Coupon{}, err
	}
	return coupon, nil
}

func (r *couponRepositoryImpl) UpdateCoupon(ctx context.Context, id string, data map[string]any) (model.Coupon, error) {
	if err := r.db.WithContext(ctx).Model(&model.Coupon{}).Where("coupon_code = ?", id).Updates(data).Error; err != nil {
		return model.Coupon{}, err
	}
	var coupon model.Coupon
	if err := r.db.WithContext(ctx).First(&coupon, "coupon_code = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.Coupon{}, errs.NotFoundError{Message: "Coupon with ID " + id + " not found"}
		}
		return model.Coupon{}, err
	}
	return coupon, nil
}

func (r *couponRepositoryImpl) DeleteCoupon(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&model.Coupon{}, "coupon_code = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *couponRepositoryImpl) SearchCouponsWithTotal(ctx context.Context, offset, limit int, couponCode string) ([]model.Coupon, int64, error) {
	var coupons []model.Coupon
	var total int64
	tx := r.db.WithContext(ctx).Where("coupon_code LIKE ?", fmt.Sprintf("%%%s%%", couponCode)).Model(&model.Coupon{}).Count(&total)
	if offset != 0 || limit != 0 {
		tx = tx.Offset(offset).Limit(limit)
	}
	err := tx.Find(&coupons).Error
	if err != nil {
		return nil, 0, err
	}
	return coupons, total, nil
}
