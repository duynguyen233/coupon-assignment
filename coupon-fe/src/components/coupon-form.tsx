'use client'

import { couponService } from '@/api/coupon'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Textarea } from '@/components/ui/textarea'
import type { CouponType, CouponUsage } from '@/types/model/coupon'
import { formatDate } from '@/utils/time'
import { AlertCircle, CheckCircle, ChevronDownIcon, Save } from 'lucide-react'
import type React from 'react'
import { useEffect, useState } from 'react'
import { Calendar } from './ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover'
import { RadioGroup, RadioGroupItem } from './ui/radio-group'

interface Coupon {
  id?: number
  coupon_code: string
  title: string
  description: string
  coupon_type: CouponType
  usage: CouponUsage
  expired_at: string
  coupon_value: number
}

interface CouponFormProps {
  coupon_code?: string | null
  onBack: () => void
}

export default function CouponForm({ coupon_code, onBack }: CouponFormProps) {
  const [formData, setFormData] = useState<Coupon>({
    coupon_code: '',
    title: '',
    description: '',
    coupon_type: 'percentage' as 'fixed' | 'percentage',
    usage: 'manual' as 'manual' | 'auto',
    expired_at: '',
    coupon_value: 0,
  })
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState(false)
  const [open, setOpen] = useState(false)

  useEffect(() => {
    const fetchCoupon = async () => {
      if (!coupon_code) return

      try {
        const response = await couponService.getCouponByCode(coupon_code)
        if ('error' in response) {
          setError(response.error)
          return
        }
        setFormData(response.data)
      } catch (error) {
        setError('Failed to fetch coupon details. Please try again later.')
      }
    }
    fetchCoupon()
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError(null)
    setSuccess(false)

    try {
      const response = coupon_code
        ? await couponService.updateCoupon(formData)
        : await couponService.createCoupon(formData)

      if ('error' in response) {
        throw new Error(response.error)
      }

      setSuccess(true)
      setTimeout(() => {
        onBack()
      }, 1500)
    } catch (error) {
      console.error('Error saving coupon:', error)
      setError(error instanceof Error ? error.message : 'Failed to save coupon')
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (field: string, value: string | number) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }))
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <div>
          <h2 className="text-2xl font-bold tracking-tight">
            {coupon_code ? 'Edit Coupon' : 'Create New Coupon'}
          </h2>
          <p className="text-muted-foreground">
            {coupon_code ? 'Update coupon details' : 'Add a new discount coupon'}
          </p>
        </div>
      </div>

      <Separator />

      {error && (
        <Alert variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}

      {success && (
        <Alert>
          <CheckCircle className="h-4 w-4" />
          <AlertDescription>
            Coupon {coupon_code ? 'updated' : 'created'} successfully! Redirecting...
          </AlertDescription>
        </Alert>
      )}

      <Card className="max-w-xl center mx-auto">
        <CardHeader>
          <CardTitle>Coupon Details</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <Label htmlFor="coupon_code">Coupon Code *</Label>
                <Input
                  id="coupon_code"
                  value={formData.coupon_code}
                  onChange={(e) => handleChange('coupon_code', e.target.value.toUpperCase())}
                  placeholder="e.g., ZLP20K"
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="title">Title *</Label>
                <Input
                  id="title"
                  value={formData.title}
                  onChange={(e) => handleChange('title', e.target.value)}
                  placeholder="e.g., 20% Off Everything"
                  required
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="description">Description *</Label>
              <Textarea
                id="description"
                value={formData.description}
                onChange={(e) => handleChange('description', e.target.value)}
                placeholder="Describe the coupon offer..."
                required
                rows={3}
              />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <Label htmlFor="coupon_type">Coupon Type *</Label>
                <RadioGroup
                  onValueChange={(value: 'fixed' | 'percentage') =>
                    handleChange('coupon_type', value as 'fixed' | 'percentage')
                  }
                  defaultValue={formData.coupon_type}
                  className="flex flex-col"
                >
                  <div className="flex items-center space-x-2">
                    <RadioGroupItem value="fixed" id="fixed" />
                    <Label htmlFor="fixed">Fixed Coupon</Label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <RadioGroupItem value="percentage" id="percentage" />
                    <Label htmlFor="percentage">Percentage Coupon</Label>
                  </div>
                </RadioGroup>
              </div>

              <div className="space-y-2">
                <Label htmlFor="coupon_value">
                  {formData.coupon_type === 'percentage' ? 'Percentage (%)' : 'Amount (VND)'} *
                </Label>
                <Input
                  id="coupon_value"
                  type="number"
                  step={formData.coupon_type === 'percentage' ? '1' : '1000'}
                  min="0"
                  max={formData.coupon_type === 'percentage' ? '100' : undefined}
                  value={formData.coupon_value}
                  onChange={(e) => handleChange('coupon_value', Number.parseFloat(e.target.value))}
                  placeholder={formData.coupon_type === 'percentage' ? '20' : '1000'}
                  required
                />
              </div>
            </div>

            <div className="space-y-4">
              <Label htmlFor="usage">Usage Type *</Label>
              <RadioGroup
                onValueChange={(value: CouponUsage) => handleChange('usage', value)}
                defaultValue={formData.usage}
                className="flex flex-col"
              >
                <div className="flex items-center space-x-2">
                  <RadioGroupItem value="manual" id="manual" />
                  <Label htmlFor="manual">Manual Apply</Label>
                </div>
                <div className="flex items-center space-x-2">
                  <RadioGroupItem value="auto" id="auto" />
                  <Label htmlFor="auto">Automatic Apply</Label>
                </div>
              </RadioGroup>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="space-y-2">
                <Label htmlFor="date-picker" className="px-1">
                  Expiry Date *
                </Label>
                <Popover open={open} onOpenChange={setOpen}>
                  <PopoverTrigger asChild>
                    <Button
                      variant="outline"
                      id="date-picker"
                      className="w-full justify-between font-normal"
                    >
                      {formData.expired_at ? formatDate(formData.expired_at) : 'Select date'}
                      <ChevronDownIcon />
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent className="w-auto overflow-hidden p-0" align="start">
                    <Calendar
                      mode="single"
                      selected={formData.expired_at ? new Date(formData.expired_at) : undefined}
                      captionLayout="dropdown"
                      onSelect={(date) => {
                        date?.setTime(new Date().getTime())
                        handleChange('expired_at', date ? date.toISOString() : '')
                        setOpen(false)
                      }}
                    />
                  </PopoverContent>
                </Popover>
              </div>
              <div className="space-y-2">
                <Label htmlFor="time-picker" className="px-1">
                  Time *
                </Label>
                <Input
                  type="time"
                  id="time-picker"
                  step="1"
                  disabled={!formData.expired_at}
                  value={
                    formData.expired_at
                      ? new Date(formData.expired_at).toTimeString().split(' ')[0]
                      : '00:00:00'
                  }
                  className="bg-background appearance-none [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
                  onChange={(e) => {
                    const time = e.target.value
                    const date = new Date(formData.expired_at) || new Date()
                    date.setHours(
                      parseInt(time.split(':')[0], 10),
                      parseInt(time.split(':')[1], 10),
                      parseInt(time.split(':')[2] || '0', 10),
                    )
                    handleChange('expired_at', date.toISOString())
                  }}
                />
              </div>
            </div>

            <Separator />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Button type="submit" disabled={loading || success}>
                <Save className="w-4 h-4 mr-2" />
                {loading ? 'Saving...' : coupon_code ? 'Update Coupon' : 'Create Coupon'}
              </Button>
              <Button type="button" variant="outline" onClick={onBack} disabled={loading}>
                Cancel
              </Button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
