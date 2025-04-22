import React, { useState, useEffect } from 'react';
import { X, Loader2 } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { useTheme } from '../../../context/ThemeContext';
import { getThemeStyles } from '../../../utility/theme';
import { connectionService, RelationConnection } from '../../../services/connectionService';

interface RelationModalProps {
  isOpen: boolean;
  onClose: () => void;
  onAdd: (relationData: RelationConnection) => Promise<void>;
}

const RelationModal: React.FC<RelationModalProps> = ({
  isOpen,
  onClose,
  onAdd
}) => {
  const { t } = useTranslation('pages');
  const { theme } = useTheme();
  const styles = getThemeStyles(theme);

  const [targetId, setTargetId] = useState<string>('');
  const [connectionType, setConnectionType] = useState<string>('outbound');
  const [endpoint, setEndpoint] = useState<string>('');
  const [active, setActive] = useState<boolean>(true);
  const [availableConnections, setAvailableConnections] = useState<{ id: string; name: string }[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isOpen) {
      loadConnections();
      resetForm();
    }
  }, [isOpen]);

  const loadConnections = async (): Promise<void> => {
    try {
      setLoading(true);
      const connections = await connectionService.getApiConnections();
      setAvailableConnections(connections.map(conn => ({
        id: conn.id,
        name: conn.name
      })));
      setLoading(false);
    } catch (err) {
      console.error('Failed to load available connections:', err);
      setError('Failed to load available connections');
      setLoading(false);
    }
  };

  const resetForm = (): void => {
    setTargetId('');
    setConnectionType('outbound');
    setEndpoint('');
    setActive(true);
    setError(null);
  };

  const handleSubmit = async (): Promise<void> => {
    if (!targetId || !connectionType) {
      setError('Please select a target connection and connection type');
      return;
    }

    try {
      setSubmitting(true);
      setError(null);

      const relationData: RelationConnection = {
        initiator_id: targetId,
        connection_type: connectionType,
        endpoint: endpoint,
        active: active,
        target_id: '' // This would typically be set by the server or from the context
      };

      await onAdd(relationData);
      onClose();
    } catch (err) {
      console.error('Failed to add relation:', err);
      setError('Failed to add relation connection');
    } finally {
      setSubmitting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 backdrop-blur-sm bg-black bg-opacity-50">
      <div className={`rounded-lg shadow-xl ${styles.card} w-full max-w-md overflow-hidden`}>
        <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
          <h2 className="text-lg font-medium">{t('connections.addRelation', 'Add Relation Connection')}</h2>
          <button onClick={onClose} className="text-gray-500 cursor-pointer hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="p-6">
          {loading ? (
            <div className="text-center py-4">
              <Loader2 className="animate-spin h-8 w-8 mx-auto text-green-500 mb-2" />
              <p className="text-gray-500 dark:text-gray-400">{t('connections.loading', 'Loading available connections...')}</p>
            </div>
          ) : (
            <div className="space-y-4">
              {error && (
                <div className={`p-3 rounded-md ${styles.error}`}>
                  <p className="text-sm">{error}</p>
                </div>
              )}

              <div>
                <label className="block text-sm font-medium mb-1">
                  {t('connections.selectTarget', 'Select Target Connection')}
                </label>
                <select
                  value={targetId}
                  onChange={(e) => setTargetId(e.target.value)}
                  className={`w-full px-3 py-2 rounded-md border ${styles.input}`}
                >
                  <option value="">{t('common.selectPlaceholder', 'Select...')}</option>
                  {availableConnections.map(conn => (
                    <option key={conn.id} value={conn.id}>
                      {conn.name} ({conn.id})
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium mb-1">
                  {t('mappings.type', 'Connection Type')}
                </label>
                <select
                  value={connectionType}
                  onChange={(e) => setConnectionType(e.target.value)}
                  className={`w-full px-3 py-2 rounded-md border ${styles.input}`}
                >
                  <option value="outbound">{t('mappings.outbound', 'Outbound')}</option>
                  <option value="inbound">{t('mappings.inbound', 'Inbound')}</option>
                  <option value="bidirectional">{t('mappings.bidirectional', 'Bidirectional')}</option>
                </select>
              </div>

              <div>
                <label className="block text-sm font-medium mb-1">
                  {t('mappings.endpoint', 'Endpoint')} ({t('common.optional', 'Optional')})
                </label>
                <input
                  type="text"
                  value={endpoint}
                  onChange={(e) => setEndpoint(e.target.value)}
                  placeholder="/api/webhook"
                  className={`w-full px-3 py-2 rounded-md border ${styles.input}`}
                />
              </div>

              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="active-toggle"
                  checked={active}
                  onChange={(e) => setActive(e.target.checked)}
                  className="mr-2 h-4 w-4"
                />
                <label htmlFor="active-toggle" className="text-sm">
                  {t('connectionList.statusActive', 'Active')}
                </label>
              </div>
            </div>
          )}
        </div>

        <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex justify-end space-x-2">
          <button
            onClick={onClose}
            className={`px-4 py-2 rounded-md text-sm font-medium cursor-pointer ${styles.secondaryButton}`}
          >
            {t('common.cancel', 'Cancel')}
          </button>
          <button
            onClick={handleSubmit}
            disabled={submitting || loading}
            className={`px-4 py-2 rounded-md text-sm font-medium cursor-pointer ${styles.primaryButton}`}
          >
            {submitting ? (
              <>
                <Loader2 className="w-4 h-4 mr-2 inline animate-spin" />
                {t('common.adding', 'Adding...')}
              </>
            ) : (
              t('common.add', 'Add')
            )}
          </button>
        </div>
      </div>
    </div>
  );
};

export default RelationModal;
