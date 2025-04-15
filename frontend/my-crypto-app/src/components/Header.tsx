import React from 'react';
import { Link } from 'react-router-dom';
import './Header.css';

const Header: React.FC = () => {
    return (
        <header className="header">
            <h1 className="logo">CryptoPulse</h1>
            <nav className="nav">
                <Link to="/">Historical</Link>
                <Link to="/live">Live</Link>
                <Link to="/about">About</Link>
            </nav>
        </header>
    );
};

export default Header;
