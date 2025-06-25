import type { Coupon } from '../types/model/coupon'
import type { ApiResponse, ErrorResponse, PaginatedResponse } from '../types/response'
import { apiClient } from './client'

function omitNil<T extends object>(obj: T): Partial<T> {
  return Object.fromEntries(
    Object.entries(obj).filter(([_, v]) => v !== undefined && v !== null),
  ) as Partial<T>
}

export const couponService = {
  getCoupons: (params?: {
    limit?: number
    offset?: number
    coupon_code?: string
  }): Promise<PaginatedResponse<Coupon> | ErrorResponse> =>
    apiClient.get('v1/coupons', omitNil(params ?? {})),
  getCouponByCode: (couponCode: string): Promise<ApiResponse<Coupon> | ErrorResponse> =>
    apiClient.get(`v1/coupons/${couponCode}`),
  createCoupon: (
    coupon: Omit<Coupon, 'created_at' | 'updated_at'>,
  ): Promise<ApiResponse<Coupon> | ErrorResponse> => apiClient.post('v1/coupons', omitNil(coupon)),
  updateCoupon: (
    coupon: Omit<Coupon, 'created_at' | 'updated_at'>,
  ): Promise<ApiResponse<Coupon> | ErrorResponse> =>
    apiClient.put(`v1/coupons/${coupon.coupon_code}`, omitNil(coupon)),
  deleteCoupon: (couponCode: string): Promise<ApiResponse<null> | ErrorResponse> =>
    apiClient.delete(`v1/coupons/${couponCode}`),
}
