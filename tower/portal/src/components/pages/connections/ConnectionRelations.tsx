import React, { JSX } from 'react';
import {
  Edit,
  Trash2,
  Plus,
  Download,
  Upload,
  ArrowRight,
  ExternalLink,
  Loader2
} from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';

import { useTheme } from '../../../context/ThemeContext';
import { getThemeStyles } from '../../../utility/theme';
import { useConnection } from '../../../context/ConnectionContext';
import RelationModal from './RelationModal';

const ConnectionRelations: React.FC = () => {
  const { t } = useTranslation('pages');
  const navigate = useNavigate();
  const { theme } = useTheme();
  const styles = getThemeStyles(theme);

  const {
    relationConnections,
    relationsLoading,
    relationsError,
    isAddModalOpen,
    setIsAddModalOpen,
    addRelationConnection,
    deleteRelationConnection,
    formatDate
  } = useConnection();

  const editRelation = (relationId: string): void => {
    navigate(`/connections/edit/${relationId}`);
  };

  const getRelationTypeIcon = (type: string): JSX.Element => {
    switch (type) {
      case 'inbound':
        return <Download size={18} className="text-blue-500" />;
      case 'outbound':
        return <Upload size={18} className="text-green-500" />;
      case 'bidirectional':
        return <ArrowRight size={18} className="text-purple-500" />;
      default:
        return <ExternalLink size={18} className="text-gray-500" />;
    }
  };

  return (
    <div className={`rounded-lg shadow-sm ${styles.card} overflow-hidden`}>
      <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
        <h2 className="text-lg font-medium">{t('connections.connectedIntegrations', 'Relation Connections')}</h2>
        <button
          className={`px-3 py-1 rounded-md text-sm flex items-center ${styles.primaryButton}`}
          onClick={() => setIsAddModalOpen(true)}
        >
          <Plus className="w-4 h-4 mr-1" />
          {t('connections.addConnection', 'New Relation')}
        </button>
      </div>

      {/* Add Relation Modal */}
      <RelationModal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        onAdd={addRelationConnection}
      />

      {relationsLoading ? (
        <div className="p-8 text-center">
          <Loader2 className="animate-spin h-8 w-8 mx-auto text-green-500 mb-2" />
          <p className="text-gray-500 dark:text-gray-400">{t('connections.loading', 'Loading relation connections...')}</p>
        </div>
      ) : relationsError ? (
        <div className={`p-6 m-4 rounded-lg border ${styles.error}`}>
          <p>{relationsError}</p>
          <button
            onClick={() => setIsAddModalOpen(true)}
            className="mt-2 px-3 py-1 rounded-md bg-gray-600 text-white hover:bg-gray-700 text-sm"
          >
            {t('common.retry', 'Retry')}
          </button>
        </div>
      ) : relationConnections.length === 0 ? (
        <div className="p-8 text-center">
          <p className="text-gray-500 dark:text-gray-400">{t('connections.noConnections', 'No relation connections found.')}</p>
          <button
            className="mt-4 px-4 py-2 rounded-md bg-green-600 text-white hover:bg-green-700"
            onClick={() => setIsAddModalOpen(true)}
          >
            {t('connections.addYourFirst', 'Add Your First Relation')}
          </button>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className={styles.tableHeader}>
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  {t('users.name', 'Name')}
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  {t('mappings.type', 'Type')}
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  {t('mappings.endpoint', 'Endpoint')}
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  {t('dashboard.lastSync', 'Last Used')}
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  {t('users.status', 'Status')}
                </th>
                <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  {t('users.actions', 'Actions')}
                </th>
              </tr>
            </thead>
            <tbody className={`divide-y ${styles.tableDivider}`}>
              {relationConnections.map((relation) => (
                <tr
                  onClick={() => editRelation(relation.id)}
                  key={relation.id}
                  className={`${styles.tableRow} cursor-pointer`}
                >
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <div className={`p-2 rounded-full mr-2 ${relation.type === 'inbound'
                        ? 'bg-blue-100 dark:bg-blue-900/30'
                        : relation.type === 'outbound'
                          ? 'bg-green-100 dark:bg-green-900/30'
                          : 'bg-purple-100 dark:bg-purple-900/30'
                        }`}>
                        {getRelationTypeIcon(relation.type)}
                      </div>
                      <div>
                        <div className="text-sm font-medium">
                          {relation.name}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`px-2 py-1 rounded-full text-xs ${relation.type === 'inbound'
                      ? styles.inbound
                      : relation.type === 'outbound'
                        ? styles.outbound
                        : styles.bidirectional
                      }`}>
                      {relation.type}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">
                    {relation.endpoint}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                    {formatDate(relation.lastUsed)}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`px-2 py-1 rounded-full text-xs ${relation.status === 'active' ? styles.active : styles.inactive
                      }`}>
                      {relation.status === 'active'
                        ? t('connectionList.statusActive', 'Active')
                        : t('connectionList.statusInactive', 'Inactive')}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <div className="flex justify-end space-x-2">
                      <button
                        className="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300"
                        onClick={(e) => {
                          e.stopPropagation();
                          deleteRelationConnection(relation.id);
                        }}
                      >
                        <Trash2 size={18} />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
};

export default ConnectionRelations;
