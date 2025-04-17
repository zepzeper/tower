// dashboard/src/components/Layout.jsx
import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Topbar from './Topbar';

const Layout = () => (
  <div className="flex h-screen">
    <Sidebar />
    <div className="flex flex-col flex-1 overflow-hidden">
      <Topbar />
      <main className="flex-1 overflow-auto bg-gray-50 p-4">
        <Outlet />
      </main>
    </div>
  </div>
);

export default Layout;
