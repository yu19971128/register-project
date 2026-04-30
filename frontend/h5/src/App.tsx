import { Routes, Route } from 'react-router-dom'
import PatientListPage from './pages/PatientListPage'
import PatientEditPage from './pages/PatientEditPage'

function App() {
  return (
    <Routes>
      <Route path="/h5/patients" element={<PatientListPage />} />
      <Route path="/h5/patients/edit" element={<PatientEditPage />} />
      <Route path="/h5/patients/edit/:id" element={<PatientEditPage />} />
    </Routes>
  )
}

export default App
