import { requestHandler } from '../handler/RequestHandler';

export interface User {
  id: string;
  name: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

class AuthService {
  private currentUser: User | null = null;

  // Store token in localStorage
  setToken(token: string): void {
    localStorage.setItem('auth_token', token);
  }

  // Get token from localStorage
  getToken(): string | null {
    return localStorage.getItem('auth_token');
  }

  // Remove token from localStorage
  removeToken(): void {
    localStorage.removeItem('auth_token');
  }

  // Check if user is authenticated
  isAuthenticated(): boolean {
    return !!this.getToken();
  }

  // Get the current user object
  getUser(): User | null {
    return this.currentUser;
  }

  setUser(user: User | null): void {
    this.currentUser = user;
  }

  // Register a new user
  async register(data: RegisterRequest): Promise<User> {
    try {
      const result = await requestHandler.post<AuthResponse>('/auth/register', data);
      this.setToken(result.token);
      this.setUser(result.user);
      return result.user;
    } catch (error) {
      throw error;
    }
  }

  async login(email: string, password: string): Promise<User> {
    try {
      console.log('Attempting login for:', email);

      const result = await requestHandler.post<AuthResponse>('/auth/login', { email, password });

      console.log(result.user)

      if (!result.token) {
        throw new Error('No token received from server');
      }

      console.log('Login successful, received token');

      // Set token first, so it's available for subsequent requests
      this.setToken(result.token);

      // If user data came with the response, store it
      if (result.user) {
        this.currentUser = result.user;
        return result.user;
      }

      // If we didn't get user data with the login response,
      // make a separate call to get it
      return await this.getCurrentUser(true); // force refresh
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  }

  // Logout user
  async logout(): Promise<void> {
    try {
      await requestHandler.post('/auth/logout', {});
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      this.removeToken();
      this.currentUser = null;
    }
  }

  // Refresh auth token
  async refreshToken(): Promise<string | null> {
    try {
      const result = await requestHandler.post<{ token: string }>('/auth/refresh', {});
      this.setToken(result.token);
      return result.token;
    } catch (error) {
      console.error('Token refresh error:', error);
      return null;
    }
  }

  // Forgot password
  async forgotPassword(email: string): Promise<void> {
    await requestHandler.post('/auth/forgot-password', { email });
  }

  // Reset password
  async resetPassword(token: string, password: string): Promise<void> {
    await requestHandler.post('/auth/reset-password', { token, password });
  }

  // Get current user data
  async getCurrentUser(forceRefresh = false): Promise<User | null> {
    // Return cached user if available and not forcing refresh
    if (this.currentUser && !forceRefresh) {
      return this.currentUser;
    }

    if (!this.isAuthenticated()) {
      return null;
    }

    try {
      console.log('Fetching current user from API');
      const user = await requestHandler.get<User>('/user/me');

      if (user) {
        console.log('User fetched successfully');
        this.currentUser = user;
        return user;
      } else {
        console.log('No user data returned');
        return null;
      }
    } catch (error: any) {
      console.error('Error fetching user:', error);
      return null;
    }
  }

  // Initialize auth state
  async initialize(): Promise<User | null> {
    if (this.isAuthenticated()) {
      return await this.getCurrentUser();
    }
    return null;
  }
}

// Create singleton instance
export const authService = new AuthService();
