import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import TopBar from './TopBar';
import { useTheme } from '../context/ThemeContext';

const Layout = () => {
  const { theme } = useTheme();

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
