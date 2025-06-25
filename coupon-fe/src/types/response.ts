export interface ApiResponse<T> {
    data: T
    message: string
    code: number
}

export interface ErrorResponse {
    error: string
    code: number
}

export interface Paging {
    total: number
    limit: number
    offset: number
}

export interface PaginatedResponse<T> {
    data: T[]
    paging: Paging
    message: string
}