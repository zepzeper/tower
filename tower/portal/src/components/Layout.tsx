import { useEffect, useState } from 'react';
import { Outlet, useNavigate } from 'react-router-dom';
import { authService } from '../services/authService';
import TopBar from './TopBar';
import Sidebar from './Sidebar';
import React from 'react';
import { useTheme } from '../context/ThemeContext';

const Layout = () => {
  const [isVerifying, setIsVerifying] = useState(true);
  const navigate = useNavigate();
  const { theme } = useTheme();

  useEffect(() => {
    const verifyAuth = async () => {
      setIsVerifying(true);
      try {
        const user = await authService.getCurrentUser(true); // Force refresh
        if (!user) {
          navigate('/login');
        }
      } catch (error) {
        console.error('Auth verification failed:', error);
        navigate('/login');
      } finally {
        setIsVerifying(false);
      }
    };

    verifyAuth();
  }, [navigate]);

  if (isVerifying) {
    return <div>Verifying your session...</div>;
  }

  return (
    <div className={`flex h-screen ${theme === 'light' ? 'bg-gray-50' : 'bg-gray-950'}`}>
      <Sidebar />
      <div className="flex-1 flex flex-col overflow-hidden">
        <TopBar />
        <main className={`flex-1 overflow-y-auto p-6 ${theme === 'light' ? 'bg-gray-50' : 'bg-gray-950'}`}>
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default Layout;
