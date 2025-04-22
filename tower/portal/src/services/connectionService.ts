import RequestHandler from '../handler/RequestHandler';
import CacheManager from './cacheService';

export interface ApiConnection {
  id: string;
  name: string;
  description: string | null;
  type: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}

export interface ApiConnectionConfig {
  connection_id: string;
  key: string;
  value: string | null;
  is_secret: boolean;
  created_at: string;
  updated_at: string;
}

export interface ApiConnectionCreateRequest {
  id: string;
  name: string;
  description?: string;
  type: string;
  configs: {
    key: string;
    value: string;
    is_secret: boolean;
  }[];
}

export interface ApiConnectionUpdateRequest {
  id: string;
  name?: string;
  description?: string;
  active?: boolean;
  configs?: {
    key: string;
    value: string;
    is_secret: boolean;
  }[];
}

export interface ConnectionType {
  id: string;
  name: string;
  description: string;
  category: string;
  categoryName: string;
  popular: boolean;
  imageUrl: string;
  authType: 'api_key' | 'oauth2' | 'basic_auth' | 'token';
  configTemplate: {
    key: string;
    name: string;
    description: string;
    type: 'string' | 'boolean' | 'number' | 'secret';
    required: boolean;
    default?: any;
  }[];
}

export interface RelationConnection {
  initiator_id: string;
  target_id: string;
  active: boolean;
  connection_type: string;
  endpoint: string;
}

export interface RelationConnectionLogs {
  id: number
  initiator_id: string;
  target_id: string;
  connection_type: string;
  message: string
  created_at: string;
}

class ConnectionService extends RequestHandler {
  constructor(cacheManager: CacheManager) {
    super(cacheManager);
  }

  async getApiConnections(forceRefresh = false): Promise<ApiConnection[]> {
    const cacheKey = 'connections:all';

    if (!forceRefresh && this.cacheManager.has(cacheKey)) {
      return this.cacheManager.get(cacheKey);
    }

    try {
      const data = await this.get<ApiConnection[]>('/connections/all');

      // Cache the result with connection IDs as dependencies
      const connectionIds = data.map(conn => `connection:${conn.id}`);
      this.cacheManager.set(cacheKey, data, ['connections', ...connectionIds]);

      return data;
    } catch (error) {
      console.error('Failed to fetch connections', error);
      throw error;
    }
  }

  async getApiConnectionWithConfig(id: string, forceRefresh = false): Promise<{
    connection: ApiConnection;
    configs: ApiConnectionConfig[];
  }> {
    const cacheKey = `connection:${id}:withConfig`;

    if (!forceRefresh && this.cacheManager.has(cacheKey)) {
      return this.cacheManager.get(cacheKey);
    }

    try {
      const data = await this.get<{
        connection: ApiConnection;
        configs: ApiConnectionConfig[];
      }>(`/connections/${id}`);

      // Cache with dependencies on both connection ID and configs
      this.cacheManager.set(cacheKey, data, [
        `connection:${id}`,
        `connection:${id}:configs`
      ]);

      return data;
    } catch (error) {
      console.error(`Failed to load connection ${id}:`, error);
      throw error;
    }
  }

  async createApiConnection(data: ApiConnectionCreateRequest): Promise<ApiConnection> {
    try {
      const result = await this.post<ApiConnection>('/connections', data);

      // Invalidate the connections list since we have a new connection
      this.cacheManager.invalidateResource('connections');

      return result;
    } catch (error) {
      console.error('Failed to create connection:', error);
      throw error;
    }
  }

  async updateApiConnection(data: ApiConnectionUpdateRequest): Promise<ApiConnection> {
    try {
      const result = await this.patch<ApiConnection>(`/connections/patch/${data.id}`, data);

      // Invalidate both the specific connection and the connections list
      this.cacheManager.invalidateResource(`connection:${data.id}`);

      // If configs were updated, invalidate config-specific cache
      if (data.configs) {
        this.cacheManager.invalidateResource(`connection:${data.id}:configs`);
      }

      return result;
    } catch (error) {
      console.error('Failed to update connection:', error);
      throw error;
    }
  }

  async deleteApiConnection(id: string): Promise<void> {
    try {
      await this.delete(`/connections/delete/${id}`);

      this.cacheManager.invalidateResource('connections');
    } catch (error) {
      console.error('Failed to delete connection:', error);
      throw error;
    }
  }

  async testApiConnection(id: string): Promise<{
    success: boolean;
    message: string;
  }> {
    return this.get<{ success: boolean; message: string }>(`/connections/test/${id}`);
  }


  async getConnectionTypesFromFile(forceRefresh = false): Promise<ConnectionType[]> {
    const cacheKey = 'connection:types';

    if (!forceRefresh && this.cacheManager.has(cacheKey)) {
      return this.cacheManager.get(cacheKey);
    }

    try {
      const res = await fetch('/schema/apis.json');
      if (!res.ok) throw new Error('Failed to load connection types');

      const data = (await res.json()) as ConnectionType[];

      // Cache with 'connection-types' as the resource identifier
      this.cacheManager.set(cacheKey, data, ['connection-types']);

      return data;
    } catch (error) {
      console.error('Failed to load connection types:', error);
      throw error;
    }
  }

  async initiateOAuthFlow(connectionType: string): Promise<{ url: string }> {
    return this.post<{ url: string }>(`/oauth/${connectionType}/initiate`, {});
  }

  async getRelationConnections(connectionId: string, forceRefresh = false): Promise<RelationConnection[]> {
    const cacheKey = `connection:${connectionId}:relations`;

    if (!forceRefresh && this.cacheManager.has(cacheKey)) {
      return this.cacheManager.get(cacheKey);
    }

    try {
      const data = await this.get<RelationConnection[]>(`/connections/${connectionId}/relations`);

      // Cache with dependencies on both the connection and its relations
      this.cacheManager.set(cacheKey, data, [
        `connection:${connectionId}`,
        'relations'
      ]);

      return data;
    } catch (error) {
      console.error(`Failed to get relation connections for ${connectionId}:`, error);
      throw error;
    }
  }

  async addRelationConnection(relationData: RelationConnection): Promise<RelationConnection> {
    try {
      const result = await this.post<RelationConnection>(
        `/connections/${relationData.initiator_id}/relations/create`,
        relationData
      );

      // Invalidate relations cache for both initiator and target
      this.cacheManager.invalidateResource(`connection:${relationData.initiator_id}`);
      this.cacheManager.invalidateResource(`connection:${relationData.target_id}`);
      this.cacheManager.invalidateResource('relations');

      return result;
    } catch (error) {
      console.error('Failed to add relation connection:', error);
      throw error;
    }
  }

  async updateRelationConnection(relationData: RelationConnection): Promise<RelationConnection> {
    try {
      const result = await this.patch<RelationConnection>(
        `/connections/${relationData.initiator_id}/relations/${relationData.target_id}`,
        relationData
      );

      // Invalidate relations cache for both initiator and target
      this.cacheManager.invalidateResource(`connection:${relationData.initiator_id}`);
      this.cacheManager.invalidateResource(`connection:${relationData.target_id}`);
      this.cacheManager.invalidateResource('relations');

      return result;
    } catch (error) {
      console.error('Failed to update relation connection:', error);
      throw error;
    }
  }

  async deleteRelationConnection(connectionId: string, relationId: string): Promise<void> {
    try {
      await this.delete(`/connections/${connectionId}/relations/${relationId}`);

      // Invalidate relations cache for both connections
      this.cacheManager.invalidateResource(`connection:${connectionId}`);
      this.cacheManager.invalidateResource(`connection:${relationId}`);
      this.cacheManager.invalidateResource('relations');
    } catch (error) {
      console.error(`Failed to delete relation connection between ${connectionId} and ${relationId}:`, error);
      throw error;
    }
  }

  async getRelationConnectionLogs(connectionId: string, forceRefresh = false): Promise<RelationConnectionLogs[]> {
    const cacheKey = `connection:${connectionId}:logs`;

    if (!forceRefresh && this.cacheManager.has(cacheKey)) {
      return this.cacheManager.get(cacheKey);
    }

    try {
      const data = await this.get<RelationConnectionLogs[]>(`/connections/${connectionId}/relations/logs`);

      // Cache with dependency on the connection
      // Use shorter TTL for logs (1 minute) by passing an explicit TTL value
      const logsExpiryTime = 60 * 1000; // 1 minute
      this.cacheManager.set(cacheKey, data, [`connection:${connectionId}`, 'logs'], logsExpiryTime);

      return data;
    } catch (error) {
      console.error(`Failed to get relation connection logs for ${connectionId}:`, error);
      throw error;
    }

  }
}

export const connectionService = new ConnectionService(new CacheManager());
