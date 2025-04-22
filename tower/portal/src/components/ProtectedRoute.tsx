import React, { useEffect, useState } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { authService } from '../services/authService';

interface ProtectedRouteProps {
  children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const [authState, setAuthState] = useState<{
    isAuthenticated: boolean;
    isLoading: boolean;
  }>({
    isAuthenticated: false,
    isLoading: true
  });

  const location = useLocation();

  useEffect(() => {
    const checkAuth = async () => {
      try {
        // First check if we have a token
        if (!authService.isAuthenticated()) {
          setAuthState({
            isAuthenticated: false,
            isLoading: false
          });
          return;
        }

        // If we have a token, try to get the current user
        // This will either return the cached user or fetch from the server
        const user = await authService.getCurrentUser();

        setAuthState({
          isAuthenticated: !!user,
          isLoading: false
        });
      } catch (error) {
        console.error("Authentication check failed:", error);
        setAuthState({
          isAuthenticated: false,
          isLoading: false
        });
      }
    };

    checkAuth();
  }, []);

  if (authState.isLoading) {
    return <div>Loading...</div>; // Or a proper loading spinner component
  }

  if (!authState.isAuthenticated) {
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  return <>{children}</>;
};

export default ProtectedRoute;
