import logo from '../assets/logo.png'
import './NavBar.css'

const NavBar = () => {
    return (
        <div className="navbar">
            <img className="logo" src={logo} alt="Logo" />
            <ul className='nav-links'>
                <li><a href="https://blog.ak0.io" className="nav-link">blog</a></li>
                <li><a href="https://blog.ak0.io" className="nav-link">contact</a></li>
                <li><a href="/resume_Alex_Krenitsky.pdf" className="nav-link">resume</a></li>
                <li><a href="/grafana/d/eeer2i27cwm4gb/ak0-overview%3ForgId%3D1%26from%3Dnow-24h%26to%3Dnow%26timezone%3Dbrowser" className="nav-link">site metrics</a></li>
            </ul>
        </div>
    )
}

export default NavBar
