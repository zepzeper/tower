import React from 'react';
import { Info } from 'lucide-react';

interface MappingPreviewProps {
  testData: Record<string, any> | null;
  previewResult: Record<string, any> | null;
  isDarkMode: boolean;
  t: (key: string) => string;
}

const MappingPreview: React.FC<MappingPreviewProps> = ({
  testData,
  previewResult,
  isDarkMode,
  t
}) => {
  if (!testData || !previewResult) return null;

  return (
    <div className={`mt-6 p-4 rounded-lg shadow-md ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
      <h2 className={`text-lg font-medium mb-4 ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
        {t('mappings.mappingPreview')}
      </h2>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <h3 className={`text-sm font-medium mb-2 ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>
            {t('mappings.sourceData')}
          </h3>
          <pre className={`p-3 rounded text-xs overflow-auto max-h-60 ${
            isDarkMode ? 'bg-gray-900 text-gray-300' : 'bg-gray-100 text-gray-800'
          }`}>
            {JSON.stringify(testData, null, 2)}
          </pre>
        </div>

        <div>
          <h3 className={`text-sm font-medium mb-2 ${isDarkMode ? 'text-gray-300' : 'text-gray-700'}`}>
            {t('mappings.transformedResult')}
          </h3>
          <pre className={`p-3 rounded text-xs overflow-auto max-h-60 ${
            isDarkMode ? 'bg-gray-900 text-gray-300' : 'bg-gray-100 text-gray-800'
          }`}>
            {JSON.stringify(previewResult, null, 2)}
          </pre>
        </div>
      </div>

      <div className={`mt-4 p-3 rounded flex items-start ${
        isDarkMode ? 'bg-blue-900 bg-opacity-30 text-blue-300' : 'bg-blue-50 text-blue-800'
      }`}>
        <Info size={18} className="mr-2 flex-shrink-0 mt-0.5" />
        <div className="text-sm">
          <p className="font-medium">{t('mappings.previewNote')}</p>
          <p className="mt-1">{t('mappings.previewExplanation')}</p>
        </div>
      </div>
    </div>
  );
};

export default MappingPreview;
