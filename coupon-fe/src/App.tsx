import CouponManagement from "./pages/coupon-page"



function App() {
  console.log("meta ", import.meta.env.REACT_APP_API_BASE_URL)
  return (
    <>
      <CouponManagement />
    </>
  )
}

export default App
