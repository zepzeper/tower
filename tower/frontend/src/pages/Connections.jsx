import { useState } from 'react';
import Layout from '../components/Layout';
import { useTranslation } from 'react-i18next';
import { Sparkles, Zap, MoveRight} from 'lucide-react';

const Connections = () => {
  const { t } = useTranslation('pages');

  const categories = [
    { id: 'all', name: t('integrations.categories.allCategories') },
    { id: 'ecommerce', name: t('integrations.categories.ecommerce') },
    { id: 'crm', name: t('integrations.categories.crm') },
    { id: 'finance', name: t('integrations.categories.finance') },
    { id: 'marketing', name: t('integrations.categories.marketing') },
  ];

  const [activeTab, setActiveTab] = useState('all');
  const [searchQuery, setSearchQuery] = useState('');

  const integrations = {
    all: [
      { id: 1, name: 'Shopify', category: 'ecommerce' },
      { id: 2, name: 'Magento', category: 'ecommerce' },
      { id: 3, name: 'Woocommerce', category: 'ecommerce' },
      { id: 4, name: 'Salesforce', category: 'crm' },
      { id: 5, name: 'Exact', category: 'finance' },
    ],
    ecommerce: [
      { id: 1, name: 'Shopify', category: 'ecommerce' },
      { id: 2, name: 'Magento', category: 'ecommerce' },
      { id: 3, name: 'Woocommerce', category: 'ecommerce' },
    ],
    crm: [{ id: 4, name: 'Salesforce', category: 'crm' }],
    finance: [{ id: 5, name: 'Exact', category: 'finance' }],
    marketing: [],
  };

  const filteredIntegrations = integrations[activeTab].filter((item) =>
    item.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
      <div className="text-center mb-14 relative">
        <div className="absolute inset-0 bg-gradient-to-r from-green-100 to-emerald-100 opacity-30 -skew-y-3 transform top-1/3 rounded-3xl" />
        <h1 className="text-5xl font-extrabold tracking-tight bg-gradient-to-r from-green-600 to-emerald-600 bg-clip-text text-transparent relative mb-4">
          {t('integrations.title')}
        </h1>
        <p className="text-lg text-gray-600 max-w-xl mx-auto relative">
          {t('integrations.subtitle')}
          <Zap className="w-5 h-5 text-green-400 inline-block ml-2 -mt-1" />
        </p>
      </div>

      <div className="mb-10 max-w-2xl mx-auto relative group">
        <div className="absolute inset-0 bg-gradient-to-r from-green-500 to-emerald-500 rounded-xl blur opacity-20 group-hover:opacity-30 transition duration-300" />
        <input
          type="text"
          placeholder={t('integrations.searchPlaceholder')}
          className="w-full px-6 py-4 border-0 ring-1 ring-gray-200/70 rounded-xl shadow-lg focus:ring-2 focus:ring-green-500 focus:outline-none text-sm relative bg-white/90 backdrop-blur-sm transition-all duration-300 placeholder-gray-400 hover:ring-gray-300"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      <div className="mb-10">
        <nav className="flex flex-wrap gap-3 justify-center">
          {categories.map((category) => (
            <button
              key={category.id}
              onClick={() => setActiveTab(category.id)}
              className={`relative px-5 py-2.5 text-sm font-medium rounded-full transition-all duration-300 ${
                activeTab === category.id
                  ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white shadow-lg'
                  : 'text-gray-600 bg-white hover:bg-gray-50 border border-gray-200 hover:border-gray-300 shadow-sm hover:shadow-md'
              }`}
            >
              {category.name}
            </button>
          ))}
        </nav>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {filteredIntegrations.length === 0 ? (
          <div className="col-span-full text-center py-20 space-y-4 animate-fade-in">
            <div className="inline-flex items-center justify-center mb-4 w-16 h-16 rounded-2xl bg-gradient-to-r from-green-100 to-emerald-100 text-green-600 transition-transform duration-300">
              <Sparkles className="w-8 h-8" />
            </div>
            <p className="text-xl font-semibold text-gray-700">
              {t('integrations.noResults')}
            </p>
            <p className="text-sm text-gray-400">
              {t('integrations.tryDifferentSearch')}
            </p>
          </div>
        ) : (
          filteredIntegrations.map((integration) => (
            <div
              key={integration.id}
              className="bg-white border border-gray-200 p-6 rounded-2xl shadow-sm hover:shadow-lg transition-all duration-300 group cursor-pointer relative overflow-hidden"
            >
              <div className="absolute inset-0 bg-gradient-to-br from-green-50/30 to-emerald-50/30 opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
              <div className="flex items-center space-x-4 relative">
                <div className="w-14 h-14 rounded-xl bg-gradient-to-br from-green-100 to-emerald-100 flex items-center justify-center text-green-600 text-xl font-semibold transition-colors duration-300">
                  {integration.name.charAt(0)}
                </div>
                <div>
                  <h3 className="text-lg font-semibold text-gray-900">
                    {integration.name}
                  </h3>
                </div>
              </div>
              <div className="flex justify-end">
                <MoveRight className="text-green-500 hover:text-emerald-600 transition-colors" />
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default Connections;
