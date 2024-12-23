import logo from '../assets/logo.png'
import './NavBar.css'

const NavBar = () => {
    return (
        <div className="navbar">
            <img className="logo" src={logo} alt="Logo" />
            <ul className='nav-links'>
                <li><a href="https://blog.ak0.io" className="nav-link">blog</a></li>
                <li><a href="https://blog.ak0.io" className="nav-link">contact</a></li>
                <li><a href="https://ak0.io/static/resume_Alex_Krenitsky.pdf" className="nav-link">resume</a></li>
            </ul>
        </div>
    )
}

export default NavBar
