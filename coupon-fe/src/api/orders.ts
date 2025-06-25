import type { MockOrder } from '@/types/model/order'
import type { ApiResponse, ErrorResponse } from '@/types/response'
import { apiClient } from './client'

export const orderService = {
  mockOrder: (
    params?: Omit<MockOrder, 'total_amount'> | ErrorResponse,
  ): Promise<ApiResponse<MockOrder> | ErrorResponse> =>
    apiClient.post('v1/orders/mock', params ?? {}),
}
