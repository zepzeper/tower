import { Routes, Route, Navigate } from 'react-router-dom';
import { ThemeProvider } from './context/ThemeContext';
import Layout from './components/Layout';
import Login from './pages/auth/Login';
import Register from './pages/auth/Register';
import ForgotPassword from './pages/auth/ForgotPassword';
import ResetPassword from './pages/auth/ResetPassword';
import LanguageWrapper from './components/LanguageWrapper';
import ProtectedRoute from './components/ProtectedRoute';
import Dashboard from './pages/Dashboard';
import Connections from './pages/Connections';
import Integrations from './pages/Integrations';
import ConnectionDetails from './pages/connections/ConnectionDetails';
import ConnectionEdit from './pages/connections/ConnectionEdit';
import Analytics from './pages/Analytics';
import Users from './pages/Users';
import React, { useState, useEffect } from 'react';
import { authService } from './services/authService';
import Account from './pages/Account';

const App: React.FC = () => {
  const [isAuthInitialized, setIsAuthInitialized] = useState(false);

  useEffect(() => {
    const initAuth = async () => {
      try {
        // This will try to load the user from token if one exists
        await authService.initialize();
      } catch (error) {
        console.error("Auth initialization failed:", error);
      } finally {
        setIsAuthInitialized(true);
      }
    };

    initAuth();
  }, []);

  // Show a loading state while auth is initializing
  if (!isAuthInitialized) {
    return (
      <ThemeProvider>
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-xl">Loading...</div>
        </div>
      </ThemeProvider>
    );
  }

  return (
    <ThemeProvider>
      <LanguageWrapper />
      <Routes>
        {['', '/:lang'].map((langPath) => (
          <Route key={langPath} path={langPath}>
            {/* Auth routes */}
            <Route path="login" element={<Login />} />
            <Route path="register" element={<Register />} />
            <Route path="forgot-password" element={<ForgotPassword />} />
            <Route path="reset-password" element={<ResetPassword />} />
            {/* Protected routes */}
            <Route
              path=""
              element={
                <ProtectedRoute>
                  <Layout />
                </ProtectedRoute>
              }
            >
              <Route index element={<Dashboard />} />
              <Route path="connections" element={<Connections />} />
              <Route path="connections/edit/:id" element={<ConnectionEdit />} />
              <Route path="connections/view/:id" element={<ConnectionDetails />} />
              <Route path="integrations" element={<Integrations />} />
              <Route path="analytics" element={<Analytics />} />
              <Route path="users" element={<Users />} />
              <Route path="account" element={<Account />} />
            </Route>
          </Route>
        ))}
        {/* Fallback */}
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </ThemeProvider>
  );
};

export default App;
