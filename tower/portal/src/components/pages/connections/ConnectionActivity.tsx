import React from 'react';
import { Clock, Loader2, X } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { useTheme } from '../../../context/ThemeContext';
import { getThemeStyles } from '../../../utility/theme';
import { useConnection } from '../../../context/ConnectionContext';

const ConnectionActivity: React.FC = () => {
  const { t } = useTranslation('pages');
  const { theme } = useTheme();
  const styles = getThemeStyles(theme);

  const {
    recentActivity,
    activityLoading,
    activityError,
    selectedActivity,
    isActivityModalOpen,
    handleViewActivityDetails,
    closeActivityModal,
    formatDate
  } = useConnection();

  return (
    <>
      <div className={`rounded-lg shadow-sm ${styles.card} overflow-hidden`}>
        <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
          <h2 className="text-lg font-medium">{t('dashboard.recentActivity')}</h2>
          <div className="flex items-center space-x-2 text-sm text-gray-500 dark:text-gray-400">
            <Clock className="w-4 h-4 mr-1" />
            {t('analytics.lastDay', 'Last 24 hours')}
          </div>
        </div>

        {activityLoading ? (
          <div className="p-8 text-center">
            <Loader2 className="animate-spin h-8 w-8 mx-auto text-green-500 mb-2" />
            <p className="text-gray-500 dark:text-gray-400">{t('connections.loading', 'Loading activity logs...')}</p>
          </div>
        ) : activityError ? (
          <div className={`p-6 m-4 rounded-lg border ${styles.error}`}>
            <p>{activityError}</p>
            <button
              onClick={() => {/* Refresh activity logs */ }}
              className="mt-2 px-3 py-1 rounded-md bg-gray-600 text-white hover:bg-gray-700 text-sm"
            >
              {t('common.retry', 'Retry')}
            </button>
          </div>
        ) : recentActivity.length === 0 ? (
          <div className="p-8 text-center">
            <p className="text-gray-500 dark:text-gray-400">{t('connections.noResults', 'No recent activity found.')}</p>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
              <thead className={styles.tableHeader}>
                <tr>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    {t('mappings.type', 'Type')}
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    {t('common.message', 'Message')}
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    {t('common.initiator', 'Initiator')}
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    {t('common.target', 'Target')}
                  </th>
                  <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                    {t('dashboard.timestamp', 'Date')}
                  </th>
                </tr>
              </thead>
              <tbody className={`divide-y ${styles.tableDivider}`}>
                {recentActivity.map((activity) => (
                  <tr
                    onClick={() => handleViewActivityDetails(activity)}
                    key={activity.id || `${activity.initiator_id}-${activity.created_at}`}
                    className={`${styles.tableRow} cursor-pointer`}
                  >
                    <td className="px-6 py-4 whitespace-nowrap">
                      <div className="flex items-center">
                        <span className={`px-2 py-1 rounded-full text-xs ${activity.connection_type === 'inbound'
                          ? styles.inbound
                          : activity.connection_type === 'outbound'
                            ? styles.outbound
                            : styles.bidirectional
                          }`}>
                          {activity.connection_type}
                        </span>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600 dark:text-gray-300">
                      {activity.message}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">
                      {activity.initiator_id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">
                      {activity.target_id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                      {formatDate(activity.created_at)}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 text-center">
          <button className={`text-sm font-medium ${theme === 'dark' ? 'text-green-400 hover:text-green-300' : 'text-green-600 hover:text-green-700'}`}>
            {t('common.viewAll', 'View All Activity')}
          </button>
        </div>
      </div>

      {/* Activity Detail Modal */}
      {isActivityModalOpen && selectedActivity && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 backdrop-blur-sm bg-opacity-30">
          <div className={`rounded-lg shadow-xl ${styles.card} w-full max-w-3xl max-h-[90vh] overflow-hidden`}>
            <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
              <h2 className="text-lg font-medium">{t('connectionDetails.activityDetail', 'Activity Detail')}</h2>
              <button onClick={closeActivityModal} className="text-gray-500 cursor-pointer hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="p-6 overflow-y-auto max-h-[calc(90vh-120px)]">
              <div className="space-y-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400">{t('mappings.type', 'Type')}</h3>
                    <div className="mt-1">
                      <span className={`px-2 py-1 rounded-full text-xs ${selectedActivity.connection_type === 'inbound'
                        ? styles.inbound
                        : selectedActivity.connection_type === 'outbound'
                          ? styles.outbound
                          : styles.bidirectional
                        }`}>
                        {selectedActivity.connection_type}
                      </span>
                    </div>
                  </div>

                  <div>
                    <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400">{t('dashboard.timestamp', 'Date')}</h3>
                    <p className="mt-1 text-sm text-gray-700 dark:text-gray-300">{formatDate(selectedActivity.created_at)}</p>
                  </div>

                  <div>
                    <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400">{t('common.initiator', 'Initiator')}</h3>
                    <p className="mt-1 text-sm font-mono break-all">{selectedActivity.initiator_id}</p>
                  </div>

                  <div>
                    <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400">{t('common.target', 'Target')}</h3>
                    <p className="mt-1 text-sm font-mono break-all">{selectedActivity.target_id}</p>
                  </div>
                </div>

                <div>
                  <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400">{t('common.message', 'Message')}</h3>
                  <p className="mt-1 text-sm text-gray-700 dark:text-gray-300">{selectedActivity.message}</p>
                </div>

                {selectedActivity && (
                  <div>
                    <h3 className="text-sm font-medium text-gray-500 dark:text-gray-400">{t('common.metadata', 'Metadata')}</h3>
                    <pre className={`mt-1 p-3 rounded-md text-xs overflow-x-auto ${theme === 'dark' ? 'bg-gray-900' : 'bg-gray-100'}`}>
                      {JSON.stringify(selectedActivity, null, 2)}
                    </pre>
                  </div>
                )}
              </div>
            </div>

            <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex justify-end">
              <button
                onClick={closeActivityModal}
                className={`px-4 py-2 rounded-md text-sm font-medium cursor-pointer ${styles.secondaryButton}`}
              >
                {t('common.close', 'Close')}
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default ConnectionActivity;
