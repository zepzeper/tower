import CacheManager from '../services/cacheService';
const apiUrl = import.meta.env.VITE_API_URL || '/api';

export interface ApiResponse<T> {
  data: T;
}

class RequestHandler {
  private baseUrl: string;
  private pendingRequests: Map<string, Promise<any>> = new Map();
  protected cacheManager: CacheManager;

  constructor(cacheManager: CacheManager) {
    this.baseUrl = apiUrl;
    this.cacheManager = cacheManager;
  }

  private getAuthHeaders(): HeadersInit {
    const token = localStorage.getItem('auth_token');
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    return headers;
  }

  private createRequestKey(method: string, endpoint: string, data?: any): string {
    if (data) {
      // For methods with request bodies, include body in the key
      return `${method}:${endpoint}:${JSON.stringify(data)}`;
    }
    return `${method}:${endpoint}`;
  }

  async get<T>(endpoint: string): Promise<T> {
    const requestKey = this.createRequestKey('GET', endpoint);

    // If we already have a pending request for this endpoint, return that Promise
    if (this.pendingRequests.has(requestKey)) {
      return this.pendingRequests.get(requestKey);
    }

    const request = this.handleUnauthorized(async () => {
      try {
        const response = await fetch(`${this.baseUrl}${endpoint}`, {
          headers: this.getAuthHeaders(),
        });
        if (!response.ok) {
          this.handleErrorResponse(response);
        }
        const json: ApiResponse<T> = await response.json();
        return json.data;
      } finally {
        // Remove from pending requests when done
        this.pendingRequests.delete(requestKey);
      }
    });

    this.pendingRequests.set(requestKey, request);

    return request;
  }

  async post<T>(endpoint: string, data: any): Promise<T> {
    const requestKey = this.createRequestKey('POST', endpoint, data);

    if (this.pendingRequests.has(requestKey)) {
      return this.pendingRequests.get(requestKey);
    }

    const request = this.handleUnauthorized(async () => {
      try {
        const response = await fetch(`${this.baseUrl}${endpoint}`, {
          method: 'POST',
          headers: this.getAuthHeaders(),
          body: JSON.stringify(data),
        });
        if (!response.ok) {
          this.handleErrorResponse(response);
        }
        const json: ApiResponse<T> = await response.json();
        return json.data;
      } finally {
        this.pendingRequests.delete(requestKey);
      }
    });

    this.pendingRequests.set(requestKey, request);

    return request;
  }

  async patch<T>(endpoint: string, data: any): Promise<T> {
    const requestKey = this.createRequestKey('PATCH', endpoint, data);

    if (this.pendingRequests.has(requestKey)) {
      return this.pendingRequests.get(requestKey);
    }

    const request = this.handleUnauthorized(async () => {
      try {
        const response = await fetch(`${this.baseUrl}${endpoint}`, {
          method: 'PATCH',
          headers: this.getAuthHeaders(),
          body: JSON.stringify(data),
        });
        if (!response.ok) {
          this.handleErrorResponse(response);
        }
        const json: ApiResponse<T> = await response.json();
        return json.data;
      } finally {
        this.pendingRequests.delete(requestKey);
      }
    });

    this.pendingRequests.set(requestKey, request);

    return request;
  }

  async delete(endpoint: string): Promise<void> {
    const requestKey = this.createRequestKey('DELETE', endpoint);

    if (this.pendingRequests.has(requestKey)) {
      return this.pendingRequests.get(requestKey);
    }

    const request = this.handleUnauthorized(async () => {
      try {
        const response = await fetch(`${this.baseUrl}${endpoint}`, {
          method: 'DELETE',
          headers: this.getAuthHeaders(),
        });
        if (!response.ok) {
          this.handleErrorResponse(response);
        }
      } finally {
        this.pendingRequests.delete(requestKey);
      }
    });

    this.pendingRequests.set(requestKey, request);

    return request;
  }

  private handleErrorResponse(response: Response): never {
    throw new Error(`Request failed with status: ${response.status} ${response.statusText}`);
  }

  private async handleUnauthorized<T>(reqFn: () => Promise<T>): Promise<T> {
    try {
      return await reqFn();
    } catch (error) {
      if (error instanceof Error && error.message.includes('401')) {
        const authService = await import('../services/authService').then((m) => m.authService);
        const newToken = await authService.refreshToken();
        if (newToken) {
          return await reqFn();
        } else {
          window.location.href = '/login';
        }
      }
      throw error;
    }
  }

}

export default RequestHandler;
