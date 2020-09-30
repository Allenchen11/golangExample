package models

import "time"

type Coupon struct {
	CouponID           int       `json:"coupon_id"`
	CouponSerialNumber string    `json:"coupon_serial_number"`
	CouponStatus       string    `json:"coupon_status"`
	CouponUpdateTime   time.Time `json:"coupon_update_time"`
}
