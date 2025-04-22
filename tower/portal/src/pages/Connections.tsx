import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../context/ThemeContext';
import { Link } from 'react-router-dom';
import { Plus, Search, Loader2, RefreshCw, AlertCircle } from 'lucide-react';
import ConnectionList from '../components/pages/connections/ConnectionList';
import { ApiConnection, connectionService } from '../services/connectionService';

const Connections: React.FC = () => {
  const { t } = useTranslation('pages');
  const { theme } = useTheme();
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [connections, setConnections] = useState<ApiConnection[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchConnections();
  }, []);

  const fetchConnections = async () => {
    try {
      setLoading(true);
      const data = await connectionService.getApiConnections();
      setConnections(data);
      setError(null);
    } catch (err) {
      console.error('Failed to fetch connections', err);
      setError(t('connections.fetchError'));
    } finally {
      setLoading(false);
    }
  };

  const filteredConnections = connections.filter(connection =>
    connection.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    connection.type.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="p-6">
      <div className="mb-6 flex flex-col sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
            {t('navigation.connections')}
          </h1>
          <p className={`mt-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-600'}`}>
            {t('connections.description')}
          </p>
        </div>
        <div className="mt-4 sm:mt-0 flex space-x-3">
          <button
            onClick={fetchConnections}
            className="inline-flex items-center cursor-pointer px-4 py-2 border rounded-md shadow-sm text-sm font-medium bg-gray-100 text-gray-700 hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600 dark:border-gray-600"
          >
            <RefreshCw size={16} className="mr-2" />
            {t('connections.refresh')}
          </button>

          <Link to="/integrations">
            <button
              className="inline-flex items-center cursor-pointer px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
            >
              <Plus size={16} className="mr-2" />
              {t('connections.addNew')}
            </button>
          </Link>
        </div>
      </div>

      {/* Search */}
      <div className={`p-4 mb-6 rounded-lg ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'} shadow-md`}>
        <div className="relative max-w-xs">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Search size={18} className={theme === 'dark' ? 'text-gray-400' : 'text-gray-500'} />
          </div>
          <input
            type="text"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className={`block w-full pl-10 pr-3 py-2 border rounded-md ${theme === 'dark'
              ? 'bg-gray-700 border-gray-600 text-white placeholder-gray-400 focus:border-green-500 focus:ring-green-500'
              : 'border-gray-300 placeholder-gray-500 focus:border-green-500 focus:ring-green-500'
              }`}
            placeholder={t('connections.searchPlaceholder')}
          />
        </div>
      </div>

      {/* Loading State */}
      {loading ? (
        <div className="flex justify-center items-center py-12">
          <Loader2 className="w-8 h-8 animate-spin text-green-500 mr-2" />
          <span className={`text-lg ${theme === 'dark' ? 'text-gray-300' : 'text-gray-600'}`}>
            {t('connections.loading')}
          </span>
        </div>
      ) : error ? (
        <div className={`p-6 rounded-lg mb-6 ${theme === 'dark' ? 'bg-red-900/30 text-red-300 border border-red-800' : 'bg-red-50 text-red-800 border border-red-200'}`}>
          <div className="flex items-center">
            <AlertCircle className="h-5 w-5 mr-2" />
            <span>{error}</span>
          </div>
          <button
            onClick={fetchConnections}
            className="mt-3 inline-flex items-center px-3 py-1.5 text-sm rounded-md bg-gray-600 text-white hover:bg-gray-700"
          >
            <RefreshCw size={14} className="mr-1.5" />
            {t('common.retry')}
          </button>
        </div>
      ) : (
        <>
          {/* Connections */}
          <ConnectionList connections={filteredConnections} />

          {/* No Results */}
          {filteredConnections.length === 0 && (
            <div className={`text-center py-12 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
              {connections.length === 0 ? (
                <>
                  <p className="text-lg font-medium">{t('connections.noConnections')}</p>
                  <p className="mt-1">{t('connections.addYourFirst')}</p>
                  <Link to="/integrations">
                    <button className="mt-4 px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 inline-flex items-center">
                      <Plus size={16} className="mr-2" />
                      {t('connections.addConnection')}
                    </button>
                  </Link>
                </>
              ) : (
                <>
                  <p className="text-lg font-medium">{t('connections.noSearchResults')}</p>
                  <p className="mt-1">{t('connections.tryDifferentSearch')}</p>
                </>
              )}
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default Connections;
