import React from 'react';
import { JSX, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../../context/ThemeContext';

interface SettingsProps {
  section?: string | null;
}

const Settings: React.FC<SettingsProps> = ({ section = null }) => {
  const navigate = useNavigate();
  const params = useParams();
  const { theme } = useTheme();
  const { t } = useTranslation();

  const currentSection = section || 'general';

  useEffect(() => {
    if (!section && location.pathname === '/settings') {
      navigate('/settings/general');
    }
  }, [section, navigate]);

  const renderSectionContent = (): JSX.Element | null => {
    switch (currentSection) {
      case 'general':
        return (
          <div>
            <h3 className={`text-lg font-semibold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
              {t('settings.general.title')}
            </h3>
            <div className={`bg-white rounded-lg shadow ${theme === 'dark' ? 'bg-gray-800' : ''}`}>
              <div className={`p-6 border-b ${theme === 'dark' ? 'border-gray-700' : 'border-gray-200'}`}>
                <h4 className={`text-md font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                  {t('settings.general.appPreferences')}
                </h4>
                <p className={`mt-1 text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                  {t('settings.general.appPreferencesDescription')}
                </p>
              </div>
              <div className="p-6">
                <div className="space-y-4">
                  <div>
                    <label className={`block text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      {t('settings.general.appName')}
                    </label>
                    <input
                      type="text"
                      className={`mt-1 block w-full rounded-md shadow-sm ${theme === 'dark'
                        ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
                        : 'border-gray-300 focus:ring-green-500 focus:border-green-500'
                        }`}
                      defaultValue={t('sidebar.logo') + 'API'}
                    />
                  </div>

                  <div>
                    <label className={`block text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      {t('settings.general.displayLanguage')}
                    </label>
                    <select
                      className={`mt-1 block w-full rounded-md shadow-sm ${theme === 'dark'
                        ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
                        : 'border-gray-300 focus:ring-green-500 focus:border-green-500'
                        }`}
                    >
                      <option>English</option>
                      <option>Nederlands</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          </div>
        );

      case 'security':
        return (
          <div>
            <h3 className={`text-lg font-semibold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
              {t('settings.security.title')}
            </h3>
            <div className={`bg-white rounded-lg shadow ${theme === 'dark' ? 'bg-gray-800' : ''}`}>
              <div className={`p-6 border-b ${theme === 'dark' ? 'border-gray-700' : 'border-gray-200'}`}>
                <h4 className={`text-md font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                  {t('settings.security.authSecurity')}
                </h4>
                <p className={`mt-1 text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                  {t('settings.security.authSecurityDescription')}
                </p>
              </div>
              <div className="p-6">
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <h5 className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                        {t('settings.security.twoFactorAuth')}
                      </h5>
                      <p className={`text-xs ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                        {t('settings.security.twoFactorAuthDescription')}
                      </p>
                    </div>
                    <div className="relative inline-block w-10 mr-2 align-middle select-none">
                      <input type="checkbox" id="toggle" className="sr-only" />
                      <label htmlFor="toggle" className={`block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer ${theme === 'dark' ? 'bg-gray-700' : ''}`}>
                        <span className={`absolute left-0 top-0 h-6 w-6 rounded-full bg-white transition-transform duration-200 ease-in ${true ? 'transform translate-x-4' : ''}`}></span>
                      </label>
                    </div>
                  </div>

                  <div>
                    <button
                      className={`px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 ${theme === 'dark' ? 'focus:ring-offset-gray-900' : ''}`}
                    >
                      {t('settings.security.changePassword')}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        );

      case 'account':
        return (
          <div>
            <h3 className={`text-lg font-semibold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
              {t('settings.account.title')}
            </h3>
            <div className={`bg-white rounded-lg shadow ${theme === 'dark' ? 'bg-gray-800' : ''}`}>
              <div className={`p-6 border-b ${theme === 'dark' ? 'border-gray-700' : 'border-gray-200'}`}>
                <h4 className={`text-md font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                  {t('settings.account.profileInfo')}
                </h4>
                <p className={`mt-1 text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                  {t('settings.account.profileInfoDescription')}
                </p>
              </div>
              <div className="p-6">
                <div className="space-y-4">
                  <div>
                    <label className={`block text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      {t('settings.account.fullName')}
                    </label>
                    <input
                      type="text"
                      className={`mt-1 block w-full rounded-md shadow-sm ${theme === 'dark'
                        ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
                        : 'border-gray-300 focus:ring-green-500 focus:border-green-500'
                        }`}
                      defaultValue={t('sidebar.admin')}
                    />
                  </div>

                  <div>
                    <label className={`block text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      {t('settings.account.emailAddress')}
                    </label>
                    <input
                      type="email"
                      className={`mt-1 block w-full rounded-md shadow-sm ${theme === 'dark'
                        ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
                        : 'border-gray-300 focus:ring-green-500 focus:border-green-500'
                        }`}
                      defaultValue={t('sidebar.adminEmail')}
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        );

      case 'notifications':
        return (
          <div>
            <h3 className={`text-lg font-semibold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
              {t('settings.notifications.title')}
            </h3>
            <div className={`bg-white rounded-lg shadow ${theme === 'dark' ? 'bg-gray-800' : ''}`}>
              <div className={`p-6 border-b ${theme === 'dark' ? 'border-gray-700' : 'border-gray-200'}`}>
                <h4 className={`text-md font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                  {t('settings.notifications.notificationPreferences')}
                </h4>
                <p className={`mt-1 text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                  {t('settings.notifications.notificationPreferencesDescription')}
                </p>
              </div>
              <div className="p-6">
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <h5 className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                        {t('settings.notifications.emailNotifications')}
                      </h5>
                      <p className={`text-xs ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                        {t('settings.notifications.emailNotificationsDescription')}
                      </p>
                    </div>
                    <div className="relative inline-block w-10 mr-2 align-middle select-none">
                      <input type="checkbox" id="email-toggle" className="sr-only" defaultChecked />
                      <label htmlFor="email-toggle" className={`block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer ${theme === 'dark' ? 'bg-gray-700' : ''}`}>
                        <span className="absolute left-0 top-0 h-6 w-6 rounded-full bg-white transition-transform duration-200 ease-in transform translate-x-4"></span>
                      </label>
                    </div>
                  </div>

                  <div className="flex items-center justify-between">
                    <div>
                      <h5 className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                        {t('settings.notifications.browserNotifications')}
                      </h5>
                      <p className={`text-xs ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                        {t('settings.notifications.browserNotificationsDescription')}
                      </p>
                    </div>
                    <div className="relative inline-block w-10 mr-2 align-middle select-none">
                      <input type="checkbox" id="browser-toggle" className="sr-only" />
                      <label htmlFor="browser-toggle" className={`block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer ${theme === 'dark' ? 'bg-gray-700' : ''}`}>
                        <span className="absolute left-0 top-0 h-6 w-6 rounded-full bg-white transition-transform duration-200 ease-in"></span>
                      </label>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        );

      case 'regional':
        return (
          <div>
            <h3 className={`text-lg font-semibold mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
              {t('settings.regional.title')}
            </h3>
            <div className={`bg-white rounded-lg shadow ${theme === 'dark' ? 'bg-gray-800' : ''}`}>
              <div className={`p-6 border-b ${theme === 'dark' ? 'border-gray-700' : 'border-gray-200'}`}>
                <h4 className={`text-md font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                  {t('settings.regional.locationFormat')}
                </h4>
                <p className={`mt-1 text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                  {t('settings.regional.locationFormatDescription')}
                </p>
              </div>
              <div className="p-6">
                <div className="space-y-4">
                  <div>
                    <label className={`block text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      {t('settings.regional.timezone')}
                    </label>
                    <select
                      className={`mt-1 block w-full rounded-md shadow-sm ${theme === 'dark'
                        ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
                        : 'border-gray-300 focus:ring-green-500 focus:border-green-500'
                        }`}
                    >
                      <option>UTC (Coordinated Universal Time)</option>
                      <option>America/New_York (Eastern Time)</option>
                      <option>America/Chicago (Central Time)</option>
                      <option>America/Denver (Mountain Time)</option>
                      <option>America/Los_Angeles (Pacific Time)</option>
                      <option>Europe/Amsterdam (Central European Time)</option>
                    </select>
                  </div>

                  <div>
                    <label className={`block text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'}`}>
                      {t('settings.regional.dateFormat')}
                    </label>
                    <select
                      className={`mt-1 block w-full rounded-md shadow-sm ${theme === 'dark'
                        ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
                        : 'border-gray-300 focus:ring-green-500 focus:border-green-500'
                        }`}
                    >
                      <option>MM/DD/YYYY</option>
                      <option>DD/MM/YYYY</option>
                      <option>YYYY-MM-DD</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          </div>
        );

      default:
        return <div>Select a settings category</div>;
    }
  };

  return (
    <div className={`px-4 py-6 ${theme === 'dark' ? 'text-white' : ''}`}>
      <h2 className="text-2xl font-bold mb-6">{t('settings.title')}</h2>
      {renderSectionContent()}
    </div>
  );
};

export default Settings;
