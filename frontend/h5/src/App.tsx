import { Routes, Route } from 'react-router-dom'
import PatientListPage from './pages/PatientListPage'
import PatientEditPage from './pages/PatientEditPage'
import RegisterPage from './pages/RegisterPage'
import RegisterConfirmPage from './pages/RegisterConfirmPage'
import RegisterTicketPage from './pages/RegisterTicketPage'

function App() {
  return (
    <Routes>
      <Route path="/h5/patients" element={<PatientListPage />} />
      <Route path="/h5/patients/edit" element={<PatientEditPage />} />
      <Route path="/h5/patients/edit/:id" element={<PatientEditPage />} />
      <Route path="/h5/register" element={<RegisterPage />} />
      <Route path="/h5/register/confirm" element={<RegisterConfirmPage />} />
      <Route path="/h5/register/ticket" element={<RegisterTicketPage />} />
    </Routes>
  )
}

export default App
