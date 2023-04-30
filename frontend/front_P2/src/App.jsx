import { useState } from 'react'
import './App.css'
import Navbar_ from './components/Navbar'
import Login from './components/Login'
import Inicio from './components/Inicio'
import Reportes from './components/Reportes'
import {BrowserRouter, Routes, Route} from 'react-router-dom'

function App() {
  const [count, setCount] = useState(0)

  return (
    <BrowserRouter>
      <Navbar_/>
      <Routes>
        <Route exact path="/" Component={Inicio}/>
        <Route exact path="/login" Component={Login}/>
        <Route exact path="/reports" Component={Reportes}/>
      </Routes>
    </BrowserRouter>
    
  )
}

export default App
