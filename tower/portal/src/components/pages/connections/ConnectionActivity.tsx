import React, { useEffect, useState } from 'react';
import {
  Clock,
  CheckCircle,
  XCircle,
  Download,
  Upload,
  ArrowRight,
  ExternalLink
} from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { connectionService, RelationConnectionLogs } from '../../../services/connectionService';

interface ConnectionActivityProps {
  connectionId: string;
  theme: string;
  formatDate: (dateString: string) => string;
}

const ConnectionActivity: React.FC<ConnectionActivityProps> = ({
  connectionId,
  theme,
  formatDate
}) => {
  const { t } = useTranslation();
  const [recentActivity, setRecentActivity] = useState<RelationConnectionLogs[]>([]);

  useEffect(() => {
    handleLogs();
  }, [connectionId]);

  const getRelationTypeIcon = (type: string) => {
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


  const handleLogs = async () => {

    try {
      const result = await connectionService.getRelationConnectionLogs(connectionId);
      setRecentActivity(result)
    } catch (err) {
      console.error(err);
    }
  };


  return (
    <div className={`rounded-lg shadow-sm ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'} overflow-hidden`}>
      <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
        <h2 className="text-lg font-medium">Recent Activity</h2>
        <div className="flex items-center space-x-2 text-sm text-gray-500 dark:text-gray-400">
          <Clock className="w-4 h-4 mr-1" />
          Last 24 hours
        </div>
      </div>

      {recentActivity.length === 0 ? (
        <div className="p-8 text-center">
          <p className="text-gray-500 dark:text-gray-400">No recent activity found.</p>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className={theme === 'dark' ? 'bg-gray-700' : 'bg-gray-50'}>
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Type
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Message
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Initiator
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Target
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Date
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
              {recentActivity.map((activity) => (
                <tr key={activity.id} className={theme === 'dark' ? 'bg-gray-800 hover:bg-gray-750' : 'bg-white hover:bg-gray-50'}>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <div className="flex items-center">
                      <div className={`p-2 rounded-full mr-2 ${activity.connection_type === 'inbound'
                        ? 'bg-blue-100 dark:bg-blue-900/30'
                        : activity.connection_type === 'outbound'
                          ? 'bg-green-100 dark:bg-green-900/30'
                          : 'bg-purple-100 dark:bg-purple-900/30'
                        }`}>
                        {getRelationTypeIcon(activity.connection_type)}
                      </div>
                      <span className={`px-2 py-1 rounded-full text-xs ${activity.connection_type === 'inbound'
                        ? 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:bg-opacity-30 dark:text-blue-300'
                        : activity.connection_type === 'outbound'
                          ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:bg-opacity-30 dark:text-green-300'
                          : 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:bg-opacity-30 dark:text-purple-300'
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
          View All Activity
        </button>
      </div>
    </div>
  );
};

export default ConnectionActivity;
