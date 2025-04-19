import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Bell, Search, Menu, LogOut, Sun, Moon } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import { useTheme } from '../context/ThemeContext';

type Language = 'en' | 'nl';

const TopBar = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const { theme, toggleTheme } = useTheme();
  const [showDropdown, setShowDropdown] = useState(false);
  const { t, i18n } = useTranslation('pages');

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  const handleLanguageChange = (lang: Language) => {
    i18n.changeLanguage(lang);
    setShowDropdown(false);
  };

  return (
    <header className={`${theme === 'light' ? 'bg-white border-b border-gray-200' : 'bg-gray-900 border-b border-gray-700'} h-16 flex items-center px-6 transition-colors duration-200`}>
      <button className="lg:hidden mr-4">
        <Menu className={`h-6 w-6 ${theme === 'light' ? 'text-gray-500' : 'text-gray-400'}`} />
      </button>
      
      <div className="relative flex-1 max-w-xs">
        <div className="absolute inset-y-0 left-0 flex items-center pl-3">
          <Search className={`h-5 w-5 ${theme === 'light' ? 'text-gray-400' : 'text-gray-500'}`} />
        </div>
        <input
          type="text"
          placeholder={t('topBar.search')}
          className={`block w-full rounded-md border-0 py-2 pl-10 ${
            theme === 'light' 
              ? 'text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-green-600 bg-white' 
              : 'text-white ring-1 ring-inset ring-gray-700 placeholder:text-gray-500 focus:ring-2 focus:ring-inset focus:ring-green-500 bg-gray-800'
          }`}
        />
      </div>
      
      <div className="ml-auto flex items-center space-x-4">
        {/* Language Selector */}
        <div className="relative">
          <button 
            onClick={() => setShowDropdown(!showDropdown)} 
            className={`px-2 py-1 text-sm font-medium rounded-md ${theme === 'light' ? 'text-gray-600 hover:bg-gray-100' : 'text-gray-300 hover:bg-gray-800'}`}
          >
            {i18n.language === 'nl' ? 'NL' : 'EN'}
          </button>
          
          {showDropdown && (
            <div className={`absolute right-0 mt-2 w-40 ${
              theme === 'light' 
                ? 'bg-white ring-1 ring-black ring-opacity-5' 
                : 'bg-gray-800 ring-1 ring-gray-700'
              } rounded-md shadow-lg py-1 z-10`}
            >
              <button
                onClick={() => handleLanguageChange('en')}
                className={`flex items-center w-full px-4 py-2 text-sm ${
                  theme === 'light' 
                    ? 'text-gray-700 hover:bg-gray-100' 
                    : 'text-gray-300 hover:bg-gray-700'
                } ${i18n.language === 'en' ? 'font-bold' : ''}`}
              >
                English
              </button>
              <button
                onClick={() => handleLanguageChange('nl')}
                className={`flex items-center w-full px-4 py-2 text-sm ${
                  theme === 'light' 
                    ? 'text-gray-700 hover:bg-gray-100' 
                    : 'text-gray-300 hover:bg-gray-700'
                } ${i18n.language === 'nl' ? 'font-bold' : ''}`}
              >
                Nederlands
              </button>
            </div>
          )}
        </div>
        
        {/* Theme Toggle Button */}
        <button 
          onClick={toggleTheme} 
          className={`p-1 rounded-md ${theme === 'light' ? 'text-gray-500 hover:bg-gray-100' : 'text-gray-400 hover:bg-gray-800'}`}
          aria-label={t('topBar.toggleDarkMode')}
        >
          {theme === 'light' ? <Moon className="h-5 w-5" /> : <Sun className="h-5 w-5" />}
        </button>
        
        <button className={`relative p-1 ${theme === 'light' ? 'text-gray-500 hover:text-gray-700' : 'text-gray-400 hover:text-gray-300'}`} aria-label={t('topBar.notifications')}>
          <Bell className="h-6 w-6" />
          <span className="absolute top-0 right-0 h-2 w-2 rounded-full bg-red-500"></span>
        </button>
        
        <div className="relative">
          <button
            onClick={() => setShowDropdown(!showDropdown)}
            className="flex items-center space-x-2"
          >
            <div className={`h-8 w-8 rounded-full ${theme === 'light' ? 'bg-green-100 text-green-600' : 'bg-green-800 text-green-100'} flex items-center justify-center`}>
              <span className="font-semibold">{user?.name?.charAt(0) || 'U'}</span>
            </div>
            <span className={`text-sm font-medium ${theme === 'light' ? 'text-gray-700' : 'text-gray-300'} hidden md:inline-block`}>
              {user?.name || 'User'}
            </span>
          </button>
          
          {showDropdown && (
            <div className={`absolute right-0 mt-2 w-48 ${
              theme === 'light' 
                ? 'bg-white ring-1 ring-black ring-opacity-5' 
                : 'bg-gray-800 ring-1 ring-gray-700'
              } rounded-md shadow-lg py-1 z-10`}
            >
              <button
                onClick={handleLogout}
                className={`flex items-center w-full px-4 py-2 text-sm ${
                  theme === 'light' 
                    ? 'text-gray-700 hover:bg-gray-100' 
                    : 'text-gray-300 hover:bg-gray-700'
                }`}
              >
                <LogOut className="h-4 w-4 mr-2" />
                {t('navigation.signOut')}
              </button>
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default TopBar;
