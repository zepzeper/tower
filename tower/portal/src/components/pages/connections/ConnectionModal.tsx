import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../../../context/ThemeContext';
import { X, Loader, Eye, EyeOff, ExternalLink } from 'lucide-react';
import { connectionService, ConnectionType, ApiConnectionCreateRequest } from '../../../services/connectionService';

interface ConnectionModalProps {
  isOpen: boolean;
  connectionId?: string;
  connectionType?: string;
  onClose: () => void;
  onSuccess: () => void;
}

interface ConfigField {
  key: string;
  name: string;
  description: string;
  type: 'string' | 'boolean' | 'number' | 'secret';
  required: boolean;
  value: string;
  isSecret: boolean;
}

const ConnectionModal: React.FC<ConnectionModalProps> = ({
  isOpen,
  connectionId,
  connectionType: initialType,
  onClose,
  onSuccess
}) => {
  const { t } = useTranslation('components');
  const { theme } = useTheme();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [selectedType, setSelectedType] = useState(initialType || '');
  const [configFields, setConfigFields] = useState<ConfigField[]>([]);
  const [connectionTypes, setConnectionTypes] = useState<ConnectionType[]>([]);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [testing, setTesting] = useState(false);
  const [showSecrets, setShowSecrets] = useState<Record<string, boolean>>({});
  const [testResult, setTestResult] = useState<{ status: 'success' | 'error' | null; message: string }>({ status: null, message: '' });

  const isEditMode = !!connectionId;
  const darkModeClasses = theme === 'dark' ? 'bg-gray-800 text-white border-gray-700' : 'bg-white text-gray-900 border-gray-200';
  const currentAuthType = connectionTypes.find(type => type.id === selectedType)?.authType || 'api_key';

  useEffect(() => {
    if (!isOpen) return;

    const loadConnectionTypes = async () => {
      try {
        const types = await connectionService.getConnectionTypesFromFile();
        setConnectionTypes(types);
        if (!isEditMode && !initialType && types.length > 0) {
          setSelectedType(types[0].id);
        }
      } catch (error) {
        console.error('Failed to load connection types:', error);
      }
    };

    loadConnectionTypes();
  }, [isOpen, initialType, isEditMode]);

  useEffect(() => {
    if (!isOpen || !isEditMode) return;

    const loadConnection = async () => {
      setLoading(true);
      try {
        const { connection, configs } = await connectionService.getApiConnectionWithConfig(connectionId);
        setName(connection.name);
        setDescription(connection.description || '');
        setSelectedType(connection.type);
      } catch (error) {
        console.error(`Failed to load connection ${connectionId}:`, error);
        setTestResult({ status: 'error', message: t('connectionModal.loadError') });
      } finally {
        setLoading(false);
      }
    };

    loadConnection();
  }, [isOpen, connectionId, isEditMode, t]);

  useEffect(() => {
    if (!selectedType) return;

    const loadConfigFields = async () => {
      const typeInfo = connectionTypes.find(type => type.id === selectedType);
      if (!typeInfo) return;

      if (isEditMode) {
        try {
          const { connection, configs } = await connectionService.getApiConnectionWithConfig(connectionId);
          const newConfigFields = typeInfo.configTemplate.map(template => ({
            ...template,
            value: configs.find(c => c.key === template.key)?.value || template.default || '',
            isSecret: template.type === 'secret'
          }));
          setConfigFields(newConfigFields);
        } catch (error) {
          console.error(`Failed to load connection configs for ${connectionId}:`, error);
        }
      } else {
        const newConfigFields = typeInfo.configTemplate.map(template => ({
          ...template,
          value: template.default || '',
          isSecret: template.type === 'secret'
        }));
        setName(typeInfo.name)
        setConfigFields(newConfigFields);
      }
    };

    loadConfigFields();
  }, [selectedType, connectionTypes, isEditMode, connectionId]);

  const handleOAuthConnect = async () => {
    try {
      const { url } = await connectionService.initiateOAuthFlow(selectedType);
      window.open(url, '_blank', 'width=600,height=600');
      // You'll need to implement a way to listen for the OAuth callback
      // This typically involves polling or websockets
    } catch (error) {
      setTestResult({ status: 'error', message: error instanceof Error ? error.message : 'Failed to initiate OAuth flow' });
    }
  };

  const handleSaveConnection = async () => {
    if (!isFormValid()) {
      setTestResult({ status: 'error', message: t('connectionModal.missingRequiredFields') });
      return;
    }

    setSaving(true);
    try {
      const connectionData: ApiConnectionCreateRequest = {
        name,
        description: description || undefined,
        type: selectedType,
        configs: configFields.map(field => ({
          key: field.key,
          value: field.value,
          is_secret: field.isSecret
        }))
      };

      if (isEditMode) {
        await connectionService.updateApiConnection({ id: connectionId, ...connectionData });
      } else {
        await connectionService.createApiConnection(connectionData);
      }

      onSuccess();
      onClose();
    } catch (error) {
      setTestResult({
        status: 'error',
        message: isEditMode ? t('connectionModal.updateError') : t('connectionModal.createError')
      });
    } finally {
      setSaving(false);
    }
  };

  const isFormValid = (): boolean => {
    if (!selectedType) return false;

    // Check if all required fields have values
    const requiredFieldsValid = configFields
      .filter(field => field.required)
      .every(field => field.value.trim() !== '');


    return requiredFieldsValid;
  };

  const handleUpdateConfigField = (key: string, value: string): void => {
    setConfigFields(prev =>
      prev.map(field =>
        field.key === key ? { ...field, value } : field
      )
    );
  };

  const toggleShowSecret = (key: string): void => {
    setShowSecrets(prev => ({
      ...prev,
      [key]: !prev[key]
    }));
  };

  // Add this test connection handler as well if you need it
  const handleTestConnection = async (): Promise<void> => {
    if (!isFormValid()) {
      setTestResult({
        status: 'error',
        message: t('connectionModal.missingRequiredFields')
      });
      return;
    }

    setTesting(true);
    setTestResult({ status: null, message: '' });

    try {
      const connectionData: ApiConnectionCreateRequest = {
        name: name,
        description: description || undefined,
        type: selectedType,
        configs: configFields.map(field => ({
          key: field.key,
          value: field.value,
          is_secret: field.isSecret
        }))
      };

      const result = await connectionService.testApiConnection(connectionData);
      setTestResult({
        status: result.success ? 'success' : 'error',
        message: result.message
      });
    } catch (error) {
      setTestResult({
        status: 'error',
        message: t('connectionModal.testError')
      });
    } finally {
      setTesting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 backdrop-blur-sm bg-opacity-30">
      <div className={`w-full max-w-xl rounded-lg shadow-lg border ${darkModeClasses} max-h-[90vh] flex flex-col`}>
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-xl font-semibold">
            {isEditMode ? t('connectionModal.editConnection') : t('connectionModal.addConnection')}
          </h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
            <X size={24} />
          </button>
        </div>

        {/* Body */}
        <div className="p-6 overflow-y-auto flex-grow">
          {loading ? (
            <div className="flex justify-center py-8">
              <Loader size={32} className="animate-spin text-green-500" />
            </div>
          ) : (
            <div className="space-y-6">
              {/* Connection Type (only when creating) */}
              {!isEditMode && (
                <div>
                  <label className={`block mb-2 text-sm font-medium ${theme === 'dark' ? 'text-gray-200' : 'text-gray-700'}`}>
                    {t('connectionModal.typeLabel')} *
                  </label>
                  <select
                    value={selectedType}
                    onChange={(e) => setSelectedType(e.target.value)}
                    className={`w-full px-3 py-2 border rounded-md ${theme === 'dark'
                      ? 'bg-gray-700 border-gray-600 text-white'
                      : 'border-gray-300 text-gray-900'
                      }`}
                  >
                    {connectionTypes.map(type => (
                      <option key={type.id} value={type.id}>{type.name}</option>
                    ))}
                  </select>
                </div>
              )}

              {/* OAuth Connect Button */}
              {currentAuthType === 'oauth2' && !isEditMode && (
                <div>
                  <button
                    onClick={handleOAuthConnect}
                    className="flex items-center justify-center w-full px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                  >
                    <ExternalLink size={16} className="mr-2" />
                    {t('connectionModal.connectWithOAuth')}
                  </button>
                  <p className={`mt-2 text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                    {t('connectionModal.oauthDescription')}
                  </p>
                </div>
              )}

              {/* Configuration Fields */}
              {configFields.map((field) => (
                <div key={field.key}>
                  <label className={`block mb-2 text-sm font-medium ${theme === 'dark' ? 'text-gray-200' : 'text-gray-700'}`}>
                    {field.name} {field.required && '*'}
                  </label>
                  <div className="relative">
                    <input
                      type={field.isSecret && !showSecrets[field.key] ? 'password' : 'text'}
                      value={field.value}
                      onChange={(e) => handleUpdateConfigField(field.key, e.target.value)}
                      className={`w-full px-3 py-2 border rounded-md ${theme === 'dark'
                        ? 'bg-gray-700 border-gray-600 text-white placeholder-gray-400'
                        : 'border-gray-300 placeholder-gray-500'
                        } ${field.isSecret ? 'pr-10' : ''}`}
                      placeholder={field.description}
                    />
                    {field.isSecret && (
                      <button
                        type="button"
                        onClick={() => toggleShowSecret(field.key)}
                        className="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-500"
                      >
                        {showSecrets[field.key] ? <EyeOff size={18} /> : <Eye size={18} />}
                      </button>
                    )}
                  </div>
                </div>
              ))}

              {/* Test Result */}
              {testResult.status && (
                <div className={`p-3 rounded-md ${testResult.status === 'success'
                  ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
                  : 'bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300'
                  }`}>
                  {testResult.message}
                </div>
              )}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex justify-between p-4 border-t border-gray-200 dark:border-gray-700">
          <button
            onClick={handleTestConnection}
            disabled={saving || loading}
            className="px-4 py-2 text-sm font-medium text-white bg-gray-600 rounded-md hover:bg-gray-700 disabled:opacity-50"
          >
            {saving ? <Loader size={16} className="animate-spin mr-2" /> : null}
            {t('connectionModal.testConnection')}
          </button>
          <button
            onClick={handleSaveConnection}
            disabled={saving || loading}
            className="px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700 disabled:opacity-50"
          >
            {saving ? <Loader size={16} className="animate-spin mr-2" /> : null}
            {isEditMode ? t('connectionModal.updateConnection') : t('connectionModal.saveConnection')}
          </button>
        </div>
      </div>
    </div>
  );
};

export default ConnectionModal;
