import React, { createContext, useContext, useState, useEffect, useMemo, ReactNode } from 'react';
import { connectionService, ApiConnection, ApiConnectionConfig, ConnectionType, RelationConnection, RelationConnectionLogs } from '../services/connectionService';

// Types
interface TestResult {
  success: boolean;
  message: string;
}

interface RelationConnectionUI {
  id: string;
  name: string;
  type: string;
  endpoint: string;
  lastUsed: string;
  status: string;
}

interface ConnectionContextType {
  // Connection data
  connection: ApiConnection | null;
  configs: ApiConnectionConfig[];
  connectionType: ConnectionType | null;
  connectionTypes: ConnectionType[];

  // UI states
  loading: boolean;
  error: string | null;
  editing: boolean;
  saving: boolean;
  testResult: TestResult | null;
  activeTab: string;
  showSecrets: Record<string, boolean>;

  // Relation connections data
  relationConnections: RelationConnectionUI[];
  relationsLoading: boolean;
  relationsError: string | null;
  isAddModalOpen: boolean;

  // Activity logs data
  recentActivity: RelationConnectionLogs[];
  activityLoading: boolean;
  activityError: string | null;
  selectedActivity: RelationConnectionLogs | null;
  isActivityModalOpen: boolean;

  // Setters and handlers
  setConnection: React.Dispatch<React.SetStateAction<ApiConnection | null>>;
  setConfigs: React.Dispatch<React.SetStateAction<ApiConnectionConfig[]>>;
  setActiveTab: React.Dispatch<React.SetStateAction<string>>;
  setEditing: React.Dispatch<React.SetStateAction<boolean>>;
  setIsAddModalOpen: React.Dispatch<React.SetStateAction<boolean>>;

  // Connection handlers
  handleUpdateConnection: () => Promise<void>;
  handleTestConnection: () => Promise<void>;
  handleInitiateOAuth: () => Promise<void>;
  toggleShowSecret: (key: string) => void;

  // Config handlers
  handleAddConfig: (key: string, value: string, isSecret: boolean) => void;
  handleRemoveConfig: (key: string) => void;
  updateConfigValue: (key: string, value: string) => void;

  // Relation handlers
  addRelationConnection: (relationData: RelationConnection) => Promise<void>;
  deleteRelationConnection: (relationId: string) => Promise<void>;

  // Activity handlers
  handleViewActivityDetails: (activity: RelationConnectionLogs) => void;
  closeActivityModal: () => void;

  // Utility
  formatDate: (dateString: string) => string;
}

// Create the context
const ConnectionContext = createContext<ConnectionContextType | null>(null);

// Custom hook to use the connection context
export const useConnection = (): ConnectionContextType => {
  const context = useContext(ConnectionContext);
  if (!context) {
    throw new Error('useConnection must be used within a ConnectionProvider');
  }
  return context;
};

// Format date utility function
export const formatDate = (dateString: string): string => {
  const date = new Date(dateString);
  return date.toLocaleString();
};

// Provider props type
interface ConnectionProviderProps {
  children: ReactNode;
  connectionId?: string;
}

// Provider component
export const ConnectionProvider: React.FC<ConnectionProviderProps> = ({
  children,
  connectionId
}) => {
  // Connection data states
  const [connection, setConnection] = useState<ApiConnection | null>(null);
  const [configs, setConfigs] = useState<ApiConnectionConfig[]>([]);
  const [connectionTypes, setConnectionTypes] = useState<ConnectionType[]>([]);
  const [connectionType, setConnectionType] = useState<ConnectionType | null>(null);

  // UI states
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [editing, setEditing] = useState<boolean>(false);
  const [saving, setSaving] = useState<boolean>(false);
  const [showSecrets, setShowSecrets] = useState<Record<string, boolean>>({});
  const [testResult, setTestResult] = useState<TestResult | null>(null);
  const [activeTab, setActiveTab] = useState<string>('details');

  // Cached data states with proper fetching logic
  const [relationConnections, setRelationConnections] = useState<RelationConnectionUI[]>([]);
  const [recentActivity, setRecentActivity] = useState<RelationConnectionLogs[]>([]);
  const [relationsLoading, setRelationsLoading] = useState<boolean>(true);
  const [activityLoading, setActivityLoading] = useState<boolean>(true);
  const [relationsError, setRelationsError] = useState<string | null>(null);
  const [activityError, setActivityError] = useState<string | null>(null);

  // Modal states
  const [isAddModalOpen, setIsAddModalOpen] = useState<boolean>(false);
  const [isActivityModalOpen, setIsActivityModalOpen] = useState<boolean>(false);
  const [selectedActivity, setSelectedActivity] = useState<RelationConnectionLogs | null>(null);

  // Fetch connection details
  useEffect(() => {
    if (!connectionId) return;

    const fetchDetails = async (): Promise<void> => {
      try {
        setLoading(true);
        const data = await connectionService.getApiConnectionWithConfig(connectionId);
        setConnection(data.connection);
        setConfigs(data.configs);

        // Load connection types to get the auth type
        const types = await connectionService.getConnectionTypesFromFile();
        setConnectionTypes(types);
        const foundType = types.find(t => t.id === data.connection.type);
        setConnectionType(foundType || null);
      } catch (err) {
        console.error(err);
        setError('Failed to load connection details.');
      } finally {
        setLoading(false);
      }
    };

    fetchDetails();
  }, [connectionId]);

  // Fetch relation connections only when tab is active
  useEffect(() => {
    if (activeTab !== 'relations' || !connectionId) return;

    const fetchRelationConnections = async (): Promise<void> => {
      try {
        setRelationsLoading(true);
        const response = await connectionService.getRelationConnections(connectionId);

        // Get connections for names
        const connections = await connectionService.getApiConnections();

        // Transform API response to UI model
        const uiConnections = response.map(relation => {
          const matchingConnection = connections.find(conn => conn.id === relation.initiator_id);
          return {
            id: relation.initiator_id,
            name: matchingConnection?.name || relation.target_id,
            type: relation.connection_type,
            endpoint: relation.endpoint,
            lastUsed: new Date().toISOString(), // This would come from the API in a real implementation
            status: relation.active ? 'active' : 'inactive'
          };
        });

        setRelationConnections(uiConnections);
        setRelationsError(null);
      } catch (err) {
        console.error('Failed to fetch relation connections:', err);
        setRelationsError('Failed to load relation connections');
      } finally {
        setRelationsLoading(false);
      }
    };

    fetchRelationConnections();
  }, [connectionId, activeTab]);

  // Fetch activity logs only when tab is active
  useEffect(() => {
    if (activeTab !== 'activity' || !connectionId) return;

    const fetchActivityLogs = async (): Promise<void> => {
      try {
        setActivityLoading(true);
        const result = await connectionService.getRelationConnectionLogs(connectionId);
        setRecentActivity(result);
        setActivityError(null);
      } catch (err) {
        console.error(err);
        setActivityError('Failed to load activity logs');
      } finally {
        setActivityLoading(false);
      }
    };

    fetchActivityLogs();
  }, [connectionId, activeTab]);

  // Update connection handler
  const handleUpdateConnection = async (): Promise<void> => {
    if (!connection) return;

    try {
      setSaving(true);
      const updateData = {
        id: connection.id,
        name: connection.name,
        description: connection.description || undefined,
        active: connection.active,
        configs: configs.map(c => ({
          key: c.key,
          value: c.value || '',
          is_secret: c.is_secret
        }))
      };

      const updated = await connectionService.updateApiConnection(updateData);
      setConnection(updated);
      setEditing(false);
      setTestResult(null);
    } catch (err) {
      console.error(err);
      setError('Failed to update connection.');
    } finally {
      setSaving(false);
    }
  };

  // Test connection handler
  const handleTestConnection = async (): Promise<void> => {
    if (!connection) return;

    try {
      setTestResult(null);
      const result = await connectionService.testApiConnection(connection.id);
      setTestResult(result);
    } catch (err) {
      console.error(err);
      setTestResult({
        success: false,
        message: 'Failed to test connection.'
      });
    }
  };

  // OAuth flow handler
  const handleInitiateOAuth = async (): Promise<void> => {
    if (!connectionType) return;

    try {
      const { url } = await connectionService.initiateOAuthFlow(connectionType.id);
      window.location.href = url;
    } catch (err) {
      console.error(err);
      setError('Failed to initiate OAuth flow.');
    }
  };

  // Toggle secret visibility
  const toggleShowSecret = (key: string): void => {
    setShowSecrets(prev => ({
      ...prev,
      [key]: !prev[key]
    }));
  };

  // Relation connection handlers
  const addRelationConnection = async (relationData: RelationConnection): Promise<void> => {
    try {
      await connectionService.addRelationConnection(relationData);

      // Refetch relations after adding
      if (connectionId) {
        setRelationsLoading(true);
        const response = await connectionService.getRelationConnections(connectionId);
        const connections = await connectionService.getApiConnections();

        const uiConnections = response.map(relation => {
          const matchingConnection = connections.find(conn => conn.id === relation.initiator_id);
          return {
            id: relation.initiator_id,
            name: matchingConnection?.name || relation.target_id,
            type: relation.connection_type,
            endpoint: relation.endpoint,
            lastUsed: new Date().toISOString(),
            status: relation.active ? 'active' : 'inactive'
          };
        });

        setRelationConnections(uiConnections);
        setRelationsLoading(false);
      }
    } catch (err) {
      console.error('Failed to add relation connection:', err);
      setRelationsError('Failed to add relation connection');
      throw err;
    }
  };

  const deleteRelationConnection = async (relationId: string): Promise<void> => {
    try {
      await connectionService.deleteRelationConnection(connectionId!, relationId);

      // Remove the deleted relation from state
      setRelationConnections(prev => prev.filter(relation => relation.id !== relationId));
    } catch (err) {
      console.error('Failed to delete relation connection:', err);
      setRelationsError('Failed to delete relation connection');
    }
  };

  // Config handlers
  const handleAddConfig = (key: string, value: string, isSecret: boolean): void => {
    if (!key.trim()) return;

    const newConfig: ApiConnectionConfig = {
      connection_id: connection?.id || '',
      key: key,
      value: value,
      is_secret: isSecret,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    };

    setConfigs([...configs, newConfig]);
  };

  const handleRemoveConfig = (key: string): void => {
    setConfigs(configs.filter(c => c.key !== key));
  };

  const updateConfigValue = (key: string, value: string): void => {
    setConfigs(configs.map(c =>
      c.key === key ? { ...c, value } : c
    ));
  };

  // Activity log handlers
  const handleViewActivityDetails = (activity: RelationConnectionLogs): void => {
    setSelectedActivity(activity);
    setIsActivityModalOpen(true);
  };

  const closeActivityModal = (): void => {
    setIsActivityModalOpen(false);
    setSelectedActivity(null);
  };

  // Context value memoization to prevent unnecessary re-renders
  const contextValue = useMemo<ConnectionContextType>(() => ({
    // Connection data
    connection,
    configs,
    connectionType,
    connectionTypes,

    // UI states
    loading,
    error,
    editing,
    saving,
    testResult,
    activeTab,
    showSecrets,

    // Relation connections data
    relationConnections,
    relationsLoading,
    relationsError,
    isAddModalOpen,

    // Activity logs data
    recentActivity,
    activityLoading,
    activityError,
    selectedActivity,
    isActivityModalOpen,

    // Setters and handlers
    setConnection,
    setConfigs,
    setActiveTab,
    setEditing,
    setIsAddModalOpen,

    // Connection handlers
    handleUpdateConnection,
    handleTestConnection,
    handleInitiateOAuth,
    toggleShowSecret,

    // Config handlers
    handleAddConfig,
    handleRemoveConfig,
    updateConfigValue,

    // Relation handlers
    addRelationConnection,
    deleteRelationConnection,

    // Activity handlers
    handleViewActivityDetails,
    closeActivityModal,

    // Utility
    formatDate
  }), [
    connection, configs, connectionType, connectionTypes,
    loading, error, editing, saving, testResult, activeTab, showSecrets,
    relationConnections, relationsLoading, relationsError, isAddModalOpen,
    recentActivity, activityLoading, activityError, selectedActivity, isActivityModalOpen
  ]);

  return (
    <ConnectionContext.Provider value={contextValue}>
      {children}
    </ConnectionContext.Provider>
  );
};
