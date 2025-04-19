import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useTheme } from '../context/ThemeContext';
import { UserPlus, Search, MoreHorizontal, CheckCircle, XCircle, Edit, Trash, UserX } from 'lucide-react';

const Users = () => {
  const { t } = useTranslation('pages');
  const { theme } = useTheme();
  const [searchTerm, setSearchTerm] = useState('');
  const [currentView, setCurrentView] = useState('active'); // 'active', 'pending', 'inactive'
  const [showModal, setShowModal] = useState(false);
  const [selectedUser, setSelectedUser] = useState(null);

  // Example users data
  const usersData = [
    { id: 1, name: 'John Doe', email: 'john.doe@example.com', role: 'Admin', status: 'active', lastLogin: '2023-05-15T10:30:00' },
    { id: 2, name: 'Jane Smith', email: 'jane.smith@example.com', role: 'Developer', status: 'active', lastLogin: '2023-05-14T16:45:00' },
    { id: 3, name: 'Alice Johnson', email: 'alice.johnson@example.com', role: 'Viewer', status: 'pending', lastLogin: null },
    { id: 4, name: 'Bob Williams', email: 'bob.williams@example.com', role: 'Developer', status: 'active', lastLogin: '2023-05-13T09:20:00' },
    { id: 5, name: 'Charlie Brown', email: 'charlie.brown@example.com', role: 'Admin', status: 'inactive', lastLogin: '2023-05-01T14:15:00' },
    { id: 6, name: 'Diana Prince', email: 'diana.prince@example.com', role: 'Viewer', status: 'active', lastLogin: '2023-05-15T08:30:00' },
    { id: 7, name: 'Edward Norton', email: 'edward.norton@example.com', role: 'Developer', status: 'active', lastLogin: '2023-05-14T11:45:00' },
    { id: 8, name: 'Fiona Apple', email: 'fiona.apple@example.com', role: 'Viewer', status: 'pending', lastLogin: null },
  ];

  // Filter users based on search term and current view
  const filteredUsers = usersData.filter(user => {
    const matchesSearch = user.name.toLowerCase().includes(searchTerm.toLowerCase()) || 
                          user.email.toLowerCase().includes(searchTerm.toLowerCase()) ||
                          user.role.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesView = currentView === 'all' || user.status === currentView;
    
    return matchesSearch && matchesView;
  });

  // Get status indicator
  const getStatusDisplay = (status) => {
    switch (status) {
      case 'active':
        return { 
          icon: <CheckCircle size={16} className="text-green-500" />, 
          text: t('users.statusActive'),
          className: 'text-green-500'
        };
      case 'inactive':
        return { 
          icon: <XCircle size={16} className="text-red-500" />, 
          text: t('users.statusInactive'),
          className: 'text-red-500'
        };
      case 'pending':
        return { 
          icon: <UserX size={16} className="text-yellow-500" />, 
          text: t('users.statusPending'),
          className: 'text-yellow-500'
        };
      default:
        return { 
          icon: <XCircle size={16} className="text-gray-500" />, 
          text: t('users.statusUnknown'),
          className: 'text-gray-500'
        };
    }
  };

  // Format date
  const formatDate = (dateString) => {
    if (!dateString) return t('users.never');
    
    const date = new Date(dateString);
    return new Intl.DateTimeFormat(undefined, {
      dateStyle: 'medium',
      timeStyle: 'short'
    }).format(date);
  };

  // Handle user actions
  const handleUserAction = (action, user) => {
    setSelectedUser(user);
    
    switch (action) {
      case 'edit':
        // Navigate to edit user page or open edit modal
        console.log('Edit user:', user);
        break;
      case 'delete':
        setShowModal(true);
        break;
      default:
        break;
    }
  };

  // Handle confirm delete
  const handleConfirmDelete = () => {
    console.log('Delete user:', selectedUser);
    // In a real app, you would call your API to delete the user
    setShowModal(false);
    setSelectedUser(null);
  };

  // Tabs for user views
  const tabs = [
    { id: 'all', name: t('users.allUsers') },
    { id: 'active', name: t('users.activeUsers') },
    { id: 'pending', name: t('users.pendingUsers') },
    { id: 'inactive', name: t('users.inactiveUsers') },
  ];

  return (
    <div className="p-6">
      <div className="mb-6 flex flex-col sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className={`text-2xl font-bold ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
            {t('navigation.users')}
          </h1>
          <p className={`mt-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-600'}`}>
            {t('users.manageUsersDescription')}
          </p>
        </div>
        <div className="mt-4 sm:mt-0">
          <button 
            className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
          >
            <UserPlus size={16} className="mr-2" />
            {t('users.addUser')}
          </button>
        </div>
      </div>

      {/* Filters and search */}
      <div className={`p-4 mb-6 rounded-lg ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'} shadow-md`}>
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div className="flex flex-wrap gap-2">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setCurrentView(tab.id)}
                className={`px-3 py-1.5 text-sm font-medium rounded-md ${
                  currentView === tab.id
                    ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:bg-opacity-30 dark:text-green-300'
                    : theme === 'dark'
                      ? 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                }`}
              >
                {tab.name}
                {tab.id !== 'all' && (
                  <span className="ml-1.5 px-2 py-0.5 text-xs rounded-full bg-opacity-80 bg-gray-200 text-gray-800 dark:bg-gray-600 dark:text-gray-200">
                    {usersData.filter(user => user.status === tab.id).length}
                  </span>
                )}
              </button>
            ))}
          </div>
          
          <div className="relative max-w-xs sm:w-64">
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
              placeholder={t('users.searchUsers')}
            />
          </div>
        </div>
      </div>

      {/* Users table */}
      <div className={`rounded-lg shadow-md overflow-hidden ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
            <thead className={theme === 'dark' ? 'bg-gray-700' : 'bg-gray-50'}>
              <tr>
                <th scope="col" className={`px-6 py-3 text-left text-xs font-medium ${theme === 'dark' ? 'text-gray-300 uppercase tracking-wider' : 'text-gray-500 uppercase tracking-wider'}`}>
                  {t('users.name')}
                </th>
                <th scope="col" className={`px-6 py-3 text-left text-xs font-medium ${theme === 'dark' ? 'text-gray-300 uppercase tracking-wider' : 'text-gray-500 uppercase tracking-wider'}`}>
                  {t('users.email')}
                </th>
                <th scope="col" className={`px-6 py-3 text-left text-xs font-medium ${theme === 'dark' ? 'text-gray-300 uppercase tracking-wider' : 'text-gray-500 uppercase tracking-wider'}`}>
                  {t('users.role')}
                </th>
                <th scope="col" className={`px-6 py-3 text-left text-xs font-medium ${theme === 'dark' ? 'text-gray-300 uppercase tracking-wider' : 'text-gray-500 uppercase tracking-wider'}`}>
                  {t('users.status')}
                </th>
                <th scope="col" className={`px-6 py-3 text-left text-xs font-medium ${theme === 'dark' ? 'text-gray-300 uppercase tracking-wider' : 'text-gray-500 uppercase tracking-wider'}`}>
                  {t('users.lastLogin')}
                </th>
                <th scope="col" className={`px-6 py-3 text-right text-xs font-medium ${theme === 'dark' ? 'text-gray-300 uppercase tracking-wider' : 'text-gray-500 uppercase tracking-wider'}`}>
                  {t('users.actions')}
                </th>
              </tr>
            </thead>
            <tbody className={`divide-y ${theme === 'dark' ? 'divide-gray-700' : 'divide-gray-200'}`}>
              {filteredUsers.map((user) => {
                const status = getStatusDisplay(user.status);
                
                return (
                  <tr key={user.id} className={theme === 'dark' ? 'bg-gray-800 hover:bg-gray-700' : 'bg-white hover:bg-gray-50'}>
                    <td className={`px-6 py-4 whitespace-nowrap text-sm font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                      {user.name}
                    </td>
                    <td className={`px-6 py-4 whitespace-nowrap text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-500'}`}>
                      {user.email}
                    </td>
                    <td className={`px-6 py-4 whitespace-nowrap text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-500'}`}>
                      <span className={`px-2 py-1 text-xs rounded-full ${
                        user.role === 'Admin'
                          ? 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:bg-opacity-30 dark:text-purple-300'
                          : user.role === 'Developer'
                            ? 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:bg-opacity-30 dark:text-blue-300'
                            : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300'
                      }`}>
                        {user.role}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <div className="flex items-center">
                        {status.icon}
                        <span className={`ml-1.5 ${status.className}`}>
                          {status.text}
                        </span>
                      </div>
                    </td>
                    <td className={`px-6 py-4 whitespace-nowrap text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-500'}`}>
                      {formatDate(user.lastLogin)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <div className="flex justify-end space-x-2">
                        <button 
                          onClick={() => handleUserAction('edit', user)}
                          className={`p-1 rounded-full ${theme === 'dark' ? 'hover:bg-gray-700' : 'hover:bg-gray-100'}`}
                          aria-label={t('users.edit')}
                        >
                          <Edit size={16} className={theme === 'dark' ? 'text-gray-400 hover:text-white' : 'text-gray-500 hover:text-gray-700'} />
                        </button>
                        <button 
                          onClick={() => handleUserAction('delete', user)}
                          className={`p-1 rounded-full ${theme === 'dark' ? 'hover:bg-gray-700' : 'hover:bg-gray-100'}`}
                          aria-label={t('users.delete')}
                        >
                          <Trash size={16} className={theme === 'dark' ? 'text-gray-400 hover:text-red-500' : 'text-gray-500 hover:text-red-500'} />
                        </button>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
        
        {filteredUsers.length === 0 && (
          <div className={`px-6 py-12 text-center ${theme === 'dark' ? 'text-gray-400' : 'text-gray-500'}`}>
            <p className="text-lg font-medium">{t('users.noResults')}</p>
            <p className="mt-1">{t('users.tryDifferentSearch')}</p>
          </div>
        )}
      </div>

      {/* Delete Confirmation Modal */}
      {showModal && (
        <div className="fixed inset-0 z-10 overflow-y-auto">
          <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
            {/* Background overlay */}
            <div 
              className="fixed inset-0 transition-opacity" 
              aria-hidden="true"
              onClick={() => setShowModal(false)}
            >
              <div className="absolute inset-0 bg-gray-500 opacity-75 dark:bg-gray-900 dark:opacity-80"></div>
            </div>

            {/* Modal panel */}
            <div className={`inline-block align-bottom rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full ${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
              <div className="px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                <div className="sm:flex sm:items-start">
                  <div className="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10">
                    <Trash size={20} className="text-red-600" />
                  </div>
                  <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                    <h3 className={`text-lg leading-6 font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
                      {t('users.deleteUserTitle')}
                    </h3>
                    <div className="mt-2">
                      <p className={`text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-500'}`}>
                        {t('users.deleteUserConfirmation', { name: selectedUser?.name })}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
              <div className="px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                <button
                  type="button"
                  className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-red-600 text-base font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm"
                  onClick={handleConfirmDelete}
                >
                  {t('users.delete')}
                </button>
                <button
                  type="button"
                  className={`mt-3 w-full inline-flex justify-center rounded-md border shadow-sm px-4 py-2 text-base font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm ${
                    theme === 'dark' 
                      ? 'border-gray-500 bg-gray-700 text-gray-200 hover:bg-gray-600 focus:ring-gray-500' 
                      : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50 focus:ring-gray-500'
                  }`}
                  onClick={() => setShowModal(false)}
                >
                  {t('common.cancel')}
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Users;
