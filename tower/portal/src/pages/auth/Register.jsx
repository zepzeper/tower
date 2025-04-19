import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../../context/AuthContext';
import { useTheme } from '../../context/ThemeContext';
import { Zap } from 'lucide-react';

const Register = () => {
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        password: '',
        confirmPassword: ''
    });
    const [error, setError] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const navigate = useNavigate();
    const { register } = useAuth();
    const { theme } = useTheme();
    const { t } = useTranslation();

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData(prev => ({ ...prev, [name]: value }));
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError('');

        // Basic validation
        if (formData.password !== formData.confirmPassword) {
            setError(t('auth.register.passwordMismatch'));
            return;
        }

        setIsLoading(true);

        try {
            await register(formData);
            navigate('/');
        } catch (err) {
            setError(err.message || 'Registration failed');
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
                        {t('auth.register.title')}
                    </h2>
                    <p className={`mt-2 text-center text-sm ${theme === 'light' ? 'text-gray-600' : 'text-gray-300'}`}>
                        {t('common.or')}{' '}
                        <Link to="/login" className={`font-medium ${theme === 'light' ? 'text-green-600 hover:text-green-500' : 'text-green-400 hover:text-green-300'}`}>
                            {t('auth.register.signInExisting')}
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
                            <label htmlFor="name" className="sr-only">{t('auth.register.nameLabel')}</label>
                            <input
                                id="name"
                                name="name"
                                type="text"
                                required
                                value={formData.name}
                                onChange={handleChange}
                                className={`appearance-none rounded-none relative block w-full px-3 py-2 border ${
                                    theme === 'light' 
                                        ? 'border-gray-300 placeholder-gray-500 text-gray-900 focus:ring-green-500 focus:border-green-500'
                                        : 'border-gray-700 placeholder-gray-400 text-white bg-gray-800 focus:ring-green-500 focus:border-green-500'
                                } rounded-t-md focus:outline-none focus:z-10 sm:text-sm`}
                                placeholder={t('auth.register.nameLabel')}
                            />
                        </div>
                        <div>
                            <label htmlFor="email-address" className="sr-only">{t('auth.register.emailLabel')}</label>
                            <input
                                id="email-address"
                                name="email"
                                type="email"
                                autoComplete="email"
                                required
                                value={formData.email}
                                onChange={handleChange}
                                className={`appearance-none rounded-none relative block w-full px-3 py-2 border ${
                                    theme === 'light' 
                                        ? 'border-gray-300 placeholder-gray-500 text-gray-900 focus:ring-green-500 focus:border-green-500'
                                        : 'border-gray-700 placeholder-gray-400 text-white bg-gray-800 focus:ring-green-500 focus:border-green-500'
                                } focus:outline-none focus:z-10 sm:text-sm`}
                                placeholder={t('auth.register.emailLabel')}
                            />
                        </div>
                        <div>
                            <label htmlFor="password" className="sr-only">{t('auth.register.passwordLabel')}</label>
                            <input
                                id="password"
                                name="password"
                                type="password"
                                autoComplete="new-password"
                                required
                                value={formData.password}
                                onChange={handleChange}
                                className={`appearance-none rounded-none relative block w-full px-3 py-2 border ${
                                    theme === 'light' 
                                        ? 'border-gray-300 placeholder-gray-500 text-gray-900 focus:ring-green-500 focus:border-green-500'
                                        : 'border-gray-700 placeholder-gray-400 text-white bg-gray-800 focus:ring-green-500 focus:border-green-500'
                                } focus:outline-none focus:z-10 sm:text-sm`}
                                placeholder={t('auth.register.passwordLabel')}
                            />
                        </div>
                        <div>
                            <label htmlFor="confirmPassword" className="sr-only">{t('auth.register.confirmPasswordLabel')}</label>
                            <input
                                id="confirmPassword"
                                name="confirmPassword"
                                type="password"
                                autoComplete="new-password"
                                required
                                value={formData.confirmPassword}
                                onChange={handleChange}
                                className={`appearance-none rounded-none relative block w-full px-3 py-2 border ${
                                    theme === 'light' 
                                        ? 'border-gray-300 placeholder-gray-500 text-gray-900 focus:ring-green-500 focus:border-green-500'
                                        : 'border-gray-700 placeholder-gray-400 text-white bg-gray-800 focus:ring-green-500 focus:border-green-500'
                                } rounded-b-md focus:outline-none focus:z-10 sm:text-sm`}
                                placeholder={t('auth.register.confirmPasswordLabel')}
                            />
                        </div>
                    </div>

                    <div>
                        <button
                            type="submit"
                            disabled={isLoading}
                            className={`group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white ${
                                theme === 'light' 
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
                                        <path fillRule="evenodd" d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z" clipRule="evenodd" />
                                    </svg>
                                </span>
                            )}
                            {isLoading ? t('auth.register.creatingAccount') : t('auth.register.createAccount')}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default Register;
