import React, { JSX, useState } from 'react';
import {
  Trash2,
  Plus,
  Key,
  EyeOff,
  Eye,
  XCircle,
  CheckCircle
} from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { useTheme } from '../../../context/ThemeContext';
import { getThemeStyles } from '../../../utility/theme';
import { useConnection } from '../../../context/ConnectionContext';
import { ApiConnectionConfig } from '../../../services/connectionService';
import { Checkbox } from '../../../utility/components/FormComponents';

const ConnectionDetailsConfiguration: React.FC = () => {
  const { t } = useTranslation('pages');
  const { theme } = useTheme();
  const styles = getThemeStyles(theme);

  const {
    connection,
    configs,
    connectionType,
    editing,
    testResult,
    showSecrets,
    setConnection,
    handleAddConfig,
    handleRemoveConfig,
    updateConfigValue,
    handleInitiateOAuth,
    toggleShowSecret,
    formatDate
  } = useConnection();

  const [newConfigKey, setNewConfigKey] = useState<string>('');
  const [newConfigValue, setNewConfigValue] = useState<string>('');
  const [newConfigIsSecret, setNewConfigIsSecret] = useState<boolean>(false);

  const handleAddNewConfig = (): void => {
    if (!newConfigKey.trim()) return;

    handleAddConfig(newConfigKey, newConfigValue, newConfigIsSecret);

    // Reset form fields
    setNewConfigKey('');
    setNewConfigValue('');
    setNewConfigIsSecret(false);
  };

  const renderConfigValue = (config: ApiConnectionConfig): JSX.Element => {
    return (
      <div className="relative w-full">
        <input
          type={config.is_secret && !showSecrets[config.key] ? 'password' : 'text'}
          value={config.value || ''}
          onChange={(e) => updateConfigValue(config.key, e.target.value)}
          className={`w-full px-3 py-2 rounded-md border ${styles.input} ${config.is_secret ? 'pr-10' : ''}`}
          readOnly={!editing}
        />
        {config.is_secret && (
          <button
            type="button"
            onClick={() => toggleShowSecret(config.key)}
            className="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-500"
          >
            {showSecrets[config.key] ? <EyeOff size={18} /> : <Eye size={18} />}
          </button>
        )}
      </div>
    );
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
      {/* Connection Information */}
      <div className={`col-span-1 rounded-lg shadow-sm ${styles.card} p-6`}>
        <h2 className="text-lg font-medium mb-4">{t('connectionDetails.configuration', 'Connection Information')}</h2>

        {editing ? (
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-1">{t('connectionModal.nameLabel')}</label>
              <input
                type="text"
                value={connection!.name}
                onChange={(e) => setConnection({ ...connection!, name: e.target.value })}
                className={`w-full px-3 py-2 rounded-md border ${styles.input}`}
              />
            </div>

            <div>
              <label className="block text-sm font-medium mb-1">{t('connectionModal.descriptionLabel')}</label>
              <textarea
                value={connection!.description || ''}
                onChange={(e) => setConnection({ ...connection!, description: e.target.value })}
                rows={4}
                className={`w-full px-3 py-2 rounded-md border ${styles.input}`}
                placeholder={t('connectionModal.descriptionPlaceholder')}
              />
            </div>

            <div className="flex items-center">
              <Checkbox
                label={t('connectionList.statusActive')}
                checked={connection!.active}
                onChange={(e) => setConnection({ ...connection!, active: e })}
                className="mr-2 h-4 w-4"
                id="active-toggle"
              />
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <div>
              <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400">{t('connectionDetails.type')}</h3>
              <p className="mt-1">{connection!.type}</p>
            </div>

            <div>
              <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400">{t('connectionModal.descriptionLabel')}</h3>
              <p className="mt-1">{connection!.description || t('connectionDetails.noDescription', 'No description provided.')}</p>
            </div>

            <div>
              <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400">{t('connectionDetails.created')}</h3>
              <p className="mt-1">{formatDate(connection!.created_at)}</p>
            </div>

            <div>
              <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400">{t('connectionDetails.updated')}</h3>
              <p className="mt-1">{formatDate(connection!.updated_at || connection!.created_at)}</p>
            </div>
          </div>
        )}

        {connectionType?.authType === 'oauth2' && (
          <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
            <button
              onClick={handleInitiateOAuth}
              className="w-full flex items-center justify-center px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700"
            >
              <Key className="w-4 h-4 mr-2" />
              {t('connectionModal.connectWithOAuth', 'Authenticate with OAuth')}
            </button>
          </div>
        )}
      </div>

      {/* Configuration */}
      <div className={`col-span-2 rounded-lg shadow-sm ${styles.card} p-6`}>
        <div className="flex justify-between items-center mb-4">
          <h2 className="text-lg font-medium">{t('connectionDetails.configuration')}</h2>
          {editing && (
            <button
              onClick={() => {
                document.getElementById('add-config-section')?.scrollIntoView({ behavior: 'smooth' });
              }}
              className={`px-3 py-1 rounded-md text-sm flex items-center ${styles.secondaryButton}`}
            >
              <Plus className="w-4 h-4 mr-1" />
              {t('common.add', 'Add Config')}
            </button>
          )}
        </div>

        {configs.length === 0 ? (
          <div className={`p-8 text-center rounded-md ${theme === 'dark' ? 'bg-gray-700' : 'bg-gray-50'}`}>
            <p className="text-gray-500 dark:text-gray-400">{t('connectionDetails.noConfig')}</p>
            {editing && (
              <button
                onClick={() => {
                  document.getElementById('add-config-section')?.scrollIntoView({ behavior: 'smooth' });
                }}
                className="mt-4 px-4 py-2 rounded-md bg-green-600 text-white hover:bg-green-700 inline-flex items-center"
              >
                <Plus className="w-4 h-4 mr-2" />
                {t('connections.addYourFirst', 'Add Your First Config')}
              </button>
            )}
          </div>
        ) : (
          <div className="overflow-hidden rounded-lg border border-gray-200 dark:border-gray-700">
            <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
              <thead className={styles.tableHeader}>
                <tr>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    {t('common.key', 'Key')}
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    {t('common.value', 'Value')}
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    {t('common.type', 'Type')}
                  </th>
                  {editing && (
                    <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                      {t('users.actions')}
                    </th>
                  )}
                </tr>
              </thead>
              <tbody className={`divide-y ${styles.tableDivider}`}>
                {configs.map((config) => (
                  <tr key={config.key} className={styles.tableRow}>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">
                      {config.key}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      {renderConfigValue(config)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <span className={`px-2 py-1 rounded-full text-xs ${config.is_secret
                        ? theme === 'dark' ? 'bg-red-900 bg-opacity-30 text-red-300' : 'bg-red-100 text-red-800'
                        : theme === 'dark' ? 'bg-blue-900 bg-opacity-30 text-blue-300' : 'bg-blue-100 text-blue-800'
                        }`}>
                        {config.is_secret ? t('common.secret', 'Secret') : t('common.plain', 'Plain')}
                      </span>
                    </td>
                    {editing && (
                      <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                        <button
                          onClick={() => handleRemoveConfig(config.key)}
                          className="text-red-500 hover:text-red-700"
                        >
                          <Trash2 className="w-5 h-5" />
                        </button>
                      </td>
                    )}
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {/* Add Configuration Form */}
        {editing && (
          <div id="add-config-section" className={`mt-8 p-6 rounded-lg border ${theme === 'dark' ? 'bg-gray-700 border-gray-600' : 'bg-gray-50 border-gray-200'}`}>
            <h3 className="text-lg font-medium mb-4 flex items-center">
              <Plus className="w-5 h-5 mr-2" />
              {t('common.add', 'Add Configuration')}
            </h3>

            <div className="grid grid-cols-12 gap-4">
              <div className="col-span-5">
                <label className="block text-sm font-medium mb-1">{t('common.key', 'Key')}</label>
                <input
                  type="text"
                  value={newConfigKey}
                  onChange={(e) => setNewConfigKey(e.target.value)}
                  className={`w-full px-3 py-2 rounded-md border ${styles.input}`}
                  placeholder="API_KEY"
                />
              </div>

              <div className="col-span-5">
                <label className="block text-sm font-medium mb-1">{t('common.value', 'Value')}</label>
                <div className="relative">
                  <input
                    type={newConfigIsSecret ? 'password' : 'text'}
                    value={newConfigValue}
                    onChange={(e) => setNewConfigValue(e.target.value)}
                    className={`w-full px-3 py-2 rounded-md border ${styles.input} ${newConfigIsSecret ? 'pr-10' : ''}`}
                    placeholder={t('common.value', 'Value')}
                  />
                  {newConfigIsSecret && (
                    <button
                      type="button"
                      onClick={() => setNewConfigIsSecret(!newConfigIsSecret)}
                      className="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-500"
                    >
                      <Eye size={18} />
                    </button>
                  )}
                </div>
              </div>

              <div className="col-span-2 flex items-center pt-6">
                <label className="flex items-center">
                  <input
                    type="checkbox"
                    checked={newConfigIsSecret}
                    onChange={(e) => setNewConfigIsSecret(e.target.checked)}
                    className="mr-2 h-4 w-4"
                  />
                  <span className="text-sm">{t('common.secret', 'Secret')}</span>
                </label>
              </div>
            </div>

            <button
              onClick={handleAddNewConfig}
              disabled={!newConfigKey.trim()}
              className={`mt-4 px-4 py-2 rounded-md text-sm font-medium ${!newConfigKey.trim()
                ? 'bg-gray-400 cursor-not-allowed text-gray-200'
                : styles.primaryButton
                }`}
            >
              {t('connections.addConnection', 'Add Configuration')}
            </button>
          </div>
        )}

        {/* Test Result */}
        {testResult && (
          <div className={`mt-6 p-4 rounded-lg border ${testResult.success ? styles.success : styles.error}`}>
            <div className="flex items-start">
              <div className={`p-2 rounded-full mr-3 ${testResult.success
                ? 'bg-green-100 dark:bg-green-900'
                : 'bg-red-100 dark:bg-red-900'}`}>
                {testResult.success
                  ? <CheckCircle className="w-5 h-5 text-green-600 dark:text-green-400" />
                  : <XCircle className="w-5 h-5 text-red-600 dark:text-red-400" />}
              </div>
              <div>
                <h3 className="font-medium text-base">
                  {testResult.success
                    ? t('common.success', 'Connection Successful')
                    : t('common.error', 'Connection Failed')}
                </h3>
                <p className="mt-1 text-sm">{testResult.message}</p>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ConnectionDetailsConfiguration;
