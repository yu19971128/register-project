import { Routes, Route } from 'react-router-dom'
import PatientListPage from './pages/PatientListPage'
import PatientEditPage from './pages/PatientEditPage'
import RegisterPage from './pages/RegisterPage'
import RegisterConfirmPage from './pages/RegisterConfirmPage'
import RegisterTicketPage from './pages/RegisterTicketPage'
import OrderListPage from './pages/OrderListPage'
import OrderDetailPage from './pages/OrderDetailPage'

function App() {
  return (
    <Routes>
      <Route path="/h5/patients" element={<PatientListPage />} />
      <Route path="/h5/patients/edit" element={<PatientEditPage />} />
      <Route path="/h5/patients/edit/:id" element={<PatientEditPage />} />
      <Route path="/h5/register" element={<RegisterPage />} />
      <Route path="/h5/register/confirm" element={<RegisterConfirmPage />} />
      <Route path="/h5/register/ticket" element={<RegisterTicketPage />} />
      <Route path="/h5/orders" element={<OrderListPage />} />
      <Route path="/h5/orders/:id" element={<OrderDetailPage />} />
    </Routes>
  )
}

export default App
