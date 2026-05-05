import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { unstableSetRender } from 'antd-mobile'
import './index.css'
import App from './App.tsx'

unstableSetRender((node, container) => {
  ;(container as any)._reactRoot ||= createRoot(container)
  const root = (container as any)._reactRoot
  root.render(node)
  return async () => {
    await new Promise((resolve) => setTimeout(resolve, 0))
    root.unmount()
  }
})

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </StrictMode>,
)
