import { Routes, Route } from 'react-router-dom'
import BasicLayout from './components/BasicLayout'
import PatientListPage from './pages/PatientListPage'
import PatientDetailPage from './pages/PatientDetailPage'

function App() {
  return (
    <BasicLayout>
      <Routes>
        <Route path="/patients" element={<PatientListPage />} />
        <Route path="/patients/:id" element={<PatientDetailPage />} />
      </Routes>
    </BasicLayout>
  )
}

export default App
