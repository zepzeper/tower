import { Link } from 'react-router-dom';
import { useState, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router-dom';
import { ChevronDown, User, LogIn } from 'lucide-react';

export const Header = ({
    logo = null,
    ctaText = "Get Started",
}) => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isProductDropdownOpen, setIsProductDropdownOpen] = useState(false);
  const { t, i18n } = useTranslation('common');
  const { lang } = useParams();
  const dropdownRef = useRef(null);

  const navLinks = [
    { 
      name: t('header.product'), 
      path: '/product',
      hasDropdown: true,
      dropdownItems: [
        { name: t('header.productDropdown.features'), path: '/features' },
        { name: t('header.productDropdown.solutions'), path: '/solutions' },
        { name: t('header.productDropdown.useCases'), path: '/use-cases' },
        { name: t('header.productDropdown.documentation'), path: '/docs' }
      ]
    },
    { name: t('header.integrations'), path: '/integrations', hasDropdown: false },
    { name: t('header.pricing'), path: '/pricing', hasDropdown: false },
  ];

  const handleLanguageChange = (newLang) => {
    i18n.changeLanguage(newLang);
    localStorage.setItem('i18nextLng', newLang);
    window.location.href = window.location.href.replace(`/${lang}/`, `/${newLang}/`);
  };

  // Close dropdown when clicking outside
  useEffect(() => {
    function handleClickOutside(event) {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setIsProductDropdownOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  return (
    <header className="bg-white shadow-sm sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-20">
          {/* Logo Section */}
          <div className="flex-shrink-0 flex items-center">
            <Link key={"/"} to={`/${lang}`}>
                {logo || (
                    <div className="flex items-center space-x-2">
                        <svg className="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
                        </svg>
                        <span className="text-2xl font-bold text-gray-900">Tower<span className="text-green-600">API</span></span>
                    </div>
                )}
            </Link>
          </div>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex space-x-8">
            {navLinks.map((link) => ( 
              link.hasDropdown ? (
                <div key={link.path} className="relative" ref={dropdownRef}>
                  <button 
                    className="flex items-center text-gray-600 hover:text-green-600 transition-colors font-medium"
                    onClick={() => setIsProductDropdownOpen(!isProductDropdownOpen)}
                  >
                    {link.name}
                    <ChevronDown className={`ml-1 w-4 h-4 transition-transform ${isProductDropdownOpen ? 'rotate-180' : ''}`} />
                  </button>
                  
                  {isProductDropdownOpen && (
                    <div className="absolute left-0 mt-2 w-56 bg-white rounded-lg shadow-lg py-2 z-10 border border-gray-100">
                      {link.dropdownItems.map((item) => (
                        <Link 
                          key={item.path} 
                          to={`/${lang}${item.path}`} 
                          className="block px-4 py-2 text-gray-700 hover:bg-gray-50 hover:text-green-600"
                          onClick={() => setIsProductDropdownOpen(false)}
                        >
                          {item.name}
                        </Link>
                      ))}
                    </div>
                  )}
                </div>
              ) : (
                <Link 
                  key={link.path} 
                  to={`/${lang}${link.path}`} 
                  className="text-gray-600 hover:text-green-600 transition-colors font-medium"
                >
                  {link.name}
                </Link>
              )
            ))}
          </nav>

          {/* Mobile Menu Button */}
          <div className="md:hidden flex items-center">
            <button
              onClick={() => setIsMenuOpen(!isMenuOpen)}
              className="inline-flex items-center justify-center p-2 rounded-md text-gray-600 hover:text-green-600 focus:outline-none"
            >
              <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                {isMenuOpen ? (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                ) : (
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
                )}
              </svg>
            </button>
          </div>

          {/* Desktop Auth and Language */}
          <div className="hidden md:flex items-center space-x-4">
            {/* Language Switcher */}
            <select 
              value={lang} 
              onChange={(e) => handleLanguageChange(e.target.value)} 
              className="bg-gray-100 rounded-md px-3 py-1.5 text-sm text-gray-700 border border-gray-200 hover:border-gray-300 focus:outline-none focus:ring-1 focus:ring-green-500"
            >
              <option value="en">EN</option>
              <option value="nl">NL</option>
            </select>
            
            {/* Sign in */}
            <Link to={`/${lang}/signin`} className="flex items-center text-gray-600 hover:text-green-600 font-medium">
              <LogIn className="w-4 h-4 mr-1" />
              {t('header.signin')}
            </Link>
            
            {/* Sign up / CTA */}
            <Link to={`/${lang}/signup`} className="bg-green-600 text-white px-5 py-2.5 rounded-lg hover:bg-green-700 transition-colors font-medium shadow-sm">
              {t('header.signup')}
            </Link>
          </div>
        </div>

        {/* Mobile Menu */}
        {isMenuOpen && (
          <div className="md:hidden pb-4 pt-2 border-t border-gray-100">
            <div className="space-y-1">
              {navLinks.map((link) => (
                link.hasDropdown ? (
                  <div key={link.path} className="py-2">
                    <div 
                      className="flex items-center justify-between px-3 py-2 text-gray-700 font-medium"
                      onClick={() => setIsProductDropdownOpen(!isProductDropdownOpen)}
                    >
                      {link.name}
                      <ChevronDown className={`w-4 h-4 transition-transform ${isProductDropdownOpen ? 'rotate-180' : ''}`} />
                    </div>
                    
                    {isProductDropdownOpen && (
                      <div className="mt-1 pl-4 border-l-2 border-gray-100 ml-3">
                        {link.dropdownItems.map((item) => (
                          <Link
                            key={item.path}
                            to={`/${lang}${item.path}`}
                            className="block px-3 py-2 text-gray-600 hover:text-green-600"
                            onClick={() => setIsMenuOpen(false)}
                          >
                            {item.name}
                          </Link>
                        ))}
                      </div>
                    )}
                  </div>
                ) : (
                  <Link
                    key={link.path}
                    to={`/${lang}${link.path}`}
                    className="block px-3 py-2 text-gray-600 hover:text-green-600 hover:bg-gray-50"
                    onClick={() => setIsMenuOpen(false)}
                  >
                    {link.name}
                  </Link>
                )
              ))}
              
              <div className="pt-4 mt-4 border-t border-gray-100">
                <div className="flex items-center px-3 mb-3">
                  <span className="text-sm text-gray-500 mr-2">{t('header.language')}:</span>
                  <select 
                    value={lang} 
                    onChange={(e) => handleLanguageChange(e.target.value)} 
                    className="bg-gray-100 rounded-md px-2 py-1 text-sm text-gray-700 border border-gray-200"
                  >
                    <option value="en">English</option>
                    <option value="nl">Nederlands</option>
                  </select>
                </div>
                
                <Link
                  to={`/${lang}/signin`}
                  className="block px-3 py-2 text-gray-600 hover:text-green-600 hover:bg-gray-50"
                  onClick={() => setIsMenuOpen(false)}
                >
                  <LogIn className="w-4 h-4 inline mr-2" />
                  {t('header.signin')}
                </Link>
                
                <Link
                  to={`/${lang}/signup`}
                  className="block mt-2 bg-green-600 text-white px-6 py-2.5 rounded-lg hover:bg-green-700 transition-colors font-medium"
                  onClick={() => setIsMenuOpen(false)}
                >
                  {t('header.signup')}
                </Link>
              </div>
            </div>
          </div>
        )}
      </div>
    </header>
  );
};
