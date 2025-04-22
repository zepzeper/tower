export interface ApiResponse<T> {
  data: T;
}

class RequestHandler {
  private baseUrl: string;

  constructor(baseUrl: string = '/api') {
    this.baseUrl = baseUrl;
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

  async get<T>(endpoint: string): Promise<T> {
    return this.handleUnauthorized(async () => {
      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        headers: this.getAuthHeaders(),
      });
      if (!response.ok) {
        this.handleErrorResponse(response);
      }
      const json: ApiResponse<T> = await response.json();
      return json.data;
    });
  }

  async post<T>(endpoint: string, data: any): Promise<T> {
    return this.handleUnauthorized(async () => {
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
    });
  }

  async patch<T>(endpoint: string, data: any): Promise<T> {
    return this.handleUnauthorized(async () => {
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
    });
  }

  async delete(endpoint: string): Promise<void> {
    return this.handleUnauthorized(async () => {
      const response = await fetch(`${this.baseUrl}${endpoint}`, {
        method: 'DELETE',
        headers: this.getAuthHeaders(),
      });
      if (!response.ok) {
        this.handleErrorResponse(response);
      }
    });
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

export const requestHandler = new RequestHandler();
