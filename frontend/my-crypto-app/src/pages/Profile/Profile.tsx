import { useState } from 'react';
import { Bell, Edit} from 'lucide-react';

export default function CryptoProfilePage() {
  const [darkMode] = useState(true);
  const [activeTab, setActiveTab] = useState(1);
  
  // Mock data - in a real app this would come from props or context
  const userData = {
    name: "Satoshi Nakamoto",
    email: "satoshi@example.com",
    profilePic: "/api/placeholder/150/150",
    coins: [
      { id: 1, code: 'BTC'},
      { id: 2, code: 'ETH'},
      { id: 3, code: 'USDT'},
      { id: 4, code: 'XRP'},
      { id: 5, code: 'BNB'},
      { id: 6, code: 'SOL'},
      { id: 7, code: 'USDC'},
      { id: 8, code: 'TRX'},
      { id: 9, code: 'DOGE'},
      { id: 10, code: 'ADA'},
    ],
    questions: [
      { id: 1, question: "What's your take on the next BTC halving?", answers: 23 },
      { id: 2, question: "Is ETH 2.0 the future of DeFi?", answers: 47 }
    ],
    alerts: [
      { id: 1, coin: 'BTC', condition: 'above', price: 60000, active: true },
      { id: 2, coin: 'ETH', condition: 'below', price: 3000, active: false }
    ]
  };

  return (
    <div className={`min-h-screen ${darkMode ? 'bg-gray-900 text-pink-300' : 'bg-gray-200 text-pink-600'}`}>
      <div className="container mx-auto px-4 py-8">
        {/* Profile Section */}
        <div className={`mb-8 p-6 rounded-lg grid grid-cols-1 md:grid-cols-3 gap-6 ${'bg-gray-800'}`}>
          <div className="flex flex-col items-center justify-center">
            <div className={`relative w-40 h-40 rounded-full overflow-hidden border-4 ${'border-pink-500'}`}>
              <img src={userData.profilePic} alt="Profile" className="w-full h-full object-cover" />
            </div>
          </div>
          
          <div className="col-span-2">
            <div className={`text-2xl font-bold mb-2 ${darkMode ? 'text-pink-300' : 'text-pink-600'}`}>
              {userData.name}
            </div>
            <div className={`mb-4 ${darkMode ? 'text-gray-400' : 'text-gray-600'}`}>
              {userData.email}
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div className={`p-4 rounded ${darkMode ? 'bg-gray-700' : 'bg-gray-200'}`}>
                <div className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-600'}`}>
                  Last Login
                </div>
                <div className={`text-xl font-bold ${darkMode ? 'text-pink-300' : 'text-pink-600'}`}>
                  Apr 26 2025
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Selected Coins */}
        <div className={`mb-8 p-6 rounded-lg ${darkMode ? 'bg-gray-800' : 'bg-gray-100'}`}>
          <h2 className={`text-xl font-bold mb-4 ${darkMode ? 'text-pink-300' : 'text-pink-600'}`}>
            SELECTED COINS
          </h2>
          
          {/* Tabs */}
          <div className="flex flex-wrap gap-2 mb-4">
            {userData.coins.map(coin => (
              <button
                key={coin.id}
                onClick={() => setActiveTab(coin.id)}
                className={`px-4 py-2 rounded-t border-2 
                  ${activeTab === coin.id ? 
                    (darkMode ? 'bg-pink-900 border-pink-500 text-pink-300' : 'bg-pink-200 border-pink-500 text-pink-700') : 
                    (darkMode ? 'bg-gray-700 border-gray-600 text-gray-400' : 'bg-gray-200 border-gray-300 text-gray-600')
                  }`}
              >
                {coin.code}
              </button>
            ))}
          </div>
        </div>

        {/* User Questions */}
        <div className={`mb-8 p-6 rounded-lg ${darkMode ? 'bg-gray-800' : 'bg-gray-100'}`}>
          <h2 className={`text-xl font-bold mb-4 ${darkMode ? 'text-pink-300' : 'text-pink-600'}`}>
            YOUR QUESTIONS
          </h2>
          
          <div className="space-y-4">
            {userData.questions.map(question => (
              <div 
                key={question.id} 
                className={`p-4 rounded border-l-4 ${darkMode ? 'bg-gray-700 border-pink-500' : 'bg-gray-50 border-pink-400'}`}
              >
                <div className={`text-lg font-bold ${darkMode ? 'text-pink-300' : 'text-pink-600'}`}>
                  {question.question}
                </div>
                <div className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-600'}`}>
                  {question.answers} answers
                </div>
              </div>
            ))}
            
            <button className={`w-full p-3 mt-3 border-2 border-dashed flex items-center justify-center 
              ${darkMode ? 'border-gray-600 text-gray-400 hover:border-pink-500 hover:text-pink-300' : 
                'border-gray-300 text-gray-500 hover:border-pink-400 hover:text-pink-600'}`}
            >
              + ASK NEW QUESTION
            </button>
          </div>
        </div>

        {/* Price Alerts */}
        <div className={`p-6 rounded-lg ${darkMode ? 'bg-gray-800' : 'bg-gray-100'}`}>
          <h2 className={`text-xl font-bold mb-4 flex items-center gap-2 ${darkMode ? 'text-pink-300' : 'text-pink-600'}`}>
            <Bell />
            PRICE ALERTS
          </h2>
          
          <div className="space-y-4">
            {userData.alerts.map(alert => (
              <div 
                key={alert.id} 
                className={`p-4 rounded flex justify-between items-center 
                  ${darkMode ? 
                    (alert.active ? 'bg-gray-700' : 'bg-gray-700 opacity-50') : 
                    (alert.active ? 'bg-gray-50' : 'bg-gray-50 opacity-50')}`}
              >
                <div>
                  <div className={`text-lg font-bold ${darkMode ? 'text-pink-300' : 'text-pink-600'}`}>
                    {alert.coin} {alert.condition === 'above' ? '↑' : '↓'} ${alert.price.toLocaleString()}
                  </div>
                  <div className={`text-sm ${darkMode ? 'text-gray-400' : 'text-gray-600'}`}>
                    Alert when {alert.coin} goes {alert.condition} ${alert.price.toLocaleString()}
                  </div>
                </div>
                <div className="flex gap-2">
                  <button className={`p-2 rounded ${darkMode ? 'bg-gray-600 hover:bg-gray-500' : 'bg-gray-200 hover:bg-gray-300'}`}>
                    <Edit size={16} className={darkMode ? 'text-gray-300' : 'text-gray-600'} />
                  </button>
                  <div className="relative inline-block w-12 align-middle select-none">
                    <input 
                      type="checkbox" 
                      id={`toggle-${alert.id}`}
                      name={`toggle-${alert.id}`} 
                      className="sr-only"
                      checked={alert.active}
                      readOnly
                    />
                    <div className={`block h-6 rounded-full ${alert.active ? 'bg-pink-500' : 'bg-gray-600'}`}></div>
                    <div className={`dot absolute left-1 top-1 h-4 w-4 rounded-full transition ${alert.active ? 'transform translate-x-6 bg-white' : 'bg-gray-300'}`}></div>
                  </div>
                </div>
              </div>
            ))}
            
            <button className={`w-full p-3 mt-3 border-2 border-dashed flex items-center justify-center 
              ${darkMode ? 'border-gray-600 text-gray-400 hover:border-pink-500 hover:text-pink-300' : 
                'border-gray-300 text-gray-500 hover:border-pink-400 hover:text-pink-600'}`}
            >
              + CREATE NEW ALERT
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}