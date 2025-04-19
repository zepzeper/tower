import React from 'react';
import { CheckCircle, AlertCircle } from 'lucide-react';

interface MappingAlertsProps {
  saveSuccess: boolean;
  missingFields: string[];
  isDarkMode: boolean;
  t: (key: string) => string;
}

const MappingAlerts: React.FC<MappingAlertsProps> = ({
  saveSuccess,
  missingFields,
  isDarkMode,
  t
}) => {
  return (
    <>
      {saveSuccess && (
        <div className={`mb-6 p-4 rounded-lg flex items-center ${
          isDarkMode ? 'bg-green-900 bg-opacity-30 text-green-300' : 'bg-green-50 text-green-800'
        }`}>
          <CheckCircle size={20} className="mr-2" />
          <span>{t('mappings.saveSuccess')}</span>
        </div>
      )}

      {missingFields.length > 0 && (
        <div className={`mt-6 p-4 rounded-lg shadow-md flex items-start ${
          isDarkMode ? 'bg-yellow-900 bg-opacity-30 text-yellow-300' : 'bg-yellow-50 text-yellow-800'
        }`}>
          <AlertCircle size={20} className="mr-3 flex-shrink-0 mt-0.5" />
          <div>
            <h3 className="font-medium">{t('mappings.missingRequiredFields')}</h3>
            <p className="mt-1 text-sm">{t('mappings.requiredFieldsWarning')}</p>
            <ul className="mt-2 ml-5 list-disc text-sm">
              {missingFields.map(field => <li key={field}>{field}</li>)}
            </ul>
          </div>
        </div>
      )}
    </>
  );
};

export default MappingAlerts;
