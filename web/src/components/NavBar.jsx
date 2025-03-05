import logo from '../assets/logo.png'
import './NavBar.css'
import { GRAFANA_URL } from './Constants'

const NavBar = () => {
    return (
        <div className="navbar">
            <img className="logo" src={logo} alt="Logo" />
            <ul className='nav-links'>
                <li><a href="https://blog.ak0.io" className="nav-link">blog</a></li>
                <li><a href="https://blog.ak0.io" className="nav-link">contact</a></li>
                <li><a href="/resume_Alex_Krenitsky.pdf" className="nav-link">resume</a></li>
                <li><a href={GRAFANA_URL} className="nav-link">site metrics</a></li>
            </ul>
        </div>
    )
}

export default NavBar
