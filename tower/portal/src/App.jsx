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

// Protected route component
const ProtectedRoute = ({ children }) => {
  const { isAuthenticated, isLoading } = useAuth();
  
  if (isLoading) {
    return <div className="flex h-screen items-center justify-center">Loading...</div>;
  }
  
  if (!isAuthenticated) {
    return <Navigate to="/login" />;
  }
  
  return children;
};

const App = () => {
  return (
    <ThemeProvider>
      {/* Add LanguageWrapper for URL-based language handling */}
      <LanguageWrapper />
      
      <Routes>
        {/* Routes with language prefixes */}
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
              {/* Dashboard routes */}
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
        
        {/* Fallback route */}
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </ThemeProvider>
  );
};

export default App;
