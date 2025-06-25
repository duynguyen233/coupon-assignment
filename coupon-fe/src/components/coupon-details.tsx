'use client'

import { couponService } from '@/api/coupon'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import type { Coupon } from '@/types/model/coupon'
import { formatDateTime } from '@/utils/time'
import { ArrowLeft, Calendar, Clock, Edit, FileText, Tag } from 'lucide-react'
import { useEffect, useState } from 'react'

interface CouponDetailsProps {
  coupon_code: string
  onBack: () => void
  onEdit: () => void
}

export default function CouponDetails({ coupon_code, onBack, onEdit }: CouponDetailsProps) {
  const [coupon, setCoupon] = useState<Coupon | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState<boolean>(true)
  const isExpired = (expiredAt: string) => {
    return new Date(expiredAt) < new Date()
  }

  useEffect(() => {
    const fetchCoupon = async () => {
      try {
        const response = await couponService.getCouponByCode(coupon_code)
        if ("error" in response) {
          setError(response.error)
          return
        } else {
          setCoupon(response.data)
        }
      } catch (error) {
        setError('Failed to fetch coupon details. Please try again later.')
      } finally {
        setLoading(false)
      }
    }

    fetchCoupon()
  }, [])

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary mx-auto mb-4"></div>
          <p className="text-muted-foreground">Loading coupon details...</p>
        </div>
      </div>
    )
  }

  if (error || !coupon) {
    return (
      <div className="text-red-500 text-center mt-10">
        <p>{error}</p>
        <Button variant="outline" onClick={onBack} className="mt-4">
          Go Back
        </Button>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" onClick={onBack} size="sm">
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to List
          </Button>
          <div>
            <h2 className="text-2xl font-bold tracking-tight">Coupon Details</h2>
            <p className="text-muted-foreground">View and manage coupon information</p>
          </div>
        </div>
        <Button onClick={onEdit}>
          <Edit className="w-4 h-4 mr-2" />
          Edit Coupon
        </Button>
      </div>

      <Separator />

      <div className="grid gap-6 lg:grid-cols-3">
        {/* Main Details */}
        <div className="lg:col-span-2 space-y-6">
          <Card>
            <CardHeader>
              <div className="flex justify-between items-start">
                <div className="space-y-3">
                  <CardTitle className="text-2xl">{coupon?.title}</CardTitle>
                  <div className="flex items-center gap-3">
                    <Badge variant="outline" className="text-base px-3 py-1">
                      <Tag className="w-4 h-4 mr-1" />
                      {coupon?.coupon_code}
                    </Badge>
                    <Badge variant={isExpired(coupon?.expired_at) ? 'destructive' : 'default'}>
                      {isExpired(coupon?.expired_at) ? 'Expired' : 'Active'}
                    </Badge>
                    <Badge
                      className={`border-transparent bg-[${coupon.coupon_type === 'fixed' ? '#00CF6A' : '#0033C9'}] text-primary-foreground [a&]:bg-green-200`}
                    >
                      {coupon.coupon_type === 'percentage' ? 'Percentage' : 'Fixed'}
                    </Badge>
                  </div>
                </div>
                <div className="text-right">
                  <div className="items-center text-primary text-3xl font-bold">
                    {coupon.coupon_type === 'percentage'
                      ? `${coupon.coupon_value}%`
                      : `${coupon.coupon_value} VND`}
                  </div>
                  <p className="text-sm text-muted-foreground mt-1">
                    {coupon.coupon_type === 'percentage' ? 'Percentage Off' : 'Fixed Discount'}
                  </p>
                </div>
              </div>
            </CardHeader>
            <CardContent className="space-y-6">
              <div>
                <div className="flex items-center mb-3">
                  <FileText className="w-5 h-5 text-muted-foreground mr-2" />
                  <h3 className="font-semibold">Description</h3>
                </div>
                <p className="text-muted-foreground leading-relaxed">{coupon.description}</p>
              </div>

              <Separator />

              <div>
                <div className="flex items-center mb-3">
                  <Clock className="w-5 h-5 text-muted-foreground mr-2" />
                  <h3 className="font-semibold">Usage Instructions</h3>
                </div>
                <p className="text-muted-foreground leading-relaxed">{coupon.usage}</p>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Sidebar Info */}
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Calendar className="w-5 h-5 mr-2" />
                Dates & Times
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Expired At</p>
                <p
                  className={`text-sm font-medium ${isExpired(coupon.expired_at) ? 'text-destructive' : 'text-green-600'}`}
                >
                  {formatDateTime(coupon.expired_at)}
                </p>
              </div>
              <Separator />
              <div>
                <p className="text-sm font-medium text-muted-foreground">Created At</p>
                <p className="text-sm">{formatDateTime(coupon.created_at)}</p>
              </div>
              <Separator />
              <div>
                <p className="text-sm font-medium text-muted-foreground">Updated At</p>
                <p className="text-sm">{formatDateTime(coupon.updated_at)}</p>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Quick Stats</CardTitle>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Status</span>
                <Badge variant={isExpired(coupon?.expired_at) ? 'destructive' : 'default'}>
                  {isExpired(coupon?.expired_at) ? 'Expired' : 'Active'}
                </Badge>
              </div>
              <Separator />
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Type</span>
                <span className="text-sm font-medium capitalize">{coupon.coupon_type}</span>
              </div>
              <Separator />
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Value</span>
                <span className="text-sm font-medium">
                  {coupon.coupon_type === 'percentage'
                    ? `${coupon.coupon_value}%`
                    : `${coupon.coupon_value} VND`}
                </span>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}
