import { Routes, Route, Navigate } from 'react-router-dom'
import BasicLayout from './components/BasicLayout'
import LoginPage from './pages/LoginPage'
import PatientListPage from './pages/PatientListPage'
import PatientDetailPage from './pages/PatientDetailPage'
import ScheduleListPage from './pages/ScheduleListPage'
import ScheduleEditPage from './pages/ScheduleEditPage'
import OrderListPage from './pages/OrderListPage'
import OrderDetailPage from './pages/OrderDetailPage'

function isLoggedIn() {
  return !!localStorage.getItem('admin_token')
}

function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route
        path="/*"
        element={
          isLoggedIn() ? (
            <BasicLayout>
              <Routes>
                <Route path="/patients" element={<PatientListPage />} />
                <Route path="/patients/:id" element={<PatientDetailPage />} />
                <Route path="/schedules" element={<ScheduleListPage />} />
                <Route path="/schedules/edit" element={<ScheduleEditPage />} />
                <Route path="/orders" element={<OrderListPage />} />
                <Route path="/orders/:id" element={<OrderDetailPage />} />
                <Route path="/" element={<Navigate to="/patients" replace />} />
              </Routes>
            </BasicLayout>
          ) : (
            <Navigate to="/login" replace />
          )
        }
      />
    </Routes>
  )
}

export default App
