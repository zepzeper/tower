import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../context/ThemeContext';
import { Plus, Search, MoreVertical, CheckCircle, XCircle, Clock } from 'lucide-react';
import { Link } from 'react-router-dom';

const Connections = () => {
    const { t } = useTranslation('pages');
    const { theme } = useTheme();
    const [searchTerm, setSearchTerm] = useState('');
    const [isActiveApiModalOpen, setActiveApiModalOpen] = useState(false);

    // Example connection data
    const connections = [
        { id: 1, name: 'Shopify Store', type: 'E-commerce', status: 'active', lastSync: '2023-05-15T10:30:00' },
        { id: 2, name: 'HubSpot CRM', type: 'CRM', status: 'active', lastSync: '2023-05-15T09:45:00' },
        { id: 3, name: 'Mailchimp', type: 'Marketing', status: 'pending', lastSync: null },
        { id: 4, name: 'Stripe Payments', type: 'Finance', status: 'active', lastSync: '2023-05-14T16:20:00' },
        { id: 5, name: 'Google Analytics', type: 'Analytics', status: 'inactive', lastSync: '2023-05-10T11:15:00' },
        { id: 6, name: 'Zendesk Support', type: 'Customer Service', status: 'active', lastSync: '2023-05-15T08:30:00' },
        { id: 7, name: 'QuickBooks', type: 'Finance', status: 'active', lastSync: '2023-05-14T14:45:00' },
        { id: 8, name: 'Salesforce', type: 'CRM', status: 'pending', lastSync: null },
    ];

    // Filter connections based on search term
    const filteredConnections = connections.filter(connection => 
        connection.name.toLowerCase().includes(searchTerm.toLowerCase()) || 
            connection.type.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const handleButtonClick = () => {
        setActiveApiModalOpen(true);
    };

    // Get status icon and color
    const getStatusDisplay = (status) => {
        switch (status) {
            case 'active':
                return { 
                    icon: <CheckCircle size={16} className="text-green-500" />, 
                    text: t('connections.statusActive'),
                    className: 'text-green-500',
                    buttonLabel: t('connections.viewData'),
                    action: 'view', // we'll check this in logic
                    route: `/connections/view/`
                };
            case 'inactive':
            case 'pending':
            default:
                return { 
                    icon: <XCircle size={16} className="text-red-500" />, 
                    text: t('connections.statusInactive'),
                    className: 'text-red-500',
                    buttonLabel: t('connections.activate'),
                    action: 'modal' // instead of navigating
                };
        }
    };

    // Format date
    const formatDate = (dateString) => {
        if (!dateString) return t('connections.never');

        const date = new Date(dateString);
        return new Intl.DateTimeFormat(undefined, {
            dateStyle: 'medium',
            timeStyle: 'short'
        }).format(date);
    };

    return (
        <div className="p-6">
            <div className="mb-6 flex flex-col sm:flex-row sm:items-center sm:justify-between">
                <div>
                    <h1 className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                        {t('navigation.connections')}
                    </h1>
                    <p className={`mt-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-600'}`}>
                        {t('connections.description')}
                    </p>
                </div>
                <div className="mt-4 sm:mt-0">
                    <Link to="/integrations">
                        <button 
                            className="inline-flex items-center cursor-pointer px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                        >
                            <Plus size={16} className="mr-2" />
                            {t('connections.addNew')}
                        </button>
                    </Link>
                </div>
            </div>

            {/* Search and filters */}
            <div className={`p-4 mb-6 rounded-lg ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'} shadow-md`}>
                <div className="relative max-w-xs">
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
                        placeholder={t('connections.searchPlaceholder')}
                    />
                </div>
            </div>

            {/* Connections grid */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                {filteredConnections.map((connection) => {
                    const status = getStatusDisplay(connection.status);

                    return (
                        <div 
                            key={connection.id} 
                            className={`p-6 rounded-lg shadow-md ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}
                        >
                            <div className="flex justify-between items-start">
                                <div>
                                    <h3 className={`text-lg font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                                        {connection.name}
                                    </h3>
                                    <p className={`text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-600'}`}>
                                        {connection.type}
                                    </p>
                                </div>
                            </div>

                            <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
                                <div className="flex justify-between items-center">
                                    <div className="flex items-center">
                                        {status.icon}
                                        <span className={`ml-1.5 text-sm ${status.className}`}>
                                            {status.text}
                                        </span>
                                    </div>
                                    <span className={`text-sm ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                                        {connection.lastSync ? (
                                            <>
                                                <span className="font-medium">{t('connections.lastSync')}:</span> {formatDate(connection.lastSync)}
                                            </>
                                        ) : (
                                                t('connections.never')
                                            )}
                                    </span>
                                </div>
                            </div>

                            <div className="mt-4 pt-4 flex justify-end gap-2">
                                <Link to={`/connections/edit/${connection.id}`}>
                                    <button 
                                        className={`px-4 py-2 text-sm font-medium rounded-md ${
theme === 'dark' 
? 'bg-gray-700 text-white hover:bg-gray-600' 
: 'bg-gray-100 text-gray-700 hover:bg-gray-200'
}`}
                                    >
                                        {t('connections.configure')}
                                    </button>
                                </Link>

                                {status.action === 'view' ? (
                                    <Link to={`${status.route}${connection.id}`}>
                                        <button 
                                            className="px-4 py-2 text-sm font-medium rounded-md bg-green-100 text-green-800 hover:bg-green-200 dark:bg-green-900 dark:bg-opacity-30 dark:text-green-300 dark:hover:bg-opacity-40"
                                        >
                                            {status.buttonLabel}
                                        </button>
                                    </Link>
                                ) : (
                                        <button 
                                            onClick={handleButtonClick}
                                            className="px-4 py-2 text-sm font-medium rounded-md bg-blue-100 text-blue-800 hover:bg-blue-200 dark:bg-blue-900 dark:bg-opacity-30 dark:text-blue-300 dark:hover:bg-opacity-40"
                                        >
                                            {status.buttonLabel}
                                        </button>
                                    )}
                            </div>
                        </div>
                    );
                })}
            </div>

            {filteredConnections.length === 0 && (
                <div className={`text-center py-12 ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
                    <p className="text-lg font-medium">{t('connections.noResults')}</p>
                    <p className="mt-1">{t('connections.tryDifferentSearch')}</p>
                </div>
            )}
        </div>


    );
};

export default Connections;
