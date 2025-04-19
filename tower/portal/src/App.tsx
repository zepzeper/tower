import { Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from './context/AuthContext';
import { ThemeProvider } from './context/ThemeContext';
import Layout from './components/Layout';
import Login from './pages/auth/Login';
import Register from './pages/auth/Register';
import ForgotPassword from './pages/auth/ForgotPassword';
import ResetPassword from './pages/auth/ResetPassword';
import LanguageWrapper from './components/LanguageWrapper';

import Dashboard from './pages/Dashboard';
import Connections from './pages/Connections';
import Integrations from './pages/Integrations';
import ConnectionDetails from './pages/connections/ConnectionDetails';
import ConnectionEdit from './pages/connections/ConnectionEdit';
import Analytics from './pages/Analytics';
import Users from './pages/Users';

import React from 'react';

interface ProtectedRouteProps {
  children: React.ReactElement;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return <div className="flex h-screen items-center justify-center">Loading...</div>;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return children;
};

const App: React.FC = () => {
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
                //<ProtectedRoute>
                  <Layout />
                //</ProtectedRoute>
              }
            >
              <Route index element={<Dashboard />} />
              <Route path="connections" element={<Connections />} />
              <Route path="connections/edit/:id" element={<ConnectionEdit />} />
              <Route path="connections/view/:id" element={<ConnectionDetails />} />
              <Route path="integrations" element={<Integrations />} />
              <Route path="analytics" element={<Analytics />} />
              <Route path="users" element={<Users />} />
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
