import { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import i18n from '../i18n';

const LanguageWrapper = (): null => {
  const navigate = useNavigate();
  const location = useLocation();
  const supportedLangs = ['en', 'nl'] as const;
  
  useEffect(() => {
    const pathLang = location.pathname.split('/')[1];
    const storedLang = localStorage.getItem('i18nextLng');
    const browserLang = navigator.language.split('-')[0];

    if (supportedLangs.includes(pathLang as typeof supportedLangs[number])) {
      i18n.changeLanguage(pathLang);
      return;
    }
    
    // Determine preferred language
    const preferredLang = supportedLangs.includes(storedLang as typeof supportedLangs[number]) 
      ? storedLang 
      : supportedLangs.includes(browserLang as typeof supportedLangs[number])
        ? browserLang 
        : 'en';
    
    // Redirect to correct language path
    navigate(`/${preferredLang}${location.pathname}`, { replace: true });
  }, [location.pathname, navigate]);

  return null;
};

export default LanguageWrapper;
