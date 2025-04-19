import { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../context/ThemeContext';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { Activity, Users, Database, ArrowUp, ArrowDown } from 'lucide-react';

const Dashboard = () => {
  const { t } = useTranslation('pages');
  const { theme } = useTheme();
  const [stats, setStats] = useState({
    activeConnections: 24,
    totalUsers: 189,
    apiCalls: 15728,
    growth: 18.2
  });

  // Example data for charts
  const connectionData = [
    { name: 'Jan', connections: 12 },
    { name: 'Feb', connections: 19 },
    { name: 'Mar', connections: 15 },
    { name: 'Apr', connections: 21 },
    { name: 'May', connections: 18 },
    { name: 'Jun', connections: 24 },
    { name: 'Jul', connections: 30 },
  ];

  return (
    <div className="p-6">
      <div className="mb-8">
        <h1 className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
          {t('navigation.dashboard')}
        </h1>
        <p className={`mt-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-600'}`}>
          {t('dashboard.overview')}
        </p>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatCard 
          title={t('dashboard.activeConnections')}
          value={stats.activeConnections}
          icon={<Activity size={24} />}
          theme={theme}
        />
        <StatCard 
          title={t('dashboard.totalUsers')}
          value={stats.totalUsers}
          icon={<Users size={24} />}
          theme={theme}
        />
        <StatCard 
          title={t('dashboard.apiCalls')}
          value={stats.apiCalls.toLocaleString()}
          icon={<Database size={24} />}
          theme={theme}
        />
        <StatCard 
          title={t('dashboard.growth')}
          value={`${stats.growth}%`}
          icon={<ArrowUp size={24} className="text-green-500" />}
          trend="up"
          theme={theme}
        />
      </div>

      {/* Charts */}
      <div className={`p-6 rounded-lg shadow-md mb-8 ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
        <h2 className={`text-lg font-medium mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
          {t('dashboard.connectionActivity')}
        </h2>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={connectionData}
              margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
            >
              <CartesianGrid strokeDasharray="3 3" className={theme === 'dark' ? 'stroke-gray-700' : 'stroke-gray-200'} />
              <XAxis 
                dataKey="name" 
                className={theme === 'dark' ? 'fill-gray-400' : 'fill-gray-600'} 
              />
              <YAxis 
                className={theme === 'dark' ? 'fill-gray-400' : 'fill-gray-600'} 
              />
              <Tooltip 
                contentStyle={{ 
                  backgroundColor: theme === 'dark' ? '#1f2937' : '#ffffff',
                  color: theme === 'dark' ? '#f3f4f6' : '#1f2937',
                  border: theme === 'dark' ? '1px solid #374151' : '1px solid #e5e7eb'
                }} 
              />
              <Legend />
              <Bar dataKey="connections" fill="#10b981" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Recent Activity */}
      <div className={`p-6 rounded-lg shadow-md ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
        <h2 className={`text-lg font-medium mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
          {t('dashboard.recentActivity')}
        </h2>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className={theme === 'dark' ? 'bg-gray-700' : 'bg-gray-50'}>
              <tr>
                <th scope="col" className={`px-6 py-3 text-left text-xs font-medium ${theme === 'dark' ? 'text-gray-300 uppercase tracking-wider' : 'text-gray-500 uppercase tracking-wider'}`}>
                  {t('dashboard.eventType')}
                </th>
                <th scope="col" className={`px-6 py-3 text-left text-xs font-medium ${theme === 'dark' ? 'text-gray-300 uppercase tracking-wider' : 'text-gray-500 uppercase tracking-wider'}`}>
                  {t('dashboard.connection')}
                </th>
                <th scope="col" className={`px-6 py-3 text-left text-xs font-medium ${theme === 'dark' ? 'text-gray-300 uppercase tracking-wider' : 'text-gray-500 uppercase tracking-wider'}`}>
                  {t('dashboard.timestamp')}
                </th>
                <th scope="col" className={`px-6 py-3 text-left text-xs font-medium ${theme === 'dark' ? 'text-gray-300 uppercase tracking-wider' : 'text-gray-500 uppercase tracking-wider'}`}>
                  {t('dashboard.status')}
                </th>
              </tr>
            </thead>
            <tbody className={`divide-y ${theme === 'dark' ? 'divide-gray-700' : 'divide-gray-200'}`}>
              {[1, 2, 3, 4, 5].map((item) => (
                <tr key={item} className={theme === 'dark' ? 'bg-gray-800' : 'bg-white'}>
                  <td className={`px-6 py-4 whitespace-nowrap text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-900'}`}>
                    {t('dashboard.apiRequest')}
                  </td>
                  <td className={`px-6 py-4 whitespace-nowrap text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-900'}`}>
                    API Integration {item}
                  </td>
                  <td className={`px-6 py-4 whitespace-nowrap text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-500'}`}>
                    {new Date().toLocaleString()}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${
                      item % 3 === 0
                        ? 'bg-red-100 text-red-800 dark:bg-red-900 dark:bg-opacity-30 dark:text-red-300'
                        : 'bg-green-100 text-green-800 dark:bg-green-900 dark:bg-opacity-30 dark:text-green-300'
                    }`}>
                      {item % 3 === 0 ? t('dashboard.failed') : t('dashboard.success')}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

// Stat Card Component
const StatCard = ({ title, value, icon, trend, theme }) => {
  return (
    <div className={`p-6 rounded-lg shadow-md ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
      <div className="flex justify-between items-start">
        <div>
          <p className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-500'}`}>{title}</p>
          <p className={`mt-2 text-3xl font-semibold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>{value}</p>
          {trend && (
            <div className="mt-2 flex items-center">
              {trend === 'up' ? (
                <ArrowUp size={16} className="text-green-500 mr-1" />
              ) : (
                <ArrowDown size={16} className="text-red-500 mr-1" />
              )}
              <span className={trend === 'up' ? 'text-green-500 text-sm' : 'text-red-500 text-sm'}>
                {trend === 'up' ? '+8.2%' : '-4.5%'}
              </span>
              <span className={`text-xs ml-1 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>vs last month</span>
            </div>
          )}
        </div>
        <div className={`p-3 rounded-full ${theme === 'dark' ? 'bg-gray-700' : 'bg-green-50'}`}>
          {icon}
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
