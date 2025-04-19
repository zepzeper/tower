import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../context/ThemeContext';
import { Search, ShoppingCart, Users, CreditCard, Mail, Database } from 'lucide-react';

const Integrations = () => {
  const { t } = useTranslation('pages');
  const { theme } = useTheme();
  const [searchTerm, setSearchTerm] = useState('');
  const [activeCategory, setActiveCategory] = useState('all');

  // Integration categories
  const categories = [
    { id: 'all', name: t('integrations.categories.allCategories') },
    { id: 'ecommerce', name: t('integrations.categories.ecommerce'), icon: ShoppingCart },
    { id: 'crm', name: t('integrations.categories.crm'), icon: Users },
    { id: 'finance', name: t('integrations.categories.finance'), icon: CreditCard },
    { id: 'marketing', name: t('integrations.categories.marketing'), icon: Mail },
  ];

  // Example integrations data
  const integrations = [
    { id: 1, name: 'Shopify', description: 'Connect your online store', category: 'ecommerce', popular: true, logo: '🛒' },
    { id: 2, name: 'HubSpot', description: 'Customer relationship management', category: 'crm', popular: true, logo: '🤝' },
    { id: 3, name: 'Mailchimp', description: 'Email marketing platform', category: 'marketing', popular: true, logo: '📧' },
    { id: 4, name: 'Stripe', description: 'Online payment processing', category: 'finance', popular: true, logo: '💳' },
    { id: 5, name: 'WooCommerce', description: 'WordPress eCommerce plugin', category: 'ecommerce', popular: false, logo: '🛍️' },
    { id: 6, name: 'Salesforce', description: 'Customer success platform', category: 'crm', popular: true, logo: '☁️' },
    { id: 7, name: 'QuickBooks', description: 'Accounting software', category: 'finance', popular: false, logo: '📊' },
    { id: 8, name: 'Google Analytics', description: 'Web analytics service', category: 'marketing', popular: true, logo: '📈' },
    { id: 9, name: 'PayPal', description: 'Online payment system', category: 'finance', popular: true, logo: '💵' },
    { id: 10, name: 'Magento', description: 'E-commerce platform', category: 'ecommerce', popular: false, logo: '🏪' },
    { id: 11, name: 'Zendesk', description: 'Customer service platform', category: 'crm', popular: false, logo: '🎧' },
    { id: 12, name: 'Mailerlite', description: 'Email marketing tool', category: 'marketing', popular: false, logo: '📨' },
  ];

  // Filter integrations based on search and category
  const filteredIntegrations = integrations.filter(integration => {
    const matchesSearch = integration.name.toLowerCase().includes(searchTerm.toLowerCase()) || 
                          integration.description.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesCategory = activeCategory === 'all' || integration.category === activeCategory;
    
    return matchesSearch && matchesCategory;
  });

  return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
          {t('integrations.title')}
        </h1>
        <p className={`mt-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-600'}`}>
          {t('integrations.subtitle')}
        </p>
      </div>

      {/* Search and filter */}
      <div className={`p-6 mb-8 rounded-lg shadow-md ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
          <div className="relative w-full md:max-w-xs">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Search size={18} className={theme === 'dark' ? 'text-gray-400' : 'text-gray-500'} />
            </div>
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className={`block w-full pl-10 pr-3 py-2 border rounded-md ${
                theme === 'dark' 
                  ? 'bg-gray-700 border-gray-600 text-white placeholder-gray-400 focus:border-green-500 focus:ring-green-500' 
                  : 'border-gray-300 placeholder-gray-500 focus:border-green-500 focus:ring-green-500'
              }`}
              placeholder={t('integrations.searchPlaceholder')}
            />
          </div>
          
          <div className="flex flex-wrap gap-2">
            {categories.map((category) => (
              <button
                key={category.id}
                onClick={() => setActiveCategory(category.id)}
                className={`px-3 py-1.5 text-sm font-medium rounded-md flex items-center ${
                  activeCategory === category.id
                    ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:bg-opacity-30 dark:text-green-300'
                    : theme === 'dark'
                      ? 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                {category.icon && <category.icon size={16} className="mr-1.5" />}
                {category.name}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Integrations grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {filteredIntegrations.map((integration) => (
          <div 
            key={integration.id} 
            className={`p-6 rounded-lg shadow-md ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}
          >
            <div className="flex items-start space-x-4">
              <div className="flex-shrink-0 w-12 h-12 flex items-center justify-center text-2xl rounded-lg bg-gray-100 dark:bg-gray-700">
                {integration.logo}
              </div>
              <div>
                <div className="flex items-center">
                  <h3 className={`text-lg font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                    {integration.name}
                  </h3>
                  {integration.popular && (
                    <span className="ml-2 px-2 py-0.5 text-xs font-medium rounded-full bg-green-100 text-green-800 dark:bg-green-900 dark:bg-opacity-30 dark:text-green-300">
                      {t('integrations.popular')}
                    </span>
                  )}
                </div>
                <p className={`mt-1 text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-600'}`}>
                  {integration.description}
                </p>
              </div>
            </div>

            <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
              <button 
                className="w-full px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
              >
                {t('integrations.connectButton')}
              </button>
            </div>
          </div>
        ))}
      </div>
      
      {filteredIntegrations.length === 0 && (
        <div className={`text-center py-12 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
          <Database size={48} className="mx-auto mb-4" />
          <p className="text-lg font-medium">{t('integrations.noResults')}</p>
          <p className="mt-1">{t('integrations.tryDifferentSearch')}</p>
        </div>
      )}
    </div>
  );
};

export default Integrations;
