import React, { useState, useEffect } from 'react';
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
import { connectionService, RelationConnection } from '../../../services/connectionService';
import RelationModal from './RelationModal';

interface RelationConnectionUI {
  id: string;
  name: string;
  type: string;
  endpoint: string;
  lastUsed: string;
  status: string;
}

interface ConnectionRelationsProps {
  connectionId: string;
  theme: string;
  formatDate: (dateString: string) => string;
}

const ConnectionRelations: React.FC<ConnectionRelationsProps> = ({
  connectionId,
  theme,
  formatDate
}) => {
  const [relationConnections, setRelationConnections] = useState<RelationConnectionUI[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [isAddModalOpen, setIsAddModalOpen] = useState<boolean>(false);
  const [addingRelation, setAddingRelation] = useState<boolean>(false);

  useEffect(() => {
    fetchRelationConnections();
  }, [connectionId]);

  const fetchRelationConnections = async () => {
    try {
      setLoading(true);
      const response = await connectionService.getRelationConnections(connectionId);

      // Transform API response to UI model if needed
      const uiConnections: RelationConnectionUI[] = response.map(relation => ({
        id: relation.initiator_id,
        name: relation.target_id, // We'll need to fetch the actual name from connections
        type: relation.connection_type,
        endpoint: relation.endpoint,
        lastUsed: new Date().toISOString(), // This would come from the API in a real implementation
        status: relation.active ? 'active' : 'inactive'
      }));

      // If we have connection names, let's update them
      const connections = await connectionService.getApiConnections();

      // Update the connection names
      const updatedConnections = uiConnections.map(relation => {
        const matchingConnection = connections.find(conn => conn.id === relation.id);
        return {
          ...relation,
          name: matchingConnection?.name || relation.id
        };
      });

      setRelationConnections(updatedConnections);
      setError(null);
    } catch (err) {
      console.error('Failed to fetch relation connections:', err);
      setError('Failed to load relation connections');
      setRelationConnections([]);
    } finally {
      setLoading(false);
    }
  };

  const addRelationConnection = async (relationData: RelationConnection) => {
    try {
      setAddingRelation(true);
      await connectionService.addRelationConnection(relationData);
      await fetchRelationConnections(); // Refresh data after adding
    } catch (err) {
      console.error('Failed to add relation connection:', err);
      setError('Failed to add relation connection');
      throw err; // Propagate error to modal
    } finally {
      setAddingRelation(false);
    }
  };

  const deleteRelationConnection = async (relationId: string) => {
    try {
      await connectionService.deleteRelationConnection(connectionId, relationId);
      await fetchRelationConnections(); // Refresh data after deletion
    } catch (err) {
      console.error('Failed to delete relation connection:', err);
      setError('Failed to delete relation connection');
    }
  };

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

  return (
    <div className={`rounded-lg shadow-sm ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'} overflow-hidden`}>
      <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
        <h2 className="text-lg font-medium">Relation Connections</h2>
        <button
          className={`px-3 py-1 rounded-md text-sm flex items-center ${theme === 'dark'
            ? 'bg-green-600 hover:bg-green-700'
            : 'bg-green-500 hover:bg-green-600 text-white'
            }`}
          onClick={() => setIsAddModalOpen(true)}
          disabled={addingRelation}
        >
          {addingRelation ? (
            <Loader2 className="w-4 h-4 mr-1 animate-spin" />
          ) : (
            <Plus className="w-4 h-4 mr-1" />
          )}
          New Relation
        </button>
      </div>

      {/* Add Relation Modal */}
      <RelationModal
        connectionId={connectionId}
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        onAdd={addRelationConnection}
        theme={theme}
      />

      {loading ? (
        <div className="p-8 text-center">
          <Loader2 className="animate-spin h-8 w-8 mx-auto text-green-500 mb-2" />
          <p className="text-gray-500 dark:text-gray-400">Loading relation connections...</p>
        </div>
      ) : error ? (
        <div className={`p-6 m-4 rounded-lg border ${theme === 'dark' ? 'bg-red-900/30 border-red-800 text-red-300' : 'bg-red-50 border-red-200 text-red-600'}`}>
          <p>{error}</p>
          <button
            onClick={fetchRelationConnections}
            className="mt-2 px-3 py-1 rounded-md bg-gray-600 text-white hover:bg-gray-700 text-sm"
          >
            Retry
          </button>
        </div>
      ) : relationConnections.length === 0 ? (
        <div className="p-8 text-center">
          <p className="text-gray-500 dark:text-gray-400">No relation connections found.</p>
          <button
            className="mt-4 px-4 py-2 rounded-md bg-green-600 text-white hover:bg-green-700"
            onClick={() => setIsAddModalOpen(true)}
          >
            Add Your First Relation
          </button>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className={theme === 'dark' ? 'bg-gray-700' : 'bg-gray-50'}>
              <tr>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Name
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Type
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Endpoint
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Last Used
                </th>
                <th scope="col" className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Status
                </th>
                <th scope="col" className="px-6 py-3 text-right text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200 dark:divide-gray-700">
              {relationConnections.map((relation) => (
                <tr key={relation.id} className={theme === 'dark' ? 'bg-gray-800 hover:bg-gray-750' : 'bg-white hover:bg-gray-50'}>
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
                      ? 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:bg-opacity-30 dark:text-blue-300'
                      : relation.type === 'outbound'
                        ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:bg-opacity-30 dark:text-green-300'
                        : 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:bg-opacity-30 dark:text-purple-300'
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
                    <span className={`px-2 py-1 rounded-full text-xs ${relation.status === 'active'
                      ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:bg-opacity-30 dark:text-green-300'
                      : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-400'
                      }`}>
                      {relation.status.charAt(0).toUpperCase() + relation.status.slice(1)}
                    </span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                    <div className="flex justify-end space-x-2">
                      <button
                        className="text-green-600 hover:text-green-900 dark:text-green-400 dark:hover:text-green-300"
                        onClick={() => {/* Implement edit functionality */ }}
                      >
                        <Edit size={18} />
                      </button>
                      <button
                        className="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300"
                        onClick={() => deleteRelationConnection(relation.id)}
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
