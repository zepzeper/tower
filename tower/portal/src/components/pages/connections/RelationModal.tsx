import React, { useState, useEffect } from 'react';
import { X, Loader2 } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { connectionService, ApiConnection, RelationConnection } from '../../../services/connectionService';

interface AddRelationModalProps {
  connectionId: string;
  isOpen: boolean;
  onClose: () => void;
  onAdd: (relationData: Omit<RelationConnection, 'id'>) => Promise<void>;
  theme: string;
}

const AddRelationModal: React.FC<AddRelationModalProps> = ({
  connectionId,
  isOpen,
  onClose,
  onAdd,
  theme
}) => {
  const { t } = useTranslation('pages');
  const [activeConnections, setActiveConnections] = useState<ApiConnection[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState<boolean>(false);

  // Form state
  const [selectedConnection, setSelectedConnection] = useState<string>('');
  const [relationType, setRelationType] = useState<string>('outbound');
  const [endpoint, setEndpoint] = useState<string>('');

  useEffect(() => {
    if (isOpen) {
      fetchActiveConnections();
    }
  }, [isOpen]);

  const fetchActiveConnections = async () => {
    try {
      setLoading(true);
      const connections = await connectionService.getApiConnections();
      // Filter to only active connections and exclude the current connection
      const filteredConnections = connections.filter(
        conn => conn.active && conn.id !== connectionId
      );
      setActiveConnections(filteredConnections);
      setError(null);
    } catch (err) {
      console.error('Failed to fetch active connections:', err);
      setError(t('connections.fetchError', 'Failed to load active connections'));
      setActiveConnections([]);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!selectedConnection) {
      setError(t('connectionModal.missingRequiredFields', 'Please select a connection'));
      return;
    }

    try {
      setSubmitting(true);

      // Get the name of the selected connection for display purposes
      const selectedConnectionObj = activeConnections.find(conn => conn.id === selectedConnection);

      const relationData = {
        initiator_id: connectionId,
        target_id: selectedConnection,
        type: relationType,
        endpoint: endpoint,
        active: true
      };

      await onAdd(relationData);
      onClose();

      // Reset form
      setSelectedConnection('');
      setRelationType('outbound');
      setEndpoint('products');
      setError(null);
    } catch (err) {
      console.error('Error adding relation:', err);
      setError(t('connections.updateError', 'Failed to add relation connection'));
    } finally {
      setSubmitting(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 backdrop-blur-sm bg-opacity-30">
      <div className={`relative w-full max-w-md p-6 rounded-lg shadow-lg ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300"
        >
          <X size={20} />
        </button>

        <h2 className="text-xl font-semibold mb-4">{t('connections.addConnection', 'Add New Relation')}</h2>

        {error && (
          <div className={`p-3 mb-4 rounded-md ${theme === 'dark' ? 'bg-red-900/30 border border-red-800 text-red-300' : 'bg-red-50 border border-red-200 text-red-600'}`}>
            {error}
          </div>
        )}

        {loading ? (
          <div className="p-8 text-center">
            <Loader2 className="animate-spin h-8 w-8 mx-auto text-green-500 mb-2" />
            <p className="text-gray-500 dark:text-gray-400">{t('connections.loading', 'Loading connections...')}</p>
          </div>
        ) : (
          <form onSubmit={handleSubmit}>
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">
                {t('common.target', 'Target Connection')}
              </label>
              {activeConnections.length === 0 ? (
                <div className="p-3 rounded-md bg-yellow-50 border border-yellow-200 text-yellow-700 dark:bg-yellow-900/30 dark:border-yellow-800 dark:text-yellow-300">
                  {t('connections.noConnections', 'No active connections available to link. Please create and activate other connections first.')}
                </div>
              ) : (
                <select
                  value={selectedConnection}
                  onChange={(e) => setSelectedConnection(e.target.value)}
                  className={`w-full px-3 py-2 rounded-md border ${theme === 'dark'
                    ? 'bg-gray-700 border-gray-600 text-white'
                    : 'bg-white border-gray-300'}`}
                  required
                >
                  <option value="">{t('connectionModal.selectType', 'Select a connection')}</option>
                  {activeConnections.map((conn) => (
                    <option key={conn.id} value={conn.id}>
                      {conn.name} ({conn.type})
                    </option>
                  ))}
                </select>
              )}
            </div>

            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">
                {t('mappings.type', 'Relation Type')}
              </label>
              <select
                value={relationType}
                onChange={(e) => setRelationType(e.target.value)}
                className={`w-full px-3 py-2 rounded-md border ${theme === 'dark'
                  ? 'bg-gray-700 border-gray-600 text-white'
                  : 'bg-white border-gray-300'}`}
                required
              >
                <option value="outbound">Outbound</option>
                <option value="inbound">Inbound</option>
                <option value="bidirectional">Bidirectional</option>
              </select>
            </div>

            <div className="mb-6">
              <label className="block text-sm font-medium mb-1">
                {t('mappings.endpoint', 'Endpoint Path')}
              </label>
              <select
                value={endpoint}
                onChange={(e) => setEndpoint(e.target.value)}
                className={`w-full px-3 py-2 rounded-md border ${theme === 'dark'
                  ? 'bg-gray-700 border-gray-600 text-white'
                  : 'bg-white border-gray-300'}`}
                required
              >
                <option value="products">Products</option>
                <option value="orders">Orders</option>
                <option value="inventory">Inventory</option>
              </select>
            </div>

            <div className="flex justify-end space-x-3">
              <button
                type="button"
                onClick={onClose}
                className={`px-4 py-2 rounded-md ${theme === 'dark'
                  ? 'bg-gray-700 hover:bg-gray-600 text-white'
                  : 'bg-gray-200 hover:bg-gray-300 text-gray-800'}`}
              >
                {t('common.cancel')}
              </button>
              <button
                type="submit"
                disabled={submitting || activeConnections.length === 0}
                className={`px-4 py-2 rounded-md flex items-center ${theme === 'dark'
                  ? 'bg-green-600 hover:bg-green-700 text-white'
                  : 'bg-green-500 hover:bg-green-600 text-white'} ${(submitting || activeConnections.length === 0) ? 'opacity-50 cursor-not-allowed' : ''
                  }`}
              >
                {submitting && <Loader2 className="w-4 h-4 mr-2 animate-spin" />}
                {t('connections.addConnection', 'Add Relation')}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
};

export default AddRelationModal;
