import CouponDetails from '@/components/coupon-details'
import CouponForm from '@/components/coupon-form'
import CouponList from '@/components/coupon-list'
import MockOrder from '@/components/mock-order'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import type { Coupon } from '@/types/model/coupon'
import { Calculator, List, Plus } from 'lucide-react'
import { useState } from 'react'
import zalopayLogo from '@/assets/Logo FA-09.png'

export default function CouponManagement() {
  const [currentView, setCurrentView] = useState('list')
  const [selectedCoupon, setSelectedCoupon] = useState<string | null>(null)
  const [editingCoupon, setEditingCoupon] = useState<string | null>(null)

  const handleViewCoupon = (coupon_code: string) => {
    setSelectedCoupon(coupon_code)
    setCurrentView('details')
  }

  const handleEditCoupon = (coupon_code: string) => {
    setEditingCoupon(coupon_code)
    setCurrentView('form')
  }

  const handleCreateNew = () => {
    setEditingCoupon(null)
    setCurrentView('form')
  }

  const handleBackToList = () => {
    setCurrentView('list')
    setSelectedCoupon(null)
    setEditingCoupon(null)
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b bg-card">
        <div className="container mx-auto px-4 py-6">
          <img src={zalopayLogo} className="w-3xs h-auto" />
          <h1 className="text-3xl font-bold text-foreground"> Coupon Management System</h1>
          <p className="text-muted-foreground mt-2">Manage your discount coupons and test orders</p>
        </div>
      </header>

      {/* Navigation */}
      <div className="border-b bg-muted/50">
        <div className="container mx-auto px-4 py-4">
          <Tabs value={currentView} onValueChange={setCurrentView} className="w-full">
            <TabsList className="grid w-full max-w-md grid-cols-3">
              <TabsTrigger value="list" className="flex items-center gap-2">
                <List className="w-4 h-4" />
                Coupons
              </TabsTrigger>
              <TabsTrigger
                value="form"
                className="flex items-center gap-2"
                onClick={handleCreateNew}
              >
                <Plus className="w-4 h-4" />
                Create
              </TabsTrigger>
              <TabsTrigger value="mock-order" className="flex items-center gap-2">
                <Calculator className="w-4 h-4" />
                Test Orders
              </TabsTrigger>
            </TabsList>
          </Tabs>
        </div>
      </div>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        {currentView === 'list' && (
          <CouponList onViewCoupon={handleViewCoupon} onEditCoupon={handleEditCoupon} />
        )}

        {currentView === 'form' && <CouponForm coupon_code={editingCoupon} onBack={handleBackToList} />}

        {currentView === 'details' && selectedCoupon && (
          <CouponDetails
            coupon_code={selectedCoupon}
            onBack={handleBackToList}
            onEdit={() => handleEditCoupon(selectedCoupon)}
          />
        )}

        {currentView === 'mock-order' && <MockOrder />}
      </main>
    </div>
  )
}
