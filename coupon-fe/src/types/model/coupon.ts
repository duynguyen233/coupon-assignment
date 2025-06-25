export type CouponType = "percentage" | "fixed"

export interface Coupon {
    coupon_code: string
    title: string
    description: string
    coupon_type: CouponType
    usage: string
    expired_at: string
    coupon_value: number
    created_at: string
    updated_at: string
}
