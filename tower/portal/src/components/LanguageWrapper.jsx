import { useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import i18n from '../i18n';
const LanguageWrapper = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const supportedLangs = ['en', 'nl'];
  useEffect(() => {
    const pathLang = location.pathname.split('/')[1];
    const storedLang = localStorage.getItem('i18nextLng');
    const browserLang = navigator.language.split('-')[0];

    if (supportedLangs.includes(pathLang)) {
      i18n.changeLanguage(pathLang);
      return;
    }
    // Determine preferred language
    const preferredLang = supportedLangs.includes(storedLang) 
      ? storedLang 
      : supportedLangs.includes(browserLang)
        ? browserLang 
        : 'en';
    // Redirect to correct language path
    navigate(`/${preferredLang}${location.pathname}`, { replace: true });
  }, [location.pathname, navigate]);
  return null;
};
export default LanguageWrapper;
