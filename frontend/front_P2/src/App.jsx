import { useState } from 'react'
import './App.css'
import Navbar_ from './components/Navbar'
import Login from './components/Login'
import Inicio from './components/Inicio'
import {BrowserRouter, Routes, Route} from 'react-router-dom'

function App() {
  const [count, setCount] = useState(0)

  return (
    <BrowserRouter>
      <Navbar_/>
      <Routes>
        <Route exact path="/" Component={Inicio}/>
        <Route exact path="/login" Component={Login}/>
      </Routes>
    </BrowserRouter>
    
  )
}

export default App
