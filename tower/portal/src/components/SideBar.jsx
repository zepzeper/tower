// dashboard/src/components/Sidebar.jsx
import { Link, useLocation } from 'react-router-dom';
import { LayoutDashboard, Link as LinkIcon, Settings, Database, BarChart, Users } from 'lucide-react';

const Sidebar = () => {
  const location = useLocation();
  
  const navigation = [
    { name: 'Dashboard', href: '/', icon: LayoutDashboard },
    { name: 'Connections', href: '/connections', icon: LinkIcon },
    { name: 'API Integrations', href: '/integrations', icon: Database },
    { name: 'Analytics', href: '/analytics', icon: BarChart },
    { name: 'Users', href: '/users', icon: Users },
    { name: 'Settings', href: '/settings', icon: Settings },
  ];
  
  return (
    <div className="bg-slate-950 text-white w-64 min-h-full hidden lg:flex flex-col">
      <div className="p-4 border-b border-green-600">
        <div className="flex items-center space-x-2">
          <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
          </svg>
          <span className="text-2xl font-bold">Tower<span className="text-green-300">API</span></span>
        </div>
      </div>
      
      <nav className="flex-1 p-4">
        <ul className="space-y-1">
          {navigation.map((item) => (
            <li key={item.name}>
              <Link
                to={item.href}
                className={`flex items-center space-x-3 px-3 py-2 rounded-lg ${
                  location.pathname === item.href
                    ? 'bg-green-800 text-white'
                    : 'text-green-100 hover:bg-green-600 hover:text-white'
                }`}
              >
                <item.icon className="h-5 w-5" />
                <span>{item.name}</span>
              </Link>
            </li>
          ))}
        </ul>
      </nav>
      
      <div className="p-4 border-t border-green-600">
        <div className="flex items-center space-x-3 p-2">
          <div className="h-8 w-8 rounded-full bg-green-500 flex items-center justify-center">
            <span className="font-bold">A</span>
          </div>
          <div>
            <p className="text-sm font-medium">Admin User</p>
            <p className="text-xs text-green-300">admin@example.com</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Sidebar;
