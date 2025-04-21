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

class ConnectionService {
  private baseUrl: string;

  constructor() {
    this.baseUrl = '/api';
  }

  async getApiConnections(): Promise<ApiConnection[]> {
    const response = await fetch(`${this.baseUrl}/connections`);
    if (!response.ok) throw new Error(`Failed to fetch API connections: ${response.statusText}`);
    return await response.json();
  }

  async getApiConnectionWithConfig(id: string): Promise<{
    connection: ApiConnection;
    configs: ApiConnectionConfig[];
  }> {
    const response = await fetch(`${this.baseUrl}/connections/${id}`);
    if (!response.ok) throw new Error(`Failed to fetch API connection ${id}: ${response.statusText}`);
    return await response.json();
  }

  async createApiConnection(data: ApiConnectionCreateRequest): Promise<ApiConnection> {
    const response = await fetch(`${this.baseUrl}/connections`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error(`Failed to create API connection: ${response.statusText}`);
    return await response.json();
  }

  async updateApiConnection(data: ApiConnectionUpdateRequest): Promise<ApiConnection> {
    const response = await fetch(`${this.baseUrl}/connections/${data.id}`, {
      method: 'PATCH',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error(`Failed to update API connection ${data.id}: ${response.statusText}`);
    return await response.json();
  }

  async deleteApiConnection(id: string): Promise<void> {
    const response = await fetch(`${this.baseUrl}/connections/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) throw new Error(`Failed to delete API connection ${id}: ${response.statusText}`);
  }

  async testApiConnection(data: ApiConnectionCreateRequest): Promise<{
    success: boolean;
    message: string;
  }> {
    const response = await fetch(`${this.baseUrl}/connections/test`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    const result = await response.json();
    if (!response.ok) return { success: false, message: result.message || `Request failed with status: ${response.status}` };
    return result;
  }

  async getConnectionTypesFromFile(): Promise<ConnectionType[]> {
    if (cachedConnectionTypes) return cachedConnectionTypes;

    const res = await fetch('/schema/apis.json');
    if (!res.ok) throw new Error('Failed to load connection types');

    const data = (await res.json()) as ConnectionType[];
    cachedConnectionTypes = data;
    return data;
  }

  async initiateOAuthFlow(connectionType: string): Promise<{ url: string }> {
    const response = await fetch(`${this.baseUrl}/oauth/${connectionType}/initiate`, {
      method: 'POST',
    });
    if (!response.ok) throw new Error(`Failed to initiate OAuth flow: ${response.statusText}`);
    return await response.json();
  }

}

export const connectionService = new ConnectionService();
