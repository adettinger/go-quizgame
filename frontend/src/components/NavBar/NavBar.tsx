import React from "react";
import { Link, NavLink } from 'react-router-dom';
import './NavBar.scss';

interface NavItem {
    path: string;
    name: string;
}

const NavBar: React.FC = () => {
    const navItems: NavItem[] = [
        { path: '/', name: 'Home' },
        { path: '/problems', name: 'Problems' },
        { path: '/problem/new', name: 'Create Problem' },
        { path: '/quiz', name: 'Take Quiz' }
    ];
    return (
        <nav className="navbar navbar-expand-lg navbar-dark bg-primary">
            <div className="container">
                <div className="collapse navbar-collapse" id="navbarNav">
                    <ul className="navbar-nav">
                        {navItems.map((item, index) => (
                            <li className="nav-item" key={index}>
                                <NavLink
                                    className={({ isActive }) => isActive ? "nav-link active" : "nav-link"}
                                    to={item.path}
                                >
                                    {item.name}
                                </NavLink>
                            </li>
                        ))}
                    </ul>
                </div>
            </div>
        </nav>
    );
};

export default NavBar;