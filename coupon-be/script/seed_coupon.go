package main

import (
	"coupon-be/internal/model"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var SEED_DATA []model.Coupon

func main() {
	db, err := gorm.Open(mysql.Open("root:123123@tcp(localhost:3306)/zalopay?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		fmt.Println("failed to connect to database:", err)
		return
	}
	for i := range 72 {
		var couponType model.CouponType
		var couponUsage model.CouponUsage
		if i%2 == 0 {
			couponType = "fixed"
			couponUsage = "manual"
		} else {
			couponType = "percentage"
			couponUsage = "auto"
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
			fmt.Printf("failed to seed coupon data: %v", err)
			return
		}
	}
}
