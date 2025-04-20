import React, { useEffect, useRef, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import {
    onAuthStateChanged,
    signOut,
    User,
    signInWithPopup,
    GoogleAuthProvider,
} from 'firebase/auth';
import { auth } from '../firebase';
import './Header.css';

const Header: React.FC = () => {
    const [user, setUser] = useState<User | null>(null);
    const [dropdownOpen, setDropdownOpen] = useState(false);
    const [loading, setLoading] = useState(true);
    const dropdownRef = useRef<HTMLDivElement>(null);
    const location = useLocation();
    const navigate = useNavigate();

    useEffect(() => {
        const unsubscribe = onAuthStateChanged(auth, (currUser) => {
            setUser(currUser);
            setLoading(false);
        });
        return () => unsubscribe();
    }, []);

    useEffect(() => {
        const handleClickOutside = (e: MouseEvent) => {
            if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
                setDropdownOpen(false);
            }
        };
        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleLogin = async () => {
        try {
            const provider = new GoogleAuthProvider();
            const result = await signInWithPopup(auth, provider);
            setUser(result.user);
            setDropdownOpen(true);
        } catch (error) {
            console.error("Login failed:", error);
        }
    };

    const handleLogout = async () => {
        await signOut(auth);
        setUser(null);
        setDropdownOpen(false);
        navigate('/');
    };

    return (
        <header className="header">
            <h1 className="logo">CryptoPulse</h1>
            <nav className="nav">
                <div className="nav-links">
                    <Link
                        to="/"
                        className={location.pathname === '/' ? 'active' : ''}
                    >
                        Dashboard
                    </Link>
                    <Link
                        to="/about"
                        className={location.pathname === '/about' ? 'active' : ''}
                    >
                        About
                    </Link>
                </div>

                <div className="auth-section">
                    {!loading && (
                        user ? (
                            <div
                                className="profile-dropdown"
                                onClick={() => setDropdownOpen(prev => !prev)}
                                ref={dropdownRef}
                            >
                                <img
                                    src={user.photoURL || ''}
                                    alt="Profile"
                                    className="profile-pic"
                                />
                                <span className="user-name">{user.displayName}</span>
                                {dropdownOpen && (
                                    <div className="dropdown-menu">
                                        <div className="dropdown-email">{user.email}</div>
                                        <button className="logout-btn" onClick={handleLogout}>
                                            Logout
                                        </button>
                                    </div>
                                )}
                            </div>
                        ) : (
                            <button className="login-link" onClick={handleLogin}>
                                Login
                            </button>
                        )
                    )}
                </div>
            </nav>
        </header>
    );
};

export default Header;
