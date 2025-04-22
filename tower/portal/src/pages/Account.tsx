import React, { useEffect, useState } from 'react';
import { Save, Loader2, User, Mail, Building, Globe, Key, CheckCircle, AlertCircle } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { authService } from '../services/authService';
import { useTheme } from '../context/ThemeContext';

const Account: React.FC = () => {
  const { t } = useTranslation('pages');
  const { theme } = useTheme();
  const [userData, setUserData] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [saving, setSaving] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);
  const [changePassword, setChangePassword] = useState<boolean>(false);
  const [passwordData, setPasswordData] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  });
  const [passwordError, setPasswordError] = useState<string | null>(null);
  const [passwordSuccess, setPasswordSuccess] = useState<boolean>(false);

  const formatDate = (dateString: string | null): string => {
    if (!dateString) return t('users.never');
    const date = new Date(dateString);
    return new Intl.DateTimeFormat(undefined, {
      dateStyle: 'medium',
      timeStyle: 'short',
    }).format(date);
  };

  useEffect(() => {
    fetchUserData();
  }, []);

  const fetchUserData = async () => {
    try {
      setLoading(true);
      // Replace with your actual service call
      const user = await authService.getCurrentUser();
      setUserData(user);
      setError(null);
    } catch (err) {
      console.error(err);
      setError(t('account.loadError', 'Failed to load account information'));
    } finally {
      setLoading(false);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    if (userData) {
      setUserData({
        ...userData,
        [name]: value,
      });
    }
    // Clear success message when form is modified
    setSuccess(false);
  };

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setPasswordData({
      ...passwordData,
      [name]: value,
    });
    // Clear success/error messages when form is modified
    setPasswordSuccess(false);
    setPasswordError(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!userData) return;

    try {
      setSaving(true);
      // Replace with your actual service call
      await authService.updateUserProfile(userData);
      setSuccess(true);
      setError(null);
    } catch (err) {
      console.error(err);
      setError(t('account.saveError', 'Failed to update account information'));
      setSuccess(false);
    } finally {
      setSaving(false);
    }
  };

  const handlePasswordSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (passwordData.newPassword !== passwordData.confirmPassword) {
      setPasswordError(t('account.passwordMismatch', 'New passwords do not match'));
      return;
    }

    try {
      setSaving(true);
      // Replace with your actual service call
      await authService.updatePassword(passwordData.currentPassword, passwordData.newPassword);
      setPasswordSuccess(true);
      setPasswordError(null);
      // Reset password form
      setPasswordData({
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      });
    } catch (err) {
      console.error(err);
      setPasswordError(t('account.passwordError', 'Failed to update password'));
      setPasswordSuccess(false);
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="container mx-auto py-8 px-4">
      <h1 className={`text-2xl font-bold mb-6 ${theme === 'dark' ? 'text-white' : 'text-gray-800'}`}>
        {t('account.title', 'Account Settings')}
      </h1>

      {loading ? (
        <div className="p-8 text-center">
          <Loader2 className="animate-spin h-8 w-8 mx-auto text-green-500 mb-2" />
          <p className="text-gray-500 dark:text-gray-400">{t('account.loading', 'Loading account information...')}</p>
        </div>
      ) : error && !userData ? (
        <div className={`p-6 rounded-lg border ${theme === 'dark' ? 'bg-red-900/30 border-red-800 text-red-300' : 'bg-red-50 border-red-200 text-red-600'}`}>
          <p>{error}</p>
          <button
            onClick={fetchUserData}
            className="mt-2 px-3 py-1 rounded-md bg-gray-600 text-white hover:bg-gray-700 text-sm"
          >
            {t('common.retry', 'Retry')}
          </button>
        </div>
      ) : userData && (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {/* Profile Information */}
          <div className="md:col-span-2">
            <div className={`rounded-lg shadow-sm ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'} overflow-hidden`}>
              <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
                <h2 className="text-lg font-medium">{t('account.profileInfo', 'Profile Information')}</h2>
              </div>

              {error && (
                <div className={`mx-6 mt-4 p-3 rounded-md ${theme === 'dark' ? 'bg-red-900/30 border border-red-800 text-red-300' : 'bg-red-50 border border-red-200 text-red-600'}`}>
                  <p className="flex items-center">
                    <AlertCircle className="w-4 h-4 mr-2" />
                    {error}
                  </p>
                </div>
              )}

              {success && (
                <div className={`mx-6 mt-4 p-3 rounded-md ${theme === 'dark' ? 'bg-green-900/30 border border-green-800 text-green-300' : 'bg-green-50 border border-green-200 text-green-600'}`}>
                  <p className="flex items-center">
                    <CheckCircle className="w-4 h-4 mr-2" />
                    {t('account.saveSuccess', 'Account information updated successfully')}
                  </p>
                </div>
              )}

              <form onSubmit={handleSubmit} className="p-6">
                <div className="space-y-4">
                  <div>
                    <label htmlFor="name" className={`block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      {t('common.name', 'Name')}
                    </label>
                    <div className="relative rounded-md shadow-sm">
                      <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                        <User className={`h-5 w-5 ${theme === 'dark' ? 'text-gray-500' : 'text-gray-400'}`} />
                      </div>
                      <input
                        type="text"
                        name="name"
                        id="name"
                        value={userData.name}
                        onChange={handleInputChange}
                        className={`block w-full pl-10 pr-3 py-2 rounded-md focus:ring-2 ${theme === 'dark'
                          ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
                          : 'bg-white border-gray-300 text-gray-900 focus:ring-green-500 focus:border-green-500'
                          }`}
                        required
                      />
                    </div>
                  </div>

                  <div>
                    <label htmlFor="email" className={`block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      {t('common.email', 'Email')}
                    </label>
                    <div className="relative rounded-md shadow-sm">
                      <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                        <Mail className={`h-5 w-5 ${theme === 'dark' ? 'text-gray-500' : 'text-gray-400'}`} />
                      </div>
                      <input
                        type="email"
                        name="email"
                        id="email"
                        value={userData.email}
                        onChange={handleInputChange}
                        className={`block w-full pl-10 pr-3 py-2 rounded-md focus:ring-2 ${theme === 'dark'
                          ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
                          : 'bg-white border-gray-300 text-gray-900 focus:ring-green-500 focus:border-green-500'
                          }`}
                        required
                      />
                    </div>
                  </div>

                  <div>
                    <label htmlFor="company" className={`block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      {t('common.company', 'Company')}
                    </label>
                    <div className="relative rounded-md shadow-sm">
                      <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                        <Building className={`h-5 w-5 ${theme === 'dark' ? 'text-gray-500' : 'text-gray-400'}`} />
                      </div>
                      <input
                        type="text"
                        name="company"
                        id="company"
                        value={userData.company}
                        onChange={handleInputChange}
                        className={`block w-full pl-10 pr-3 py-2 rounded-md focus:ring-2 ${theme === 'dark'
                          ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
                          : 'bg-white border-gray-300 text-gray-900 focus:ring-green-500 focus:border-green-500'
                          }`}
                      />
                    </div>
                  </div>

                  <div>
                    <label htmlFor="website" className={`block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      {t('common.website', 'Website')}
                    </label>
                    <div className="relative rounded-md shadow-sm">
                      <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                        <Globe className={`h-5 w-5 ${theme === 'dark' ? 'text-gray-500' : 'text-gray-400'}`} />
                      </div>
                      <input
                        type="url"
                        name="website"
                        id="website"
                        value={userData.website || ''}
                        onChange={handleInputChange}
                        className={`block w-full pl-10 pr-3 py-2 rounded-md focus:ring-2 ${theme === 'dark'
                          ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
                          : 'bg-white border-gray-300 text-gray-900 focus:ring-green-500 focus:border-green-500'
                          }`}
                        placeholder="https://"
                      />
                    </div>
                  </div>
                </div>

                <div className="mt-6">
                  <button
                    type="submit"
                    disabled={saving}
                    className={`inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 ${saving ? 'opacity-75 cursor-not-allowed' : ''
                      }`}
                  >
                    {saving ? (
                      <>
                        <Loader2 className="animate-spin -ml-1 mr-2 h-4 w-4" />
                        {t('common.saving', 'Saving...')}
                      </>
                    ) : (
                      <>
                        <Save className="-ml-1 mr-2 h-4 w-4" />
                        {t('common.save', 'Save Changes')}
                      </>
                    )}
                  </button>
                </div>
              </form>
            </div>

            {/* Password Change Section */}
            <div className={`mt-6 rounded-lg shadow-sm ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'} overflow-hidden`}>
              <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
                <h2 className="text-lg font-medium">{t('account.changePassword', 'Change Password')}</h2>
              </div>

              {passwordError && (
                <div className={`mx-6 mt-4 p-3 rounded-md ${theme === 'dark' ? 'bg-red-900/30 border border-red-800 text-red-300' : 'bg-red-50 border border-red-200 text-red-600'}`}>
                  <p className="flex items-center">
                    <AlertCircle className="w-4 h-4 mr-2" />
                    {passwordError}
                  </p>
                </div>
              )}

              {passwordSuccess && (
                <div className={`mx-6 mt-4 p-3 rounded-md ${theme === 'dark' ? 'bg-green-900/30 border border-green-800 text-green-300' : 'bg-green-50 border border-green-200 text-green-600'}`}>
                  <p className="flex items-center">
                    <CheckCircle className="w-4 h-4 mr-2" />
                    {t('account.passwordSuccess', 'Password updated successfully')}
                  </p>
                </div>
              )}

              {!changePassword ? (
                <div className="p-6">
                  <button
                    onClick={() => setChangePassword(true)}
                    className={`inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium ${theme === 'dark'
                      ? 'bg-gray-700 text-white hover:bg-gray-600'
                      : 'bg-gray-200 text-gray-800 hover:bg-gray-300'
                      }`}
                  >
                    <Key className="-ml-1 mr-2 h-4 w-4" />
                    {t('account.changePasswordBtn', 'Change Password')}
                  </button>
                </div>
              ) : (
                <form onSubmit={handlePasswordSubmit} className="p-6">
                  <div className="space-y-4">
                    <div>
                      <label htmlFor="currentPassword" className={`block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                        {t('account.currentPassword', 'Current Password')}
                      </label>
                      <input
                        type="password"
                        name="currentPassword"
                        id="currentPassword"
                        value={passwordData.currentPassword}
                        onChange={handlePasswordChange}
                        className={`block w-full px-3 py-2 rounded-md focus:ring-2 ${theme === 'dark'
                          ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
                          : 'bg-white border-gray-300 text-gray-900 focus:ring-green-500 focus:border-green-500'
                          }`}
                        required
                      />
                    </div>

                    <div>
                      <label htmlFor="newPassword" className={`block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                        {t('account.newPassword', 'New Password')}
                      </label>
                      <input
                        type="password"
                        name="newPassword"
                        id="newPassword"
                        value={passwordData.newPassword}
                        onChange={handlePasswordChange}
                        className={`block w-full px-3 py-2 rounded-md focus:ring-2 ${theme === 'dark'
                          ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
                          : 'bg-white border-gray-300 text-gray-900 focus:ring-green-500 focus:border-green-500'
                          }`}
                        required
                      />
                    </div>

                    <div>
                      <label htmlFor="confirmPassword" className={`block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                        {t('account.confirmPassword', 'Confirm New Password')}
                      </label>
                      <input
                        type="password"
                        name="confirmPassword"
                        id="confirmPassword"
                        value={passwordData.confirmPassword}
                        onChange={handlePasswordChange}
                        className={`block w-full px-3 py-2 rounded-md focus:ring-2 ${theme === 'dark'
                          ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
                          : 'bg-white border-gray-300 text-gray-900 focus:ring-green-500 focus:border-green-500'
                          }`}
                        required
                      />
                    </div>
                  </div>

                  <div className="mt-6 flex space-x-3">
                    <button
                      type="submit"
                      disabled={saving}
                      className={`inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 ${saving ? 'opacity-75 cursor-not-allowed' : ''
                        }`}
                    >
                      {saving ? (
                        <>
                          <Loader2 className="animate-spin -ml-1 mr-2 h-4 w-4" />
                          {t('common.updating', 'Updating...')}
                        </>
                      ) : (
                        <>
                          <Key className="-ml-1 mr-2 h-4 w-4" />
                          {t('account.updatePassword', 'Update Password')}
                        </>
                      )}
                    </button>
                    <button
                      type="button"
                      onClick={() => setChangePassword(false)}
                      className={`inline-flex items-center px-4 py-2 border rounded-md shadow-sm text-sm font-medium ${theme === 'dark'
                        ? 'border-gray-600 bg-transparent text-gray-300 hover:bg-gray-700'
                        : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50'
                        }`}
                    >
                      {t('common.cancel', 'Cancel')}
                    </button>
                  </div>
                </form>
              )}
            </div>
          </div>

          {/* Account Info Sidebar */}
          <div className="md:col-span-1">
            <div className={`rounded-lg shadow-sm ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'} overflow-hidden`}>
              <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
                <h2 className="text-lg font-medium">{t('account.accountInfo', 'Account Information')}</h2>
              </div>
              <div className="p-6">
                <div className="space-y-4">
                  <div>
                    <h3 className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                      {t('common.userId', 'User ID')}
                    </h3>
                    <p className={`mt-1 text-sm font-mono break-all ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      {userData.id}
                    </p>
                  </div>
                  <div>
                    <h3 className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                      {t('common.role', 'Role')}
                    </h3>
                    <p className={`mt-1 text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      {userData.role}
                    </p>
                  </div>
                  <div>
                    <h3 className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                      {t('common.memberSince', 'Member Since')}
                    </h3>
                    <p className={`mt-1 text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      {formatDate(userData.created_at)}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Account;
