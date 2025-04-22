import { requestHandler } from '../handler/RequestHandler';

let cachedConnectionTypes: ConnectionType[] | null = null;

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
  initiator_id: string;
  target_id: string;
  connection_type: string;
  message: string
  created_at: string;
}

class ConnectionService {
  async getApiConnections(): Promise<ApiConnection[]> {
    return requestHandler.get<ApiConnection[]>('/connections/all');
  }

  async getApiConnectionWithConfig(id: string): Promise<{
    connection: ApiConnection;
    configs: ApiConnectionConfig[];
  }> {
    return requestHandler.get<{
      connection: ApiConnection;
      configs: ApiConnectionConfig[];
    }>(`/connections/${id}`);
  }

  async createApiConnection(data: ApiConnectionCreateRequest): Promise<ApiConnection> {
    return requestHandler.post<ApiConnection>('/connections', data);
  }

  async updateApiConnection(data: ApiConnectionUpdateRequest): Promise<ApiConnection> {
    return requestHandler.patch<ApiConnection>(`/connections/patch/${data.id}`, data);
  }

  async deleteApiConnection(id: string): Promise<void> {
    return requestHandler.delete(`/connections/delete/${id}`);
  }

  async testApiConnection(id: string): Promise<{
    success: boolean;
    message: string;
  }> {
    return requestHandler.get<{ success: boolean; message: string }>(`/connections/test/${id}`);
  }

  async getConnectionTypesFromFile(): Promise<ConnectionType[]> {
    if (cachedConnectionTypes) return cachedConnectionTypes;

    // This one still uses fetch directly since it's accessing a local file, not an API endpoint
    const res = await fetch('/schema/apis.json');
    if (!res.ok) throw new Error('Failed to load connection types');

    const data = (await res.json()) as ConnectionType[];
    cachedConnectionTypes = data;
    return data;
  }

  async initiateOAuthFlow(connectionType: string): Promise<{ url: string }> {
    return requestHandler.post<{ url: string }>(`/oauth/${connectionType}/initiate`, {});
  }


  async getRelationConnections(connectionId: string): Promise<RelationConnection[]> {
    return requestHandler.get<RelationConnection[]>(`/connections/${connectionId}/relations`);
  }

  async addRelationConnection(
    relationData: RelationConnection
  ): Promise<RelationConnection> {
    return requestHandler.post<RelationConnection>(
      `/connections/${relationData.initiator_id}/relations/create`,
      relationData
    );
  }

  async updateRelationConnection(
    relationData: RelationConnection
  ): Promise<RelationConnection> {
    return requestHandler.patch<RelationConnection>(
      `/connections/${relationData.initiator_id}/relations/${relationData.target_id}`,
      relationData
    );
  }

  async deleteRelationConnection(connectionId: string, relationId: string): Promise<void> {
    return requestHandler.delete(`/connections/${connectionId}/relations/${relationId}`);
  }


  async getRelationConnectionLogs(connectionId: string): Promise<RelationConnectionLogs[]> {
    return requestHandler.get(`/connections/${connectionId}/relations/logs`);
  }

}

export const connectionService = new ConnectionService();
