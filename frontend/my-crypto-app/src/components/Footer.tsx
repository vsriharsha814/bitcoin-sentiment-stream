import React from 'react';
import './Footer.css';

const Footer: React.FC = () => {
    return (
        <footer className={`p-4 text-center border-t ${'border-pink-700 text-gray-400'} ${'bg-gray-900 text-pink-300'}`}>
            <div className="text-xl font-bold tracking-widest mb-2">CRYPTO PULSE</div>
            <div>Your Crypto-Invest Mate</div>
            <div className="mt-2 text-sm">© 2025 Crypto Pulse. All rights reserved.</div>
        </footer>
    );
};

export default Footer;