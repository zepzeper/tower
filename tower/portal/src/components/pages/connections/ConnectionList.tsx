import React, { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../../../context/ThemeContext';
import { Edit, Trash2, Loader, MoreVertical, ExternalLink, Power, PowerOff, Plus } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';
import { ApiConnection, connectionService } from '../../../services/connectionService';

type ConnectionListProps = {
  connections: ApiConnection[];
};

const ConnectionList: React.FC<ConnectionListProps> = ({ connections }) => {
  const { t } = useTranslation('pages');
  const { theme } = useTheme();
  const navigate = useNavigate();

  const [isDeleting, setIsDeleting] = useState<string | null>(null);
  const [isToggling, setIsToggling] = useState<string | null>(null);
  const [activeDropdown, setActiveDropdown] = useState<string | null>(null);

  const formatDate = (date: string) =>
    new Date(date).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' });

  const handleDelete = async (id: string) => {
    setIsDeleting(id);
    try {
      await connectionService.deleteApiConnection(id);
    } catch (err) {
      console.error(err);
    } finally {
      setIsDeleting(null);
    }
  };

  const handleToggleActive = async (id: string, current: boolean) => {
    setIsToggling(id);
    try {
      await connectionService.updateApiConnection({ id, active: !current });
    } catch (err) {
      console.error(err);
    } finally {
      setIsToggling(null);
    }
  };

  const handleView = (connection: ApiConnection) => navigate(`/connections/view/${connection.id}`);

  const toggleDropdown = (id: string) => setActiveDropdown(prev => (prev === id ? null : id));

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {connections.map(connection => (
        <div
          key={connection.id}
          className={`p-4 rounded-lg border shadow cursor-pointer hover:shadow-md transition ${theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-200'}`}
          onClick={() => handleView(connection)}
        >
          <div className="flex items-start justify-between">
            <div className="flex items-center">
              <div>
                <h3 className={`text-lg font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                  {connection.name}
                </h3>
                <p className={`text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                </p>
              </div>
            </div>

            <div className="relative dropdown-parent">
              <button
                onClick={(e) => { e.stopPropagation(); toggleDropdown(connection.id); }}
                className={`p-1 rounded-md ${theme === 'dark' ? 'hover:bg-gray-700 text-gray-400' : 'hover:bg-gray-100 text-gray-600'}`}
              >
                <MoreVertical size={20} />
              </button>

              {activeDropdown === connection.id && (
                <div className={`absolute right-0 mt-1 w-48 rounded-md shadow-lg z-10 ${theme === 'dark' ? 'bg-gray-700' : 'bg-white'} ring-1 ring-black ring-opacity-5`}>
                  <div className="py-1">
                    <button onClick={(e) => { e.stopPropagation(); handleToggleActive(connection.id, connection.active); }}
                      disabled={isToggling === connection.id}
                      className="flex items-center w-full px-4 py-2 text-sm hover:bg-gray-100 dark:hover:bg-gray-600">
                      {isToggling === connection.id
                        ? <Loader size={16} className="animate-spin mr-2" />
                        : connection.active
                          ? <PowerOff size={16} className="mr-2" />
                          : <Power size={16} className="mr-2" />}
                      {connection.active ? t('connectionList.deactivate') : t('connectionList.activate')}
                    </button>
                    <button onClick={(e) => { e.stopPropagation(); handleDelete(connection.id); }}
                      disabled={isDeleting === connection.id}
                      className="flex items-center w-full px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-600">
                      {isDeleting === connection.id
                        ? <Loader size={16} className="animate-spin mr-2" />
                        : <Trash2 size={16} className="mr-2" />}
                      {t('connectionList.delete')}
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>

          {connection.description && (
            <p className={`mt-2 text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-600'}`}>
              {connection.description}
            </p>
          )}

          <div className="mt-3 pt-3 border-t border-gray-200 dark:border-gray-700 text-xs text-gray-500 dark:text-gray-400">
            <div className="flex justify-between items-center">
              <span>{t('connectionList.created')}: {formatDate(connection.created_at)}</span>
              <span className={`px-2 py-0.5 rounded font-medium ${connection.active
                ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300'
                : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'}`}>
                {connection.active ? t('connectionList.statusActive') : t('connectionList.statusInactive')}
              </span>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default ConnectionList;
