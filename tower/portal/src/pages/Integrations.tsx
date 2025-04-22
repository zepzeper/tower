import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../context/ThemeContext';
import { Search, Database } from 'lucide-react';
import ConnectionModal from '../components/pages/connections/ConnectionModal';
import { connectionService } from '../services/connectionService';

interface Integration {
  id: string;
  name: string;
  description: string;
  category: string;
  categoryName: string;
  popular: boolean;
  imageUrl: string;
  type: string;
}

const Integrations: React.FC = () => {
  const { t } = useTranslation('pages');
  const { theme } = useTheme();
  const [searchTerm, setSearchTerm] = useState<string>('');
  const [activeCategory, setActiveCategory] = useState<string>('all');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedIntegration, setSelectedIntegration] = useState<Integration | null>(null);
  const [integrations, setIntegrations] = useState<Integration[]>([]);
  const [categories, setCategories] = useState<{ id: string, name: string }[]>([]);
  const [loading, setLoading] = useState(true);

  // Load connection types on component mount
  useEffect(() => {
    const loadConnectionTypes = async () => {
      try {
        const types = await connectionService.getConnectionTypesFromFile();
        const transformed = types.map(type => ({
          id: type.id,
          name: type.name,
          description: type.description,
          category: type.category || 'other',
          categoryName: type.categoryName || 'Other',
          popular: type.popular || false,
          imageUrl: type.imageUrl || getDefaultImageUrl(type.id),
          type: type.id
        }));

        setIntegrations(transformed);

        // Generate categories from unique category values
        const uniqueCategoryIds = Array.from(
          new Set(types.map(t => t.category || 'other'))
        );

        // Then map to category objects
        const uniqueCategories = uniqueCategoryIds.map(categoryId => ({
          id: categoryId,
          name: types.find(t => t.category === categoryId)?.categoryName ||
            categoryId.charAt(0).toUpperCase() + categoryId.slice(1) // Default name if not found
        }));

        setCategories([{ id: 'all', name: t('integrations.categories.allCategories') }, ...uniqueCategories]);
      } catch (error) {
        console.error('Failed to load connection types:', error);
      } finally {
        setLoading(false);
      }
    };

    loadConnectionTypes();
  }, [t]);

  const getDefaultImageUrl = (typeId: string): string => {
    // Fallback to default image if none provided
    return `/images/integrations/${typeId}.png`;
  };

  const filteredIntegrations = integrations.filter((integration) => {
    const matchesSearch = integration.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      integration.description.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesCategory = activeCategory === 'all' || integration.category === activeCategory;
    return matchesSearch && matchesCategory;
  });

  const handleConnectClick = (integration: Integration) => {
    setSelectedIntegration(integration);
    setIsModalOpen(true);
  };

  const handleModalClose = () => {
    setIsModalOpen(false);
    setSelectedIntegration(null);
  };

  const handleConnectionSuccess = () => {
    setIsModalOpen(false);
    setSelectedIntegration(null);
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Database className="animate-spin" size={32} />
      </div>
    );
  }

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

      {/* Search + Filters */}
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
              className={`block w-full pl-10 pr-3 py-2 border rounded-md ${theme === 'dark'
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
                className={`px-3 py-1.5 text-sm font-medium rounded-md ${activeCategory === category.id
                  ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:bg-opacity-30 dark:text-green-300'
                  : theme === 'dark'
                    ? 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
              >
                {category.name}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Integrations Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {filteredIntegrations.map((integration) => (
          <div
            key={integration.id}
            className={`p-6 rounded-lg shadow-md ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}
          >
            <div className="flex flex-col items-center">
              <div className="w-40 h-20 mb-4 flex items-center justify-center">
                <img
                  src={integration.imageUrl}
                  alt={integration.name}
                  className="max-w-full max-h-full object-contain"
                />
              </div>
              <div className="text-center">
                <div className="flex items-center justify-center">
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
                <div className="mt-2">
                  <span className={`text-xs px-2 py-1 rounded ${theme === 'dark' ? 'bg-gray-700 text-gray-300' : 'bg-gray-100 text-gray-700'}`}>
                    {integration.categoryName}
                  </span>
                </div>
              </div>
            </div>

            <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
              <button
                onClick={() => handleConnectClick(integration)}
                className="w-full px-4 py-2 text-sm font-medium text-white bg-green-600 rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
              >
                {t('integrations.connectButton')}
              </button>
            </div>
          </div>
        ))}
      </div>

      {filteredIntegrations.length === 0 && !loading && (
        <div className={`text-center py-12 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
          <Database size={48} className="mx-auto mb-4" />
          <p className="text-lg font-medium">{t('integrations.noResults')}</p>
          <p className="mt-1">{t('integrations.tryDifferentSearch')}</p>
        </div>
      )}

      {/* Connection Modal */}
      {selectedIntegration && (
        <ConnectionModal
          isOpen={isModalOpen}
          connectionType={selectedIntegration.type}
          onClose={handleModalClose}
          onSuccess={handleConnectionSuccess}
        />
      )}
    </div>
  );
};

export default Integrations;
