import { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { 
  LayoutDashboard, 
  Link as LinkIcon, 
  Settings, 
  Database, 
  BarChart, 
  Users,
  Zap,
  ChevronDown,
  Shield,
  UserCog,
  Bell,
  Globe
} from 'lucide-react';
import { useTheme } from '../context/ThemeContext';

const Sidebar = () => {
  const location = useLocation();
  const { theme } = useTheme();
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);
  const { t, i18n } = useTranslation('pages');
  
  const navigation = [
    { name: t('navigation.dashboard'), href: '/', icon: LayoutDashboard },
    { name: t('navigation.connections'), href: '/connections', icon: LinkIcon },
    { name: t('navigation.apiIntegrations'), href: '/integrations', icon: Database },
    { name: t('navigation.analytics'), href: '/analytics', icon: BarChart },
    { name: t('navigation.users'), href: '/users', icon: Users },
  ];
  
  const settingsNavigation = [
    { name: t('navigation.general'), href: '/settings/general', icon: Settings },
    { name: t('navigation.security'), href: '/settings/security', icon: Shield },
    { name: t('navigation.account'), href: '/settings/account', icon: UserCog },
    { name: t('navigation.notifications'), href: '/settings/notifications', icon: Bell },
    { name: t('navigation.regional'), href: '/settings/regional', icon: Globe },
  ];
  
  return (
    <div className={`${theme === 'light' ? 'bg-white text-gray-800 border-r border-gray-200' : 'bg-gray-900 text-white border-r border-gray-700'} w-64 min-h-full hidden lg:flex flex-col transition-colors duration-200`}>
      <div className={`p-4 border-b ${theme === 'light' ? 'border-gray-200' : 'border-gray-700'}`}>
        <div className="flex items-center space-x-2">
          <Zap className={`w-8 h-8 ${theme === 'light' ? 'text-green-600' : 'text-green-400'}`} />
          <span className="text-2xl font-bold">{t('sidebar.logo')}<span className={`${theme === 'light' ? 'text-green-600' : 'text-green-400'}`}>API</span></span>
        </div>
      </div>
      
      <nav className="flex-1 p-4">
        <ul className="space-y-1">
          {navigation.map((item) => (
            <li key={item.name}>
              <Link
                to={item.href}
                className={`flex items-center space-x-3 px-3 py-2 rounded-lg transition-colors duration-200 ${
                  location.pathname === item.href
                    ? theme === 'light' 
                      ? 'bg-green-50 text-green-700 border border-green-100' 
                      : 'bg-gray-800 text-green-400 border border-gray-700'
                    : theme === 'light'
                      ? 'text-gray-700 hover:bg-gray-100 hover:text-green-600'
                      : 'text-gray-300 hover:bg-gray-800 hover:text-white'
                }`}
              >
                <item.icon className="h-5 w-5" />
                <span>{item.name}</span>
              </Link>
            </li>
          ))}
          
          {/* Settings collapsible section */}
          <li>
            <button
              onClick={() => setIsSettingsOpen(!isSettingsOpen)}
              className={`flex items-center justify-between w-full px-3 py-2 rounded-lg transition-colors duration-200 ${
                location.pathname.startsWith('/settings')
                  ? theme === 'light' 
                    ? 'bg-green-50 text-green-700 border border-green-100' 
                    : 'bg-gray-800 text-green-400 border border-gray-700'
                  : theme === 'light'
                    ? 'text-gray-700 hover:bg-gray-100 hover:text-green-600'
                    : 'text-gray-300 hover:bg-gray-800 hover:text-white'
              }`}
            >
              <div className="flex items-center space-x-3">
                <Settings className="h-5 w-5" />
                <span>{t('navigation.settings')}</span>
              </div>
              <ChevronDown className={`h-4 w-4 transition-transform ${isSettingsOpen ? 'rotate-180' : ''}`} />
            </button>
            
            {/* Settings submenu */}
            {isSettingsOpen && (
              <ul className={`mt-1 ml-4 pl-2 space-y-1 ${theme === 'light' ? 'border-l border-gray-200' : 'border-l border-gray-700'}`}>
                {settingsNavigation.map((item) => (
                  <li key={item.name}>
                    <Link
                      to={item.href}
                      className={`flex items-center space-x-3 px-3 py-2 rounded-lg transition-colors duration-200 ${
                        location.pathname === item.href
                          ? theme === 'light' 
                            ? 'bg-green-50 text-green-700 border border-green-100' 
                            : 'bg-gray-800 text-green-400 border border-gray-700'
                          : theme === 'light'
                            ? 'text-gray-700 hover:bg-gray-100 hover:text-green-600'
                            : 'text-gray-300 hover:bg-gray-800 hover:text-white'
                      }`}
                    >
                      <item.icon className="h-4 w-4" />
                      <span className="text-sm">{item.name}</span>
                    </Link>
                  </li>
                ))}
              </ul>
            )}
          </li>
        </ul>
      </nav>
      
    </div>
  );
};

export default Sidebar;

//<div className={`p-4 border-t ${theme === 'light' ? 'border-gray-200' : 'border-gray-700'}`}>
//  <div className="flex items-center space-x-3 p-2">
//    <div className={`h-8 w-8 rounded-full ${theme === 'light' ? 'bg-green-100 text-green-600' : 'bg-green-800 text-green-100'} flex items-center justify-center`}>
//      <span className="font-bold">A</span>
//    </div>
//    <div>
//      <p className="text-sm font-medium">Admin User</p>
//      <p className={`text-xs ${theme === 'light' ? 'text-gray-500' : 'text-gray-400'}`}>admin@example.com</p>
//    </div>
//  </div>
//</div>
