import React from 'react';
import { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../../context/ThemeContext';
import { authService } from '../../services/authService';
import { Zap } from 'lucide-react';

const ResetPassword: React.FC = () => {
  const [password, setPassword] = useState<string>('');
  const [confirmPassword, setConfirmPassword] = useState<string>('');
  const [token, setToken] = useState<string>('');
  const [error, setError] = useState<string>('');
  const [isSuccess, setIsSuccess] = useState<boolean>(false);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const location = useLocation();
  const { theme } = useTheme();
  const { t } = useTranslation();

  useEffect(() => {
    const searchParams = new URLSearchParams(location.search);
    const tokenParam = searchParams.get('token');
    if (tokenParam) {
      setToken(tokenParam);
    } else {
      setError(t('auth.resetPassword.invalidToken'));
    }
  }, [location.search, t]);

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');

    if (password !== confirmPassword) {
      setError(t('auth.resetPassword.passwordMismatch'));
      return;
    }

    if (password.length < 8) {
      setError(t('auth.resetPassword.passwordTooShort'));
      return;
    }

    setIsLoading(true);

    try {
      await authService.resetPassword(token, password);
      setIsSuccess(true);
    } catch (err: any) {
      setError(err.message || t('auth.resetPassword.genericError'));
    } finally {
      setIsLoading(false);
    }
  };

  if (isSuccess) {
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
              {t('auth.resetPassword.passwordResetSuccess')}
            </h2>
            <div className="mt-4 text-center">
              <p className={`text-sm ${theme === 'light' ? 'text-gray-600' : 'text-gray-300'}`}>
                {t('auth.resetPassword.passwordResetDescription')}
              </p>
              <Link
                to="/login"
                className={`mt-4 inline-block font-medium ${theme === 'light' ? 'text-green-600 hover:text-green-500' : 'text-green-400 hover:text-green-300'}`}
              >
                {t('auth.resetPassword.backToLogin')}
              </Link>
            </div>
          </div>
        </div>
      </div>
    );
  }

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
            {t('auth.resetPassword.title')}
          </h2>
          <p className={`mt-2 text-center text-sm ${theme === 'light' ? 'text-gray-600' : 'text-gray-300'}`}>
            {t('auth.resetPassword.enterNewPassword')}
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
              <label htmlFor="password" className="sr-only">{t('auth.resetPassword.newPasswordLabel')}</label>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="new-password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className={`appearance-none rounded-none relative block w-full px-3 py-2 border ${theme === 'light'
                  ? 'border-gray-300 placeholder-gray-500 text-gray-900 focus:ring-green-500 focus:border-green-500'
                  : 'border-gray-700 placeholder-gray-400 text-white bg-gray-800 focus:ring-green-500 focus:border-green-500'
                  } rounded-t-md focus:outline-none focus:z-10 sm:text-sm`}
                placeholder={t('auth.resetPassword.newPasswordLabel')}
              />
            </div>
            <div>
              <label htmlFor="confirm-password" className="sr-only">{t('auth.resetPassword.confirmPasswordLabel')}</label>
              <input
                id="confirm-password"
                name="confirmPassword"
                type="password"
                autoComplete="new-password"
                required
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                className={`appearance-none rounded-none relative block w-full px-3 py-2 border ${theme === 'light'
                  ? 'border-gray-300 placeholder-gray-500 text-gray-900 focus:ring-green-500 focus:border-green-500'
                  : 'border-gray-700 placeholder-gray-400 text-white bg-gray-800 focus:ring-green-500 focus:border-green-500'
                  } rounded-b-md focus:outline-none focus:z-10 sm:text-sm`}
                placeholder={t('auth.resetPassword.confirmPasswordLabel')}
              />
            </div>
          </div>

          <div>
            <button
              type="submit"
              disabled={isLoading || !token}
              className={`group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white ${theme === 'light'
                ? 'bg-green-600 hover:bg-green-700 focus:ring-green-500'
                : 'bg-green-600 hover:bg-green-700 focus:ring-green-400'
                } focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:bg-gray-400 disabled:cursor-not-allowed`}
            >
              {isLoading ? t('auth.resetPassword.resettingPassword') : t('auth.resetPassword.resetPassword')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default ResetPassword;
