import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../context/ThemeContext';
import { 
  LineChart, Line, BarChart, Bar, PieChart, Pie, Cell, 
  XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer
} from 'recharts';
import { Calendar, Filter, Download, ArrowDown, ArrowUp } from 'lucide-react';

const Analytics = () => {
  const { t } = useTranslation('pages');
  const { theme } = useTheme();
  const [timeRange, setTimeRange] = useState('month'); // 'week', 'month', 'year'
  
  // Example data for charts
  const apiCallsData = [
    { name: '1 May', calls: 1200, errors: 80 },
    { name: '8 May', calls: 1800, errors: 120 },
    { name: '15 May', calls: 1600, errors: 90 },
    { name: '22 May', calls: 2400, errors: 110 },
    { name: '29 May', calls: 2100, errors: 130 }
  ];
  
  const performanceData = [
    { name: '1 May', responseTime: 120, throughput: 850 },
    { name: '8 May', responseTime: 140, throughput: 920 },
    { name: '15 May', responseTime: 110, throughput: 1050 },
    { name: '22 May', responseTime: 105, throughput: 1200 },
    { name: '29 May', responseTime: 115, throughput: 1100 }
  ];
  
  const integrationDistribution = [
    { name: 'Shopify', value: 35, color: '#10B981' },
    { name: 'HubSpot', value: 25, color: '#3B82F6' },
    { name: 'Stripe', value: 20, color: '#6366F1' },
    { name: 'Mailchimp', value: 15, color: '#F59E0B' },
    { name: 'Others', value: 5, color: '#71717A' }
  ];
  
  const COLORS = ['#10B981', '#3B82F6', '#6366F1', '#F59E0B', '#71717A'];
  
  const stats = [
    { 
      id: 1, 
      title: t('analytics.totalApiCalls'), 
      value: '52,543', 
      change: '+12.3%', 
      trend: 'up',
      changeText: t('analytics.vsLastPeriod')
    },
    { 
      id: 2, 
      title: t('analytics.avgResponseTime'), 
      value: '118ms', 
      change: '-8.1%', 
      trend: 'down',
      changeText: t('analytics.vsLastPeriod')
    },
    { 
      id: 3, 
      title: t('analytics.errorRate'), 
      value: '4.2%', 
      change: '+0.8%', 
      trend: 'up',
      changeText: t('analytics.vsLastPeriod'),
      bad: true
    },
    { 
      id: 4, 
      title: t('analytics.activeIntegrations'), 
      value: '8', 
      change: '+2', 
      trend: 'up',
      changeText: t('analytics.vsLastPeriod')
    }
  ];

  return (
    <div className="p-6">
      <div className="mb-6 flex flex-col sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
            {t('navigation.analytics')}
          </h1>
          <p className={`mt-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-600'}`}>
            {t('analytics.description')}
          </p>
        </div>
        <div className="mt-4 sm:mt-0 flex space-x-3">
          <button className={`flex items-center px-3 py-1.5 rounded border ${
            theme === 'dark' ? 'border-gray-600 bg-gray-700 text-gray-200' : 'border-gray-300 bg-white text-gray-700'
          }`}>
            <Calendar size={16} className="mr-2" />
            {timeRange === 'week' && t('analytics.lastWeek')}
            {timeRange === 'month' && t('analytics.lastMonth')}
            {timeRange === 'year' && t('analytics.lastYear')}
          </button>
          <button className={`flex items-center px-3 py-1.5 rounded border ${
            theme === 'dark' ? 'border-gray-600 bg-gray-700 text-gray-200' : 'border-gray-300 bg-white text-gray-700'
          }`}>
            <Filter size={16} className="mr-2" />
            {t('analytics.filters')}
          </button>
          <button className={`flex items-center px-3 py-1.5 rounded border ${
            theme === 'dark' ? 'border-gray-600 bg-gray-700 text-gray-200' : 'border-gray-300 bg-white text-gray-700'
          }`}>
            <Download size={16} className="mr-2" />
            {t('analytics.export')}
          </button>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {stats.map((stat) => (
          <div 
            key={stat.id} 
            className={`p-6 rounded-lg shadow-md ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}
          >
            <p className={`text-sm font-medium ${theme === 'dark' ? 'text-gray-300' : 'text-gray-500'}`}>
              {stat.title}
            </p>
            <p className={`mt-2 text-3xl font-semibold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
              {stat.value}
            </p>
            <div className="mt-2 flex items-center">
              {stat.trend === 'up' ? (
                <ArrowUp size={16} className={stat.bad ? 'text-red-500 mr-1' : 'text-green-500 mr-1'} />
              ) : (
                <ArrowDown size={16} className={stat.bad ? 'text-green-500 mr-1' : 'text-red-500 mr-1'} />
              )}
              <span className={
                stat.trend === 'up' 
                  ? stat.bad ? 'text-red-500 text-sm' : 'text-green-500 text-sm'
                  : stat.bad ? 'text-green-500 text-sm' : 'text-red-500 text-sm'
              }>
                {stat.change}
              </span>
              <span className={`text-xs ml-1 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                {stat.changeText}
              </span>
            </div>
          </div>
        ))}
      </div>

      {/* API Calls Chart */}
      <div className={`p-6 rounded-lg shadow-md mb-8 ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
        <h2 className={`text-lg font-medium mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
          {t('analytics.apiCallsOverTime')}
        </h2>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart
              data={apiCallsData}
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
              <Line type="monotone" dataKey="calls" stroke="#10b981" strokeWidth={2} activeDot={{ r: 8 }} name={t('analytics.calls')} />
              <Line type="monotone" dataKey="errors" stroke="#ef4444" strokeWidth={2} name={t('analytics.errors')} />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Performance & Distribution Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
        {/* Performance Chart */}
        <div className={`p-6 rounded-lg shadow-md ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
          <h2 className={`text-lg font-medium mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
            {t('analytics.performanceMetrics')}
          </h2>
          <div className="h-72">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart
                data={performanceData}
                margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
              >
                <CartesianGrid strokeDasharray="3 3" className={theme === 'dark' ? 'stroke-gray-700' : 'stroke-gray-200'} />
                <XAxis 
                  dataKey="name" 
                  className={theme === 'dark' ? 'fill-gray-400' : 'fill-gray-600'} 
                />
                <YAxis 
                  yAxisId="left"
                  orientation="left"
                  stroke="#10b981"
                  className={theme === 'dark' ? 'fill-gray-400' : 'fill-gray-600'}
                />
                <YAxis 
                  yAxisId="right"
                  orientation="right"
                  stroke="#3b82f6"
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
                <Bar yAxisId="left" dataKey="responseTime" fill="#10b981" name={t('analytics.responseTime')} />
                <Bar yAxisId="right" dataKey="throughput" fill="#3b82f6" name={t('analytics.throughput')} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Integration Distribution Chart */}
        <div className={`p-6 rounded-lg shadow-md ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
          <h2 className={`text-lg font-medium mb-4 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
            {t('analytics.integrationDistribution')}
          </h2>
          <div className="h-72">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={integrationDistribution}
                  cx="50%"
                  cy="50%"
                  labelLine={true}
                  label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                  outerRadius={80}
                  fill="#8884d8"
                  dataKey="value"
                  nameKey="name"
                >
                  {integrationDistribution.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: theme === 'dark' ? '#1f2937' : '#ffffff',
                    color: theme === 'dark' ? '#f3f4f6' : '#1f2937',
                    border: theme === 'dark' ? '1px solid #374151' : '1px solid #e5e7eb'
                  }} 
                />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      {/* Tips Card */}
      <div className={`p-6 rounded-lg shadow-md ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
        <h2 className={`text-lg font-medium mb-2 ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
          {t('analytics.performanceTips')}
        </h2>
        <ul className={`list-disc pl-5 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-600'}`}>
          <li className="mb-1">{t('analytics.tip1')}</li>
          <li className="mb-1">{t('analytics.tip2')}</li>
          <li className="mb-1">{t('analytics.tip3')}</li>
          <li>{t('analytics.tip4')}</li>
        </ul>
      </div>
    </div>
  );
};

export default Analytics;
