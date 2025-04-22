import React from 'react';
import { useParams, useNavigate } from 'react-router-dom';
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

import { useTheme } from '../../context/ThemeContext';
import { getThemeStyles } from '../../utility/theme';
import { ConnectionProvider, useConnection } from '../../context/ConnectionContext';

// The main component just sets up the context provider
const ConnectionDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();

  return (
    <ConnectionProvider connectionId={id}>
      <ConnectionDetailsContent />
    </ConnectionProvider>
  );
};

// The inner component consumes the context
const ConnectionDetailsContent: React.FC = () => {
  const { t } = useTranslation('pages');
  const { theme } = useTheme();
  const navigate = useNavigate();
  const styles = getThemeStyles(theme);

  const {
    connection,
    loading,
    error,
    editing,
    saving,
    activeTab,
    setActiveTab,
    setEditing,
    setConnection,
    handleUpdateConnection,
    handleTestConnection
  } = useConnection();

  if (loading) {
    return (
      <div className="flex items-center justify-center h-screen">
        <Loader2 className="animate-spin h-10 w-10 text-green-500" />
        <span className="ml-2">{t('connectionDetails.loading')}</span>
      </div>
    );
  }

  if (error || !connection) {
    return (
      <div className="p-6">
        <div className={`p-6 mx-auto max-w-4xl rounded-lg border ${styles.error}`}>
          <h3 className="text-lg font-medium mb-2">{t('common.error')}</h3>
          <p>{error || t('connectionDetails.notFound')}</p>
          <button
            onClick={() => navigate('/connections')}
            className="mt-4 px-4 py-2 rounded-md bg-gray-600 text-white hover:bg-gray-700"
          >
            {t('connectionDetails.backToList')}
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className={`min-h-screen ${styles.pageBackground}`}>
      {/* Header */}
      <div className={`py-6 px-8 border-b ${styles.card} ${styles.border} shadow-sm`}>
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
                    className={`text-xl font-bold px-3 py-2 rounded-md border ${styles.input}`}
                  />
                ) : (
                  <span>{connection.name}</span>
                )}
                <span className={`ml-3 text-sm px-2 py-0.5 rounded-full ${connection.active ? styles.active : styles.inactive}`}>
                  {connection.active ? t('connectionDetails.statusActive') : t('connectionDetails.statusInactive')}
                </span>
              </h1>
              <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                {t('connectionDetails.type')}: {connection.type} • {t('connectionDetails.id')}: {connection.id}
              </p>
            </div>
          </div>

          <div className="flex space-x-3">
            {editing ? (
              <>
                <button
                  onClick={handleUpdateConnection}
                  disabled={saving}
                  className={`px-4 py-2 rounded-md ${styles.primaryButton} flex items-center cursor-pointer`}
                >
                  {saving ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : <Save className="w-4 h-4 mr-2" />}
                  {t('common.save')}
                </button>
                <button
                  onClick={() => setEditing(false)}
                  className={`px-4 py-2 rounded-md ${styles.secondaryButton} flex items-center cursor-pointer`}
                >
                  <X className="w-4 h-4 mr-2" />
                  {t('common.cancel')}
                </button>
              </>
            ) : (
              <>
                <button
                  onClick={handleTestConnection}
                  className="px-4 py-2 rounded-md bg-blue-600 text-white hover:bg-blue-700 flex items-center cursor-pointer"
                >
                  <TestTube2 className="w-4 h-4 mr-2" />
                  {t('connectionModal.testConnection')}
                </button>
                <button
                  onClick={() => setEditing(true)}
                  className={`px-4 py-2 rounded-md ${styles.primaryButton} flex items-center cursor-pointer`}
                >
                  <Edit className="w-4 h-4 mr-2" />
                  {t('connectionDetails.edit')}
                </button>
              </>
            )}
          </div>
        </div>
      </div>

      {/* Tabs Navigation */}
      <div className={`border-b ${styles.border}`}>
        <div className="flex px-8">
          <button
            onClick={() => setActiveTab('details')}
            className={`py-4 px-6 font-medium text-sm border-b-2 ${activeTab === 'details' ? styles.tab.active : styles.tab.inactive}`}
          >
            <Settings className="w-4 h-4 inline mr-2" />
            {t('connectionDetails.configuration')}
          </button>
          <button
            onClick={() => setActiveTab('relations')}
            className={`py-4 px-6 font-medium text-sm border-b-2 ${activeTab === 'relations' ? styles.tab.active : styles.tab.inactive}`}
          >
            <ArrowRight className="w-4 h-4 inline mr-2" />
            {t('connections.connectedIntegrations')}
          </button>
          <button
            onClick={() => setActiveTab('activity')}
            className={`py-4 px-6 font-medium text-sm border-b-2 ${activeTab === 'activity' ? styles.tab.active : styles.tab.inactive}`}
          >
            <Activity className="w-4 h-4 inline mr-2" />
            {t('dashboard.recentActivity')}
          </button>
        </div>
      </div>

      {/* Main Content Area */}
      <div className="p-8">
        {/* Details & Configuration Tab */}
        {activeTab === 'details' && <ConnectionDetailsConfiguration />}

        {/* Relations Tab */}
        {activeTab === 'relations' && <ConnectionRelations />}

        {/* Activity Tab */}
        {activeTab === 'activity' && <ConnectionActivity />}
      </div>
    </div>
  );
};

export default ConnectionDetails;
