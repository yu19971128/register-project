import { Routes, Route } from 'react-router-dom'
import BasicLayout from './components/BasicLayout'
import PatientListPage from './pages/PatientListPage'
import PatientDetailPage from './pages/PatientDetailPage'
import ScheduleListPage from './pages/ScheduleListPage'
import ScheduleEditPage from './pages/ScheduleEditPage'

function App() {
  return (
    <BasicLayout>
      <Routes>
        <Route path="/patients" element={<PatientListPage />} />
        <Route path="/patients/:id" element={<PatientDetailPage />} />
        <Route path="/schedules" element={<ScheduleListPage />} />
        <Route path="/schedules/edit" element={<ScheduleEditPage />} />
      </Routes>
    </BasicLayout>
  )
}

export default App
