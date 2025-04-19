import React from 'react';
import { Eye, Download, Save, RefreshCw } from 'lucide-react';

interface MappingHeaderProps {
  connectionName: string;
  description: string;
  isSaving: boolean;
  isDarkMode: boolean;
  t: (key: string) => string;
  onSave: () => void;
  onTest: () => void;
  onExport: () => void;
}

const MappingHeader: React.FC<MappingHeaderProps> = ({
  connectionName,
  description,
  isSaving,
  isDarkMode,
  t,
  onSave,
  onTest,
  onExport
}) => {
  return (
    <div className="mb-6 flex flex-col sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 className={`text-2xl font-bold ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
          {connectionName} {t('mappings.title')}
        </h1>
        <p className={`mt-1 ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
          {description}
        </p>
      </div>
      <div className="mt-4 sm:mt-0 flex gap-2">
        <button
          onClick={onTest}
          className={`inline-flex items-center px-4 py-2 border rounded-md shadow-sm text-sm font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 ${
            isDarkMode
              ? 'border-gray-600 bg-gray-700 text-white hover:bg-gray-600'
              : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50'
          }`}
        >
          <Eye size={16} className="mr-2" />
          {t('mappings.testMappings')}
        </button>

        <button
          onClick={onExport}
          className={`inline-flex items-center px-4 py-2 border rounded-md shadow-sm text-sm font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 ${
            isDarkMode
              ? 'border-gray-600 bg-gray-700 text-white hover:bg-gray-600'
              : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50'
          }`}
        >
          <Download size={16} className="mr-2" />
          {t('mappings.exportMappings')}
        </button>

        <button
          onClick={onSave}
          disabled={isSaving}
          className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
        >
          {isSaving ? (
            <RefreshCw size={16} className="animate-spin mr-2" />
          ) : (
            <Save size={16} className="mr-2" />
          )}
          {t('mappings.saveMappings')}
        </button>
      </div>
    </div>
  );
};

export default MappingHeader;
