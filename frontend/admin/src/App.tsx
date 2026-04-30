import { Routes, Route } from 'react-router-dom'
import BasicLayout from './components/BasicLayout'
import PatientListPage from './pages/PatientListPage'
import PatientDetailPage from './pages/PatientDetailPage'
import ScheduleListPage from './pages/ScheduleListPage'
import ScheduleEditPage from './pages/ScheduleEditPage'
import OrderListPage from './pages/OrderListPage'
import OrderDetailPage from './pages/OrderDetailPage'

function App() {
  return (
    <BasicLayout>
      <Routes>
        <Route path="/patients" element={<PatientListPage />} />
        <Route path="/patients/:id" element={<PatientDetailPage />} />
        <Route path="/schedules" element={<ScheduleListPage />} />
        <Route path="/schedules/edit" element={<ScheduleEditPage />} />
        <Route path="/orders" element={<OrderListPage />} />
        <Route path="/orders/:id" element={<OrderDetailPage />} />
      </Routes>
    </BasicLayout>
  )
}

export default App
