import React from 'react';
import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../../context/ThemeContext';
import { authService } from '../../services/authService';
import { Zap } from 'lucide-react';

const ForgotPassword = () => {
  const { t } = useTranslation('pages');
  const [email, setEmail] = useState('');
  const [isSubmitted, setIsSubmitted] = useState(false);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { theme } = useTheme();

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await authService.forgotPassword(email);
      setIsSubmitted(true);
    } catch (err: any) {
      setError(err.message || t('auth.forgotPassword.failedToSendReset'));
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
            {t('auth.forgotPassword.title')}
          </h2>
          <p className={`mt-2 text-center text-sm ${theme === 'light' ? 'text-gray-600' : 'text-gray-300'}`}>
            {t('common.or')}{' '}
            <Link to="/login" className={`font-medium ${theme === 'light' ? 'text-green-600 hover:text-green-500' : 'text-green-400 hover:text-green-300'}`}>
              {t('auth.forgotPassword.returnToSignIn')}
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

        {isSubmitted ? (
          <div className={`rounded-md ${theme === 'light' ? 'bg-green-50' : 'bg-green-900 bg-opacity-20'} p-4`}>
            <div className="flex">
              <div className="flex-shrink-0">
                <svg className={`h-5 w-5 ${theme === 'light' ? 'text-green-400' : 'text-green-500'}`} viewBox="0 0 20 20" fill="currentColor">
                  <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                </svg>
              </div>
              <div className="ml-3">
                <p className={`text-sm font-medium ${theme === 'light' ? 'text-green-800' : 'text-green-300'}`}>
                  {t('auth.forgotPassword.resetLinkSent', { email })}
                </p>
              </div>
            </div>
            <div className="mt-4 text-center">
              <Link
                to="/login"
                className={`font-medium ${theme === 'light' ? 'text-green-600 hover:text-green-500' : 'text-green-400 hover:text-green-300'}`}
              >
                {t('auth.forgotPassword.returnToSignIn')}
              </Link>
            </div>
          </div>
        ) : (
          <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
            <div>
              <label htmlFor="email-address" className={`block text-sm font-medium ${theme === 'light' ? 'text-gray-700' : 'text-gray-300'}`}>
                {t('auth.forgotPassword.emailLabel')}
              </label>
              <div className="mt-1">
                <input
                  id="email-address"
                  name="email"
                  type="email"
                  autoComplete="email"
                  required
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className={`appearance-none block w-full px-3 py-2 border ${theme === 'light'
                    ? 'border-gray-300 placeholder-gray-400 text-gray-900 focus:ring-green-500 focus:border-green-500 bg-white'
                    : 'border-gray-700 placeholder-gray-500 text-white focus:ring-green-500 focus:border-green-500 bg-gray-800'
                    } rounded-md shadow-sm focus:outline-none sm:text-sm`}
                  placeholder={t('auth.forgotPassword.emailLabel')}
                />
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
                {isLoading ? t('auth.forgotPassword.sending') : t('auth.forgotPassword.sendResetLink')}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
};

export default ForgotPassword;
