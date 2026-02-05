import React from "react";
import { Link, NavLink } from 'react-router-dom';
import './NavBar.scss';

const NavBar: React.FC = () => {
    return (
        <nav className="navbar navbar-expand-lg navbar-dark bg-primary">
            <div className="container">
                <div className="collapse navbar-collapse" id="navbarNav">
                    <ul className="navbar-nav">
                        <li className="nav-item">
                            <NavLink
                                className={({ isActive }) => isActive ? "nav-link active" : "nav-link"}
                                to="/"
                            >
                                Home
                            </NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink
                                className={({ isActive }) => isActive ? "nav-link active" : "nav-link"}
                                to="/problems"
                            >
                                Problems
                            </NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink
                                className={({ isActive }) => isActive ? "nav-link active" : "nav-link"}
                                to="/problem/new"
                            >
                                Create Problem
                            </NavLink>
                        </li>
                        <li className="nav-item">
                            <NavLink
                                className={({ isActive }) => isActive ? "nav-link active" : "nav-link"}
                                to="/quiz"
                            >
                                Take Quiz
                            </NavLink>
                        </li>
                    </ul>
                </div>
            </div>
        </nav>
    );
};

export default NavBar;