import { couponService } from '@/api/coupon'
import { Alert, AlertDescription } from '@/components/ui/alert'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import { PAGINATION_LIMIT } from '@/constants/paginate'
import type { Coupon } from '@/types/model/coupon'
import { formatDateTime, isExpired } from '@/utils/time'
import _ from 'lodash'
import { AlertCircle, Edit, Eye, Loader2, Package, SearchIcon, Tag, Trash2 } from 'lucide-react'
import { useCallback, useEffect, useRef, useState } from 'react'
import { Input } from './ui/input'

interface CouponListProps {
  onViewCoupon: (coupon_code: string) => void
  onEditCoupon: (coupon_code: string) => void
}

export default function CouponList({ onViewCoupon, onEditCoupon }: CouponListProps) {
  const [coupons, setCoupons] = useState<Coupon[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [total, setTotal] = useState<number>(0)
  const [currentPage, setCurrentPage] = useState(1)
  const [couponCodeQuery, setCouponCodeQuery] = useState<string | undefined>('')
  const [hasMore, setHasMore] = useState(true)
  const [loadingMore, setLoadingMore] = useState(false)
  const [searchLoading, setSearchLoading] = useState(false)

  const observerTarget = useRef<HTMLDivElement>(null)

  const fetchCoupons = useCallback(
    async (page: number, searchQuery: string = '', isNewSearch = false) => {
      if (isNewSearch || page === 1) {
        setSearchLoading(true)
      } else if (page === 1) {
        setLoading(true)
      } else {
        setLoadingMore(true)
      }
      setError(null)
      try {
        const data = await couponService.getCoupons({
          limit: PAGINATION_LIMIT,
          offset: (page - 1) * PAGINATION_LIMIT,
          coupon_code: searchQuery || undefined,
        })
        if ('error' in data) {
          setError(data.error)
        } else {
          if (isNewSearch || page === 1) {
            setCoupons(data.data as Coupon[])
            setTotal(data.paging.total as number)
          } else {
            await new Promise((resolve) => setTimeout(resolve, 1000)) // Simulate network delay
            setCoupons((prev) => [...prev, ...(data.data as Coupon[])])
            setTotal(data.paging.total as number)
          }
          setHasMore(data.paging.total > page * PAGINATION_LIMIT)
        }
      } catch (err) {
        console.error('Failed to fetch coupons:', err)
        setError('Failed to fetch coupons')
      } finally {
        setLoading(false)
        setLoadingMore(false)
        setSearchLoading(false)
      }
    },
    [],
  )

  const debouncedSearch = useCallback(
    _.debounce(async (coupon_code: string) => {
      setCurrentPage(1)
      await fetchCoupons(1, coupon_code, true)
    }, 300),
    [fetchCoupons],
  )

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        const target = entries[0]
        if (target.isIntersecting && hasMore && !loading && !loadingMore) {
          const nextPage = currentPage + 1
          setCurrentPage(nextPage)
          fetchCoupons(nextPage, couponCodeQuery)
        }
      },
      {
        root: null,
        rootMargin: '0px',
        threshold: 0.1,
      },
    )

    const currentTarget = observerTarget.current
    if (currentTarget) {
      observer.observe(currentTarget)
    }
    return () => {
      if (currentTarget) {
        observer.unobserve(currentTarget)
      }
    }
  }, [hasMore, loading, loadingMore, currentPage, couponCodeQuery, fetchCoupons])

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setCouponCodeQuery(value)
    debouncedSearch(value)
  }

  useEffect(() => {
    fetchCoupons(1)
  }, [fetchCoupons])

  const handleDelete = async (couponCode: string) => {
    setLoading(true)
    setError(null)
    try {
      const response = await couponService.deleteCoupon(couponCode)
      if ('error' in response) {
        setError(response.error)
      } else {
        setCoupons((prev) => prev.filter((coupon) => coupon.coupon_code !== couponCode))
        setTotal((prev) => prev - 1)
      }
    } catch (err) {
      console.error('Failed to delete coupon:', err)
      setError('Failed to delete coupon')
    } finally {
      setLoading(false)
    }
  }

  if (loading && !searchLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading coupons...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <Alert variant="destructive">
        <AlertCircle className="h-4 w-4" />
        <AlertDescription>{error}</AlertDescription>
      </Alert>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">All Coupons</h2>
          <p className="text-muted-foreground">Manage your discount coupons</p>
        </div>
        <div className="flex items-right gap-6">
          <div className="grid w-full max-w-sm items-center gap-1.5">
            <div className="relative">
              <div className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground">
                <SearchIcon className="h-4 w-4" />
              </div>
              <Input
                id="search"
                type="search"
                placeholder="Search..."
                value={couponCodeQuery}
                className="w-full rounded-lg bg-background pl-8"
                onChange={handleSearchChange}
              />
            </div>
          </div>
          <Badge variant="secondary" className="text-sm">
            {total} Total
          </Badge>
        </div>
      </div>

      <Separator />

      {coupons.length === 0 ? (
        <Card>
          <CardContent className="flex flex-col items-center justify-center py-12">
            <Package className="h-12 w-12 text-muted-foreground mb-4" />
            <h3 className="text-lg font-semibold mb-2">No coupons found</h3>
            <p className="text-muted-foreground text-center mb-4">
              Create your first coupon to get started with discount management
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {coupons.map((coupon) => (
            <Card key={coupon.coupon_code} className="hover:shadow-md transition-shadow">
              <CardHeader className="pb-3">
                <div className="flex justify-between items-start">
                  <div className="space-y-2">
                    <CardTitle className="text-lg">{coupon.title}</CardTitle>
                    <div className="flex items-center gap-2">
                      <Badge variant="outline">
                        <Tag className="h-4 w-4 mr-1" />
                        <span>{coupon.coupon_code}</span>
                      </Badge>
                      <Badge variant={isExpired(coupon.expired_at) ? 'destructive' : 'default'}>
                        {isExpired(coupon.expired_at) ? 'Expired' : 'Active'}
                      </Badge>
                      <Badge
                        className={`border-transparent text-primary-foreground ${
                          coupon.coupon_type === 'fixed' ? 'bg-[#00CF6A]' : 'bg-[#0033C9]'
                        }`}
                      >
                        {coupon.coupon_type === 'percentage' ? 'Percentage' : 'Fixed'}
                      </Badge>
                    </div>
                  </div>
                  <div className="flex items-center text-primary">
                    <span className="font-bold ml-1">
                      {coupon.coupon_type === 'percentage'
                        ? `${coupon.coupon_value}%`
                        : `${coupon.coupon_value} VND`}
                    </span>
                  </div>
                </div>
              </CardHeader>
              <CardContent className="space-y-4">
                <p className="text-sm text-muted-foreground line-clamp-2">{coupon.description}</p>
                <p className="text-xs text-muted-foreground">
                  Expires: {formatDateTime(coupon.expired_at)}
                </p>
                <Separator />
                <div className="flex gap-2">
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => onViewCoupon(coupon.coupon_code)}
                    className="flex-1"
                  >
                    <Eye className="w-4 h-4 mr-1" />
                    View
                  </Button>
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => onEditCoupon(coupon.coupon_code)}
                    className="flex-1"
                  >
                    <Edit className="w-4 h-4 mr-1" />
                    Edit
                  </Button>
                  <AlertDialog>
                    <AlertDialogTrigger asChild>
                      <Button size="sm" variant="outline">
                        <Trash2 className="w-4 h-4" />
                      </Button>
                    </AlertDialogTrigger>
                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>Delete Coupon</AlertDialogTitle>
                        <AlertDialogDescription>
                          Are you sure you want to delete "{coupon.title}"? This action cannot be
                          undone.
                        </AlertDialogDescription>
                      </AlertDialogHeader>
                      <AlertDialogFooter>
                        <AlertDialogCancel>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={() => handleDelete(coupon.coupon_code)}>
                          Delete
                        </AlertDialogAction>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
      {/* Loading indicator for lazy loading */}
      {loadingMore && (
        <div className="flex items-center justify-center py-8">
          <div className="flex items-center gap-2">
            <Loader2 className="h-4 w-4 animate-spin" />
            <span className="text-sm text-muted-foreground">Loading more coupons...</span>
          </div>
        </div>
      )}

      {/* Intersection observer target */}
      {hasMore && !loadingMore && <div ref={observerTarget} className="h-4" />}

      {/* End of list indicator */}
      {!hasMore && coupons.length > 0 && (
        <div className="flex items-center justify-center py-8">
          <p className="text-sm text-muted-foreground">You've reached the end of the list</p>
        </div>
      )}
    </div>
  )
}
