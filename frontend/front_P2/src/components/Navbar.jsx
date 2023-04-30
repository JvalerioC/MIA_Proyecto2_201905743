import React, { useState } from 'react';
import {
  Navbar,
  NavbarBrand,
  Nav,
  NavItem,
  NavLink,
  NavbarText,
} from 'reactstrap';

function Navbar_() {
  const [isOpen, setIsOpen] = useState(false);


  return (
    <div>
      <Navbar className='navbar navbar-expand-lg navbar-light bg-light'>

        <Nav className="me-auto" navbar>
          <NavLink href="/" className='font-weight-bold' style={{fontSize: "18px", color: "black"}}>
            <img src='/th.webp' width={30} />Inicio
          </NavLink>
          <NavItem>
            <NavLink href="/login">Login</NavLink>
          </NavItem>
          <NavItem>
            <NavLink href="/">
              Cargar Archivo
            </NavLink>
          </NavItem>
          <NavItem>
            <NavLink href="/reports">
              Reportes
            </NavLink>
          </NavItem>
        </Nav>
      </Navbar>
    </div>
  );
  /* href en navlink y to en link */
}

export default Navbar_;