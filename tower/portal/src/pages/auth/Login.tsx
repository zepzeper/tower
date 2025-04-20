import React from 'react';
import { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../../context/AuthContext';
import { useTheme } from '../../context/ThemeContext';
import { Zap } from 'lucide-react';

interface LocationState {
  from?: {
    pathname: string;
  };
}

const Login: React.FC = () => {
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [error, setError] = useState<string>('');
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();
  const { theme } = useTheme();
  const { t } = useTranslation();

  const from = (location.state as LocationState)?.from?.pathname || '/';

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await login(email, password);
      navigate(from, { replace: true });
    } catch (err: any) {
      setError(err.message || t('auth.login.invalidCredentials'));
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className={`min-h-screen flex items-center justify-center ${theme === 'light' ? 'bg-gray-50' : 'bg-gray-900'} py-12 px-4 sm:px-6 lg:px-8 transition-colors duration-200`}>
      <div className={`max-w-md w-full space-y-8 ${theme === 'dark' ? 'text-white' : ''}`}>
        <div>
          <div className="flex justify-center">
            <div className="flex items-center space-x-2">
              <Zap className={`w-10 h-10 ${theme === 'light' ? 'text-green-600' : 'text-green-400'}`} />
              <span className={`text-3xl font-bold ${theme === 'light' ? 'text-gray-900' : 'text-white'}`}>
                {t('sidebar.logo')}<span className={`${theme === 'light' ? 'text-green-600' : 'text-green-400'}`}>API</span>
              </span>
            </div>
          </div>
          <h2 className={`mt-6 text-center text-3xl font-extrabold ${theme === 'light' ? 'text-gray-900' : 'text-white'}`}>
            {t('auth.login.title')}
          </h2>
          <p className={`mt-2 text-center text-sm ${theme === 'light' ? 'text-gray-600' : 'text-gray-300'}`}>
            {t('common.or')}{' '}
            <Link to="/register" className={`font-medium ${theme === 'light' ? 'text-green-600 hover:text-green-500' : 'text-green-400 hover:text-green-300'}`}>
              {t('auth.login.createAccount')}
            </Link>
          </p>
        </div>

        {error && (
          <div className={`rounded-md ${theme === 'light' ? 'bg-red-50' : 'bg-red-900 bg-opacity-20'} p-4`}>
            <div className="flex">
              <div className="flex-shrink-0">
                <svg className={`h-5 w-5 ${theme === 'light' ? 'text-red-400' : 'text-red-500'}`} viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                </svg>
              </div>
              <div className="ml-3">
                <p className={`text-sm font-medium ${theme === 'light' ? 'text-red-800' : 'text-red-300'}`}>{error}</p>
              </div>
            </div>
          </div>
        )}

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="rounded-md shadow-sm -space-y-px">
            <div>
              <label htmlFor="email-address" className="sr-only">{t('auth.login.emailLabel')}</label>
              <input
                id="email-address"
                name="email"
                type="email"
                autoComplete="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className={`appearance-none rounded-none relative block w-full px-3 py-2 border ${theme === 'light'
                    ? 'border-gray-300 placeholder-gray-500 text-gray-900 focus:ring-green-500 focus:border-green-500'
                    : 'border-gray-700 placeholder-gray-400 text-white bg-gray-800 focus:ring-green-500 focus:border-green-500'
                  } rounded-t-md focus:outline-none focus:z-10 sm:text-sm`}
                placeholder={t('auth.login.emailLabel')}
              />
            </div>
            <div>
              <label htmlFor="password" className="sr-only">{t('auth.login.passwordLabel')}</label>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="current-password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className={`appearance-none rounded-none relative block w-full px-3 py-2 border ${theme === 'light'
                    ? 'border-gray-300 placeholder-gray-500 text-gray-900 focus:ring-green-500 focus:border-green-500'
                    : 'border-gray-700 placeholder-gray-400 text-white bg-gray-800 focus:ring-green-500 focus:border-green-500'
                  } rounded-b-md focus:outline-none focus:z-10 sm:text-sm`}
                placeholder={t('auth.login.passwordLabel')}
              />
            </div>
          </div>

          <div className="flex items-center justify-between">
            <div className="flex items-center">
              <input
                id="remember-me"
                name="remember-me"
                type="checkbox"
                className={`h-4 w-4 ${theme === 'light' ? 'text-green-600 focus:ring-green-500 border-gray-300' : 'text-green-500 focus:ring-green-600 border-gray-700 bg-gray-700'} rounded`}
              />
              <label htmlFor="remember-me" className={`ml-2 block text-sm ${theme === 'light' ? 'text-gray-900' : 'text-gray-300'}`}>
                {t('auth.login.rememberMe')}
              </label>
            </div>

            <div className="text-sm">
              <Link to="/forgot-password" className={`font-medium ${theme === 'light' ? 'text-green-600 hover:text-green-500' : 'text-green-400 hover:text-green-300'}`}>
                {t('auth.login.forgotPassword')}
              </Link>
            </div>
          </div>

          <div>
            <button
              type="submit"
              disabled={isLoading}
              className={`group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white ${theme === 'light'
                  ? 'bg-green-600 hover:bg-green-700 focus:ring-green-500'
                  : 'bg-green-600 hover:bg-green-700 focus:ring-green-400'
                } focus:outline-none focus:ring-2 focus:ring-offset-2`}
            >
              {isLoading ? (
                <span className="absolute left-0 inset-y-0 flex items-center pl-3">
                  <svg className="animate-spin h-5 w-5 text-green-300" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                </span>
              ) : (
                <span className="absolute left-0 inset-y-0 flex items-center pl-3">
                  <svg className={`h-5 w-5 ${theme === 'light' ? 'text-green-500 group-hover:text-green-400' : 'text-green-300 group-hover:text-green-200'}`} viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fillRule="evenodd" d="M10 1a4.5 4.5 0 00-4.5 4.5V9H5a2 2 0 00-2 2v6a2 2 0 002 2h10a2 2 0 002-2v-6a2 2 0 00-2-2h-.5V5.5A4.5 4.5 0 0010 1zm3 8V5.5a3 3 0 10-6 0V9h6z" clipRule="evenodd" />
                  </svg>
                </span>
              )}
              {isLoading ? t('auth.login.signingIn') : t('auth.login.signIn')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Login;
