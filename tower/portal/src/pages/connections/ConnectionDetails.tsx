import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  connectionService,
  ApiConnection,
  ApiConnectionConfig,
  ConnectionType,
  ApiConnectionUpdateRequest,
  RelationConnectionLogs
} from '../../services/connectionService';
import { useTheme } from '../../context/ThemeContext';
import { useTranslation } from 'react-i18next';
import {
  Loader2,
  ChevronLeft,
  Edit,
  Save,
  X,
  TestTube2,
  ArrowRight,
  Activity,
  Settings
} from 'lucide-react';

import ConnectionDetailsConfiguration from '../../components/pages/connections/ConnectionDetailsConfiguration';
import ConnectionRelations from '../../components/pages/connections/ConnectionRelations';
import ConnectionActivity from '../../components/pages/connections/ConnectionActivity';

interface TestResult {
  success: boolean;
  message: string;
}

interface Theme {
  theme: string;
}

const ConnectionDetails: React.FC = () => {
  const { t } = useTranslation();
  const { id } = useParams<{ id: string }>();
  const { theme } = useTheme() as Theme;
  const navigate = useNavigate();

  const [activeTab, setActiveTab] = useState<string>('details');
  const [connection, setConnection] = useState<ApiConnection | null>(null);
  const [configs, setConfigs] = useState<ApiConnectionConfig[]>([]);
  const [connectionTypes, setConnectionTypes] = useState<ConnectionType[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [editing, setEditing] = useState<boolean>(false);
  const [testResult, setTestResult] = useState<TestResult | null>(null);
  const [connectionType, setConnectionType] = useState<ConnectionType | null>(null);
  const [saving, setSaving] = useState<boolean>(false);
  const [showSecrets, setShowSecrets] = useState<Record<string, boolean>>({});

  // We'll keep these here but they'll be used by child components
  const [relationConnections, setRelationConnections] = useState<any[]>([]);

  useEffect(() => {
    if (!id) return;

    const fetchDetails = async () => {
      try {
        setLoading(true);
        const data = await connectionService.getApiConnectionWithConfig(id);
        setConnection(data.connection);
        setConfigs(data.configs);

        // Load connection types to get the auth type
        const types = await connectionService.getConnectionTypesFromFile();
        setConnectionTypes(types);
        const foundType = types.find(t => t.id === data.connection.type);
        setConnectionType(foundType || null);
      } catch (err) {
        console.error(err);
        setError('Failed to load connection details.');
      } finally {
        setLoading(false);
      }
    };

    fetchDetails();
  }, [id]);

  const handleUpdateConnection = async () => {
    if (!connection) return;

    try {
      setSaving(true);
      const updateData: ApiConnectionUpdateRequest = {
        id: connection.id,
        name: connection.name,
        description: connection.description || undefined,
        active: connection.active,
        configs: configs.map(c => ({
          key: c.key,
          value: c.value || '',
          is_secret: c.is_secret
        }))
      };

      const updated = await connectionService.updateApiConnection(updateData);
      setConnection(updated);
      setEditing(false);
      setTestResult(null);
    } catch (err) {
      console.error(err);
      setError('Failed to update connection.');
    } finally {
      setSaving(false);
    }
  };

  const handleTestConnection = async (id: string) => {
    if (!connection) return;

    try {
      setTestResult(null);
      const result = await connectionService.testApiConnection(id);
      setTestResult(result);
    } catch (err) {
      console.error(err);
      setTestResult({
        success: false,
        message: 'Failed to test connection.'
      });
    }
  };

  const handleInitiateOAuth = async () => {
    if (!connectionType) return;

    try {
      const { url } = await connectionService.initiateOAuthFlow(connectionType.id);
      window.location.href = url;
    } catch (err) {
      console.error(err);
      setError('Failed to initiate OAuth flow.');
    }
  };

  const toggleShowSecret = (key: string) => {
    setShowSecrets(prev => ({
      ...prev,
      [key]: !prev[key]
    }));
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <Loader2 className="animate-spin h-10 w-10 text-green-500" />
      </div>
    );
  }

  if (error || !connection) {
    return (
      <div className="p-6">
        <div className={`p-6 mx-auto max-w-4xl rounded-lg border ${theme === 'dark' ? 'bg-red-900/30 border-red-800 text-red-300' : 'bg-red-50 border-red-200 text-red-600'}`}>
          <h3 className="text-lg font-medium mb-2">Error</h3>
          <p>{error || 'Connection not found.'}</p>
          <button
            onClick={() => navigate('/connections')}
            className="mt-4 px-4 py-2 rounded-md bg-gray-600 text-white hover:bg-gray-700"
          >
            Back to Connections
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className={`min-h-screen ${theme === 'dark' ? 'bg-gray-900 text-white' : 'bg-gray-50 text-gray-900'}`}>
      {/* Header */}
      <div className={`py-6 px-8 border-b ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-200'} shadow-sm`}>
        <div className="flex items-center justify-between">
          <div className="flex items-center">
            <button
              onClick={() => navigate('/connections')}
              className="mr-4 p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-700"
            >
              <ChevronLeft className="w-5 h-5" />
            </button>
            <div>
              <h1 className="text-2xl font-bold flex items-center">
                {editing ? (
                  <input
                    type="text"
                    value={connection.name}
                    onChange={(e) => setConnection({ ...connection, name: e.target.value })}
                    className={`text-xl font-bold px-3 py-2 rounded-md border ${theme === 'dark'
                      ? 'bg-gray-700 border-gray-600 text-white'
                      : 'bg-white border-gray-300 text-gray-900'
                      }`}
                  />
                ) : (
                  <span>{connection.name}</span>
                )}
                <span className={`ml-3 text-sm px-2 py-0.5 rounded-full ${connection.active
                  ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
                  : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'
                  }`}>
                  {connection.active ? 'Active' : 'Inactive'}
                </span>
              </h1>
              <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                Type: {connection.type} • ID: {connection.id}
              </p>
            </div>
          </div>

          <div className="flex space-x-3">
            {editing ? (
              <>
                <button
                  onClick={handleUpdateConnection}
                  disabled={saving}
                  className="px-4 py-2 rounded-md bg-green-600 text-white hover:bg-green-700 flex items-center"
                >
                  {saving ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : <Save className="w-4 h-4 mr-2" />}
                  Save Changes
                </button>
                <button
                  onClick={() => {
                    setEditing(false);
                    // Reset to original values
                    if (id) {
                      connectionService.getApiConnectionWithConfig(id).then(data => {
                        setConnection(data.connection);
                        setConfigs(data.configs);
                      });
                    }
                  }}
                  className={`px-4 py-2 rounded-md ${theme === 'dark'
                    ? 'bg-gray-700 hover:bg-gray-600'
                    : 'bg-gray-200 hover:bg-gray-300'
                    } text-gray-800 dark:text-white flex items-center`}
                >
                  <X className="w-4 h-4 mr-2" />
                  Cancel
                </button>
              </>
            ) : (
              <>
                <button
                  onClick={() => { handleTestConnection(connection.id) }}
                  className="px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700 flex items-center"
                >
                  <TestTube2 className="w-4 h-4 mr-2" />
                  Test Connection
                </button>
                <button
                  onClick={() => setEditing(true)}
                  className="px-4 py-2 rounded-md bg-green-600 text-white hover:bg-green-700 flex items-center"
                >
                  <Edit className="w-4 h-4 mr-2" />
                  Edit
                </button>
              </>
            )}
          </div>
        </div>
      </div>

      {/* Tabs Navigation */}
      <div className={`border-b ${theme === 'dark' ? 'border-gray-700' : 'border-gray-200'}`}>
        <div className="flex px-8">
          <button
            onClick={() => setActiveTab('details')}
            className={`py-4 px-6 font-medium text-sm border-b-2 ${activeTab === 'details'
              ? theme === 'dark'
                ? 'border-green-500 text-green-400'
                : 'border-green-600 text-green-600'
              : theme === 'dark'
                ? 'border-transparent text-gray-400 hover:text-gray-300'
                : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
          >
            <Settings className="w-4 h-4 inline mr-2" />
            Details & Configuration
          </button>
          <button
            onClick={() => setActiveTab('relations')}
            className={`py-4 px-6 font-medium text-sm border-b-2 ${activeTab === 'relations'
              ? theme === 'dark'
                ? 'border-green-500 text-green-400'
                : 'border-green-600 text-green-600'
              : theme === 'dark'
                ? 'border-transparent text-gray-400 hover:text-gray-300'
                : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
          >
            <ArrowRight className="w-4 h-4 inline mr-2" />
            Relation Connections
          </button>
          <button
            onClick={() => setActiveTab('activity')}
            className={`py-4 px-6 font-medium text-sm border-b-2 ${activeTab === 'activity'
              ? theme === 'dark'
                ? 'border-green-500 text-green-400'
                : 'border-green-600 text-green-600'
              : theme === 'dark'
                ? 'border-transparent text-gray-400 hover:text-gray-300'
                : 'border-transparent text-gray-500 hover:text-gray-700'
              }`}
          >
            <Activity className="w-4 h-4 inline mr-2" />
            Recent Activity
          </button>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="p-8">
        {/* Details & Configuration Tab */}
        {activeTab === 'details' && (
          <ConnectionDetailsConfiguration
            connection={connection}
            configs={configs}
            connectionType={connectionType}
            editing={editing}
            theme={theme}
            testResult={testResult}
            setConnection={setConnection}
            setConfigs={setConfigs}
            handleInitiateOAuth={handleInitiateOAuth}
            formatDate={formatDate}
            toggleShowSecret={toggleShowSecret}
            showSecrets={showSecrets}
          />
        )}

        {/* Relations Tab */}
        {activeTab === 'relations' && (
          <ConnectionRelations
            connectionId={connection.id}
            theme={theme}
            formatDate={formatDate}
          />
        )}

        {/* Activity Tab */}
        {activeTab === 'activity' && (
          <ConnectionActivity
            connectionId={connection.id}
            theme={theme}
            formatDate={formatDate}
          />
        )}
      </div>
    </div>
  );
};

export default ConnectionDetails;
