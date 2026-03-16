import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './styles/main.css'
import "./styles/event.css"
import "./styles/menu.css"
import App from './App.tsx'

createRoot(document.getElementById('main')!).render(
  <StrictMode>
    <App />
  </StrictMode>,
)
