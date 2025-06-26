export type CouponType = 'percentage' | 'fixed'
export type CouponUsage = 'manual' | 'auto'

export interface Coupon {
  coupon_code: string
  title: string
  description: string
  coupon_type: CouponType
  usage: CouponUsage
  expired_at: string
  coupon_value: number
  created_at: string
  updated_at: string
}
