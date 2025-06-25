import type { Coupon } from './coupon'

export interface MockOrder {
  cost: number
  created_at: string
  coupon_code: string | null
  total_amount: number
  coupon: Coupon | null
}
