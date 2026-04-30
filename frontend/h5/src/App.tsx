import { Routes, Route, Navigate, useLocation } from 'react-router-dom'
import BottomNav from './components/BottomNav'
import PatientListPage from './pages/PatientListPage'
import PatientEditPage from './pages/PatientEditPage'
import RegisterPage from './pages/RegisterPage'
import RegisterConfirmPage from './pages/RegisterConfirmPage'
import RegisterTicketPage from './pages/RegisterTicketPage'
import OrderListPage from './pages/OrderListPage'
import OrderDetailPage from './pages/OrderDetailPage'
import PersonalCenterPage from './pages/PersonalCenterPage'

function App() {
  const location = useLocation()

  // Show bottom nav on main tab pages
  const showBottomNav = [
    '/h5/register',
    '/h5/me',
    '/h5/patients',
    '/h5/orders',
  ].includes(location.pathname)

  const wrap = (page: React.ReactNode) =>
    showBottomNav ? <BottomNav>{page}</BottomNav> : <>{page}</>

  return (
    <Routes>
      <Route path="/" element={<Navigate to="/h5/register" replace />} />
      <Route path="/h5/register" element={wrap(<RegisterPage />)} />
      <Route path="/h5/register/confirm" element={<RegisterConfirmPage />} />
      <Route path="/h5/register/ticket" element={<RegisterTicketPage />} />
      <Route path="/h5/patients" element={wrap(<PatientListPage />)} />
      <Route path="/h5/patients/edit" element={<PatientEditPage />} />
      <Route path="/h5/patients/edit/:id" element={<PatientEditPage />} />
      <Route path="/h5/orders" element={wrap(<OrderListPage />)} />
      <Route path="/h5/orders/:id" element={<OrderDetailPage />} />
      <Route path="/h5/me" element={wrap(<PersonalCenterPage />)} />
    </Routes>
  )
}

export default App
