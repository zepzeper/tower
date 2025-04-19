import { createContext, useState, useEffect, useContext } from 'react';

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    // Check if user is logged in on initial load
    const checkAuth = async () => {
      try {
        const token = localStorage.getItem('token');
        
        if (token) {
          // In a real application, you would validate the token with your API
          // const response = await api.get('/auth/me', {
          //   headers: { Authorization: `Bearer ${token}` }
          // });
          // setUser(response.data);
          
          // For demo purposes, set a demo user
          setUser({
            id: '1',
            name: 'Admin User',
            email: 'admin@example.com',
            role: 'admin'
          });
          setIsAuthenticated(true);
        }
      } catch (error) {
        // If token is invalid, clear it
        localStorage.removeItem('token');
        setUser(null);
        setIsAuthenticated(false);
      } finally {
        setIsLoading(false);
      }
    };

    checkAuth();
  }, []);

    const login = async (email, password) => {
        try {
            const response = await fetch('http://localhost:8080/internal/auth/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ email, password }),
            });

            if (!response.ok) {
                throw new Error('Login failed');
            }

            const data = await response.json();

            // Store the token
            localStorage.setItem('token', data.token);
            setUser(data.user);
            setIsAuthenticated(true);

            return data.user;
        } catch (error) {
            throw new Error(error.message || 'Login failed');
        }
    };

  const register = async (userData) => {
    // For demo purposes
    return new Promise((resolve, reject) => {
      setTimeout(() => {
        const token = 'demo-token-' + Math.random().toString(36).substring(2);
        const user = {
          id: '2',
          name: userData.name,
          email: userData.email,
          role: 'user'
        };
        
        localStorage.setItem('token', token);
        setUser(user);
        setIsAuthenticated(true);
        resolve(user);
      }, 1000);
    });
  };

  const logout = () => {
    localStorage.removeItem('token');
    setUser(null);
    setIsAuthenticated(false);
  };

  const value = {
    user,
    isAuthenticated,
    isLoading,
    login,
    register,
    logout
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
