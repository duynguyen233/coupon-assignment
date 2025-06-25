package controller

import (
	"context"
	"coupon-be/internal/model"
	"coupon-be/internal/repositories"
	"coupon-be/internal/schema"
	"coupon-be/internal/services"
	"coupon-be/pkg/logger"
	"coupon-be/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type CouponController interface {
	CreateCoupon(ctx context.Context, coupon schema.CreateCouponRequest) (model.Coupon, error)
	GetCouponsWithTotal(ctx context.Context, offset, limit int, couponCode string) ([]model.Coupon, int64, error)
	GetCouponByID(ctx context.Context, id string) (model.Coupon, error)
	UpdateCoupon(ctx context.Context, id string, coupon schema.UpdateCouponRequest) (model.Coupon, error)
	DeleteCoupon(ctx context.Context, id string) error
}

type couponControllerImpl struct {
	l     logger.Interface
	cs    services.CouponService
	cr    repositories.CouponRepository
	redis *redis.Client
}

func NewCouponController(l logger.Interface, cs services.CouponService, cr repositories.CouponRepository, rc *redis.Client) CouponController {
	return &couponControllerImpl{
		l:     l,
		cs:    cs,
		cr:    cr,
		redis: rc,
	}
}

func (c *couponControllerImpl) CreateCoupon(ctx context.Context, coupon schema.CreateCouponRequest) (model.Coupon, error) {
	couponModel := model.Coupon{
		CouponCode:  *coupon.CouponCode,
		Title:       *coupon.Title,
		Description: *coupon.Description,
		CouponType:  *coupon.CouponType,
		Usage:       *coupon.Usage,
		ExpiredAt:   *coupon.ExpiredAt,
		CouponValue: *coupon.CouponValue,
	}

	couponResponse, err := c.cr.CreateCoupon(ctx, couponModel)
	if err != nil {
		c.l.Error("Failed to create coupon", "error", err, "coupon", couponModel)
		return model.Coupon{}, err
	}
	return couponResponse, nil
}

func (c *couponControllerImpl) GetCouponsWithTotal(ctx context.Context, offset, limit int, couponCode string) ([]model.Coupon, int64, error) {
	if couponCode != "" {
		return c.cr.SearchCouponsWithTotal(ctx, offset, limit, couponCode)
	}
	return c.cr.GetCouponsWithTotal(ctx, offset, limit)
}

func (c *couponControllerImpl) GetCouponByID(ctx context.Context, id string) (model.Coupon, error) {
	hashKey := "coupon:" + id
	couponHash, err := c.redis.HGetAll(ctx, hashKey).Result()
	if err == nil && len(couponHash) > 0 {
		return c.getCouponFromCache(ctx, id)
	}
	coupon, err := c.cr.GetCouponByID(ctx, id)
	if err != nil {
		c.l.Error("Failed to get coupon by ID", "error", err, "id", id)
		return model.Coupon{}, err
	}
	go func() {
		couponMap := utils.StructToMapGetNull(coupon)
		err := c.redis.HSet(context.Background(), hashKey, couponMap).Err()
		if err != nil {
			c.l.Error("Failed to cache coupon", "error", err, "id", id)
			return
		} else {
			c.l.Info("Cached coupon successfully", "id", id)
		}
		_, err = c.redis.Expire(context.Background(), hashKey, CACHE_EXPIRATION*time.Second).Result()
		if err != nil {
			c.l.Error("Failed to set expiration for cached coupon", "error", err, "id", id)
		} else {
			c.l.Info("Set expiration for cached coupon successfully", "id", id, "expiration", CACHE_EXPIRATION)
		}
	}()
	return coupon, nil
}

func (c *couponControllerImpl) UpdateCoupon(ctx context.Context, id string, coupon schema.UpdateCouponRequest) (model.Coupon, error) {
	couponMap := utils.StructToMapGetNull(coupon)
	couponMap["updated_at"] = time.Now()
	couponResponse, err := c.cr.UpdateCoupon(ctx, id, couponMap)
	if err != nil {
		return model.Coupon{}, err
	}
	go func() {
		ctx1 := context.Background()
		hashKey := "coupon:" + id
		couponMap := utils.StructToMapGetNull(couponResponse)
		err := c.redis.HSet(ctx1, hashKey, couponMap).Err()
		if err != nil {
			c.l.Error("Failed to update cached coupon", "error", err, "id", id)
			return
		} else {
			c.l.Info("Updated cached coupon successfully", "id", id)
		}
		_, err = c.redis.Expire(ctx1, hashKey, CACHE_EXPIRATION*time.Second).Result()
		if err != nil {
			c.l.Error("Failed to set expiration for updated cached coupon", "error", err, "id", id)
		} else {
			c.l.Info("Set expiration for updated cached coupon successfully", "id", id, "expiration", CACHE_EXPIRATION)
		}
	}()
	return couponResponse, nil
}

func (c *couponControllerImpl) DeleteCoupon(ctx context.Context, id string) error {
	if err := c.cr.DeleteCoupon(ctx, id); err != nil {
		c.l.Error("Failed to delete coupon", "error", err, "id", id)
		return err
	}
	go func() {
		hashKey := "coupon:" + id
		err := c.redis.Del(ctx, hashKey).Err()
		if err != nil {
			c.l.Error("Failed to delete cached coupon", "error", err, "id", id)
			return
		} else {
			c.l.Info("Deleted cached coupon successfully", "id", id)
		}
	}()
	return nil
}

func (c *couponControllerImpl) getCouponFromCache(ctx context.Context, id string) (model.Coupon, error) {
	hashKey := "coupon:" + id
	couponHash, err := c.redis.HGetAll(ctx, hashKey).Result()
	if err != nil || len(couponHash) == 0 {
		return model.Coupon{}, fmt.Errorf("coupon not found in cache: %w", err)
	}

	couponValue, err := strconv.ParseFloat(couponHash["coupon_value"], 64)
	if err != nil {
		return model.Coupon{}, fmt.Errorf("invalid coupon value in cache: %w", err)
	}
	expiredAt, err := utils.ParseTime(couponHash["expired_at"])
	if err != nil {
		return model.Coupon{}, fmt.Errorf("invalid expired_at in cache: %w", err)
	}
	createdAt, err := utils.ParseTime(couponHash["created_at"])
	if err != nil {
		return model.Coupon{}, fmt.Errorf("invalid created_at in cache: %w", err)
	}
	updatedAt, err := utils.ParseTime(couponHash["updated_at"])
	if err != nil {
		return model.Coupon{}, fmt.Errorf("invalid updated_at in cache: %w", err)
	}

	return model.Coupon{
		CouponCode:  couponHash["coupon_code"],
		Title:       couponHash["title"],
		Description: couponHash["description"],
		CouponType:  model.CouponType(couponHash["coupon_type"]),
		Usage:       model.CouponUsage(couponHash["usage"]),
		CouponValue: couponValue,
		ExpiredAt:   expiredAt,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}
