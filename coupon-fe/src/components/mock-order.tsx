import { orderService } from '@/api/orders'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import type { MockOrder as MockOrderModel } from '@/types/model/order'
import { formatDate, formatDateTime } from '@/utils/time'
import {
  AlertCircle,
  Calculator,
  Calendar as CalendarLogo,
  CheckCircle,
  ChevronDownIcon,
  Receipt,
} from 'lucide-react'
import type React from 'react'
import { useState } from 'react'
import { Badge } from './ui/badge'
import { Calendar } from './ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover'

export default function MockOrder() {
  const [result, setResult] = useState<MockOrderModel | null>(null)
  const [loading, setLoading] = useState(false)
  const [open, setOpen] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [formData, setFormData] = useState({
    cost: 0,
    coupon_code: null,
    created_at: new Date().toISOString(),
    coupon: null,
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setResult(null)

    try {
      const data = await orderService.mockOrder(formData as Omit<MockOrderModel, 'total_amount'>)
      if ('error' in data) {
        setError(data.error)
      } else {
        setResult(data.data)
        setError(null)
      }
    } catch (error) {
      console.error('Error testing order:', error)
      setError('Failed to calculate order: ' + (error as Error).message)
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (field: string, value: string | number | null) => {
    setFormData((prev) => ({
      ...prev,
      [field]: value,
    }))
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold tracking-tight">Test Order Calculator</h2>
        <p className="text-muted-foreground">
          Test how coupons apply to orders with different amounts and dates
        </p>
      </div>

      <Separator />

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Input Form */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <Calculator className="w-5 h-5 mr-2" />
              Order Details
            </CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit} className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="cost">Order Amount (VND) *</Label>
                <Input
                  id="cost"
                  type="number"
                  step="1000"
                  min="0"
                  value={formData.cost}
                  onChange={(e) => handleChange('cost', parseInt(e.target.value))}
                  placeholder="10000"
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="coupon_code">Coupon Code</Label>
                <Input
                  id="coupon_code"
                  value={formData.coupon_code || ''}
                  onChange={(e) => handleChange('coupon_code', e.target.value || null)}
                  placeholder="ZLP8k"
                />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-2">
                  <Label htmlFor="date-picker" className="px-1">
                    Order Date *
                  </Label>
                  <Popover open={open} onOpenChange={setOpen}>
                    <PopoverTrigger asChild>
                      <Button
                        variant="outline"
                        id="date-picker"
                        className="w-full justify-between font-normal"
                      >
                        {formData.created_at ? formatDate(formData.created_at) : 'Select date'}
                        <ChevronDownIcon />
                      </Button>
                    </PopoverTrigger>
                    <PopoverContent className="w-auto overflow-hidden p-0" align="start">
                      <Calendar
                        mode="single"
                        selected={formData.created_at ? new Date(formData.created_at) : undefined}
                        captionLayout="dropdown"
                        onSelect={(date) => {
                          handleChange('created_at', date ? date.toISOString() : '')
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
                    value={
                      formData.created_at
                        ? new Date(formData.created_at).toTimeString().split(' ')[0]
                        : '10:30:00'
                    }
                    className="bg-background appearance-none [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
                    onChange={(e) => {
                      const time = e.target.value
                      const date = new Date(formData.created_at)
                      date.setHours(
                        parseInt(time.split(':')[0], 10),
                        parseInt(time.split(':')[1], 10),
                        parseInt(time.split(':')[2] || '0', 10),
                      )
                      handleChange('created_at', date.toISOString())
                    }}
                  />
                </div>
              </div>

              <Button type="submit" disabled={loading} className="w-full">
                <Calculator className="w-4 h-4 mr-2" />
                {loading ? 'Calculating...' : 'Calculate Discount'}
              </Button>
            </form>
          </CardContent>
        </Card>

        {/* Results */}
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center">
              <Receipt className="w-5 h-5 mr-2" />
              Calculation Results
            </CardTitle>
          </CardHeader>
          <CardContent>
            {!result && !error ? (
              <div className="flex flex-col items-center justify-center py-12 text-center">
                <Calculator className="w-12 h-12 text-muted-foreground mb-4" />
                <p className="text-muted-foreground">
                  Enter order details and click calculate to see results
                </p>
              </div>
            ) : (
              <div className="space-y-6">
                {/* Status Alert */}
                {!error ? (
                  <Alert>
                    <CheckCircle className="h-4 w-4" />
                    <AlertDescription>Coupon Applied Successfully</AlertDescription>
                  </Alert>
                ) : (
                  <Alert variant="destructive">
                    <AlertCircle className="h-4 w-4" />
                    <AlertDescription>{error || 'Coupon Not Applied'}</AlertDescription>
                  </Alert>
                )}
              </div>
            )}
            {result && (
              <div className="space-y-3">
                <div className="bg-gray-50 p-4 rounded-md space-y-3">
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600">Original Amount</span>
                    <span className="font-semibold">{result.cost} VND</span>
                  </div>

                  <div className="flex justify-between items-center">
                    <span className="text-gray-600">Discount Amount</span>
                    <span className="font-semibold text-green-600">
                      {result.total_amount - result.cost} VND
                    </span>
                  </div>

                  <hr className="border-gray-200" />

                  <div className="flex justify-between items-center text-lg">
                    <span className="font-semibold text-gray-900">Final Amount</span>
                    <span className="font-bold text-blue-600">{result.total_amount} VND</span>
                  </div>
                </div>

                {/* Coupon Details */}
                {result.coupon && (
                  <div className="bg-blue-50 p-4 rounded-md">
                    <h4 className="font-medium text-blue-900 mb-2">Applied Coupon</h4>
                    <Separator className="my-2" />
                    <div className="space-y-2 text-sm">
                      <div className="flex justify-between">
                        <span className="text-blue-700">Coupon Code</span>
                        <Badge variant="secondary" className="bg-blue-100 text-blue-800 font-mono">
                          {result.coupon.coupon_code}
                        </Badge>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-blue-700">Title</span>
                        <span className="text-blue-900">{result.coupon.title}</span>
                      </div>
                      <div className="flex justify-between">
                        <span className="text-blue-700">Coupon Type</span>
                        <div className="flex items-center text-blue-900">
                          {result.coupon.coupon_type === 'percentage' ? (
                            <>{result.coupon.coupon_value}%</>
                          ) : (
                            <>{result.coupon.coupon_value}</>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                {/* Expiry Info */}
                {result.coupon?.expired_at && (
                  <div className="bg-yellow-50 p-3 rounded-md">
                    <div className="flex items-center text-sm text-yellow-800">
                      <CalendarLogo className="w-4 h-4 mr-2" />
                      Coupon expired on {formatDateTime(result.coupon.expired_at)}
                    </div>
                  </div>
                )}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
