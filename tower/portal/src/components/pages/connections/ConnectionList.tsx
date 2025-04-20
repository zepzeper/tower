import React, { useState, useRef, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../../../context/ThemeContext';
import { Edit, Trash2, Loader, MoreVertical, ExternalLink, Power, PowerOff } from 'lucide-react';
import { ApiConnection } from '../../../services/ConnectionService';

interface ConnectionListProps {
  connections: ApiConnection[];
  onEdit: (connection: ApiConnection) => void;
  onDelete: (connectionId: string) => Promise<void>;
  onToggleActive: (connectionId: string, active: boolean) => Promise<void>;
}

const ConnectionList: React.FC<ConnectionListProps> = ({
  connections,
  onEdit,
  onDelete,
  onToggleActive
}) => {
  const { t } = useTranslation('components');
  const { theme } = useTheme();
  const [isDeleting, setIsDeleting] = useState<string | null>(null);
  const [isToggling, setIsToggling] = useState<string | null>(null);
  const [activeDropdown, setActiveDropdown] = useState<string | null>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setActiveDropdown(null);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  // Dictionary of connection types
  const connectionTypeIcons: Record<string, string> = {
    woocommerce: '🛒',
    shopify: '🏪',
    hubspot: '📊',
    brincr: '🤝',
    quickbooks: '💵',
    mailchimp: '📧',
    salesforce: '☁️',
    xero: '📒',
  };

  const connectionTypeNames: Record<string, string> = {
    woocommerce: 'WooCommerce',
    shopify: 'Shopify',
    hubspot: 'HubSpot',
    brincr: 'Brincr',
    quickbooks: 'QuickBooks',
    mailchimp: 'Mailchimp',
    salesforce: 'Salesforce',
    xero: 'Xero',
  };

  // Format date for display
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString(undefined, {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  // Handle connection deletion
  const handleDelete = async (connectionId: string) => {
    setIsDeleting(connectionId);
    try {
      await onDelete(connectionId);
    } finally {
      setIsDeleting(null);
    }
  };

  // Handle toggling connection active state
  const handleToggleActive = async (connectionId: string, currentActive: boolean) => {
    setIsToggling(connectionId);
    try {
      await onToggleActive(connectionId, !currentActive);
    } finally {
      setIsToggling(null);
    }
  };

  // Toggle dropdown menu
  const toggleDropdown = (connectionId: string) => {
    setActiveDropdown(prev => prev === connectionId ? null : connectionId);
  };

  if (connections.length === 0) {
    return (
      <div className={`text-center py-6 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
        <p>{t('connectionList.noConnections')}</p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {connections.map((connection) => (
        <div
          key={connection.id}
          className={`p-4 rounded-lg border shadow ${theme === 'dark'
            ? 'bg-gray-800 border-gray-700'
            : 'bg-white border-gray-200'
            }`}
        >
          <div className="flex items-start justify-between">
            <div className="flex items-center">
              <div className="flex-shrink-0 w-10 h-10 flex items-center justify-center text-xl rounded-lg bg-gray-100 dark:bg-gray-700 mr-3">
                {connectionTypeIcons[connection.type] || '🔌'}
              </div>
              <div>
                <h3 className={`text-lg font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                  {connection.name}
                </h3>
                <p className={`text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                  {connectionTypeNames[connection.type] || connection.type}
                </p>
              </div>
            </div>

            <div className="relative" ref={dropdownRef}>
              <button
                onClick={() => toggleDropdown(connection.id)}
                className={`p-1 rounded-md ${theme === 'dark'
                  ? 'hover:bg-gray-700 text-gray-400'
                  : 'hover:bg-gray-100 text-gray-600'
                  }`}
              >
                <MoreVertical size={20} />
              </button>

              {activeDropdown === connection.id && (
                <div className={`absolute right-0 mt-1 w-48 rounded-md shadow-lg z-10 ${theme === 'dark' ? 'bg-gray-700' : 'bg-white'
                  } ring-1 ring-black ring-opacity-5`}>
                  <div className="py-1">
                    <button
                      onClick={() => onEdit(connection)}
                      className={`flex items-center w-full px-4 py-2 text-sm ${theme === 'dark'
                        ? 'text-gray-200 hover:bg-gray-600'
                        : 'text-gray-700 hover:bg-gray-100'
                        }`}
                    >
                      <Edit size={16} className="mr-2" />
                      {t('connectionList.edit')}
                    </button>

                    <button
                      onClick={() => handleToggleActive(connection.id, connection.active)}
                      disabled={isToggling === connection.id}
                      className={`flex items-center w-full px-4 py-2 text-sm ${theme === 'dark'
                        ? 'text-gray-200 hover:bg-gray-600'
                        : 'text-gray-700 hover:bg-gray-100'
                        }`}
                    >
                      {isToggling === connection.id ? (
                        <Loader size={16} className="animate-spin mr-2" />
                      ) : connection.active ? (
                        <PowerOff size={16} className="mr-2" />
                      ) : (
                        <Power size={16} className="mr-2" />
                      )}
                      {connection.active
                        ? t('connectionList.deactivate')
                        : t('connectionList.activate')}
                    </button>

                    <button
                      onClick={() => handleDelete(connection.id)}
                      disabled={isDeleting === connection.id}
                      className={`flex items-center w-full px-4 py-2 text-sm ${theme === 'dark'
                        ? 'text-red-400 hover:bg-gray-600'
                        : 'text-red-600 hover:bg-gray-100'
                        }`}
                    >
                      {isDeleting === connection.id ? (
                        <Loader size={16} className="animate-spin mr-2" />
                      ) : (
                        <Trash2 size={16} className="mr-2" />
                      )}
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

          <div className="mt-3 pt-3 border-t border-gray-200 dark:border-gray-700">
            <div className="flex justify-between items-center">
              <div className={`text-xs ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                {t('connectionList.created')}: {formatDate(connection.created_at)}
              </div>

              <div className="flex items-center">
                <span
                  className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${connection.active
                    ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:bg-opacity-30 dark:text-green-300'
                    : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'
                    }`}
                >
                  {connection.active
                    ? t('connectionList.statusActive')
                    : t('connectionList.statusInactive')}
                </span>
              </div>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default ConnectionList;
