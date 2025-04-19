import React from 'react';
import { ArrowRight, Trash2, ChevronRight, ChevronDown, RefreshCw } from 'lucide-react';

interface Mapping {
  id: string;
  sourceField: string;
  targetField: string;
  transform: string | null;
}

interface Field {
  id: string;
  name: string;
  type: string;
  path: string;
  required?: boolean;
}

interface Transformation {
  id: string;
  name: string;
  description: string;
  isAdvanced?: boolean;
}

interface MappingPanelProps {
  mappings: Mapping[];
  sourceFields: Field[];
  targetFields: Field[];
  transformations: Transformation[];
  activeTransformation: string | null;
  isDarkMode: boolean;
  showAdvanced: boolean;
  mappingHistory: Mapping[][];
  t: (key: string) => string;
  getFieldById: (id: string, type: 'source' | 'target') => Field | undefined;
  getTransformName: (id: string | null) => string;
  setActiveTransformation: (id: string | null) => void;
  updateTransformation: (mappingId: string, transform: string | null) => void;
  removeMapping: (id: string) => void;
  undoMappingChange: () => void;
  autoMapFields: () => void;
  setShowAdvanced: (value: boolean) => void;
}

const MappingPanel: React.FC<MappingPanelProps> = ({
  mappings,
  sourceFields,
  targetFields,
  transformations,
  activeTransformation,
  isDarkMode,
  showAdvanced,
  mappingHistory,
  t,
  getFieldById,
  getTransformName,
  setActiveTransformation,
  updateTransformation,
  removeMapping,
  undoMappingChange,
  autoMapFields,
  setShowAdvanced,
  loadingAutoMapping
}) => {
  return (
    <div className={`p-4 rounded-lg shadow-md ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
      <h2 className={`text-lg font-medium mb-4 ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
        {t('mappings.fieldMappings')}
      </h2>

      <div className="flex items-center justify-between mb-4">
        <button
          onClick={() => setShowAdvanced(!showAdvanced)}
          className={`text-sm flex items-center ${isDarkMode ? 'text-blue-400 hover:text-blue-300' : 'text-blue-600 hover:text-blue-700'}`}
        >
          {showAdvanced ? <ChevronDown size={16} className="mr-1" /> : <ChevronRight size={16} className="mr-1" />}
          {t('mappings.advancedOptions')}
        </button>

        <button
          onClick={autoMapFields}
          className={`text-sm ${isDarkMode ? 'text-green-400 hover:text-green-300' : 'text-green-600 hover:text-green-700'}`}
        >
          {t('mappings.autoMapFields')}
        </button>
      </div>

      <div className="overflow-y-auto max-h-96">
        {loadingAutoMapping ? (

                    <div className="p-6 flex justify-center items-center">
                        <div className="text-center">
                            <RefreshCw size={28} className={`animate-spin mb-4 mx-auto ${isDarkMode ? 'text-green-400' : 'text-green-600'}`} />
                            <p className={isDarkMode ? 'text-gray-300' : 'text-gray-700'}>
                                {t('mappings.loading')}
                            </p>
                        </div>
                    </div>

        ): 
          mappings.length > 0 ? mappings.map(mapping => {
          const sourceField = getFieldById(mapping.sourceField, 'source');
          const targetField = getFieldById(mapping.targetField, 'target');
          if (!sourceField || !targetField) return null;

          const isActive = activeTransformation === mapping.id;

          return (
            <div key={mapping.id} className="p-3 mb-3 rounded-md border">
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <div className={`mb-2 ${isDarkMode ? 'text-gray-200' : 'text-gray-800'}`}>
                    <span className="font-medium">{sourceField.name}</span>
                    <ArrowRight size={14} className="inline mx-2" />
                    <span className="font-medium">{targetField.name}</span>
                    {targetField.required && (
                      <span className={`ml-2 px-1.5 py-0.5 text-xs rounded ${
                        isDarkMode ? 'bg-blue-900 text-blue-300' : 'bg-blue-100 text-blue-800'
                      }`}>
                        {t('mappings.required')}
                      </span>
                    )}
                  </div>

                  <div className={`flex items-center text-xs ${isDarkMode ? 'text-gray-400' : 'text-gray-500'}`}>
                    <span className="font-medium">{t('mappings.transformation')}:</span>
                    <span className={`ml-1 ${mapping.transform ? (isDarkMode ? 'text-green-400' : 'text-green-600') : ''}`}>
                      {getTransformName(mapping.transform)}
                    </span>
                    <button
                      onClick={() => setActiveTransformation(mapping.id)}
                      className={`ml-2 underline text-xs ${isDarkMode ? 'text-blue-400 hover:text-blue-300' : 'text-blue-600 hover:text-blue-500'}`}
                    >
                      {mapping.transform ? t('mappings.change') : t('mappings.add')}
                    </button>
                  </div>

                  {isActive && (
                    <div className={`mt-2 p-2 rounded text-sm ${isDarkMode ? 'bg-gray-700' : 'bg-gray-100'}`}>
                      <div className="grid grid-cols-2 gap-2">
                        {transformations.filter(t => !t.isAdvanced).map(tf => (
                          <button
                            key={tf.id}
                            onClick={() => updateTransformation(mapping.id, tf.id === 'identity' ? null : tf.id)}
                            className={`px-2 py-1 rounded text-xs ${
                              mapping.transform === tf.id || (!mapping.transform && tf.id === 'identity')
                                ? isDarkMode ? 'bg-green-800 text-green-200' : 'bg-green-200 text-green-800'
                                : isDarkMode ? 'bg-gray-600 text-gray-200 hover:bg-gray-500' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                            }`}
                          >
                            {tf.name}
                          </button>
                        ))}
                      </div>

                      {showAdvanced && (
                        <div className="mt-2 grid grid-cols-2 gap-2">
                          {transformations.filter(t => t.isAdvanced).map(tf => (
                            <button
                              key={tf.id}
                              onClick={() => updateTransformation(mapping.id, tf.id)}
                              className={`px-2 py-1 rounded text-xs ${
                                mapping.transform === tf.id
                                  ? isDarkMode ? 'bg-green-800 text-green-200' : 'bg-green-200 text-green-800'
                                  : isDarkMode ? 'bg-gray-600 text-gray-200 hover:bg-gray-500' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                              }`}
                            >
                              {tf.name}
                            </button>
                          ))}
                        </div>
                      )}
                    </div>
                  )}
                </div>

                <button
                  onClick={() => removeMapping(mapping.id)}
                  className={`p-1 rounded-full ${
                    isDarkMode
                      ? 'text-gray-400 hover:text-red-400 hover:bg-gray-700'
                      : 'text-gray-500 hover:text-red-500 hover:bg-gray-100'
                  }`}
                >
                  <Trash2 size={16} />
                </button>
              </div>
            </div>
          );
        }) : (
          <div className={`text-center py-8 ${isDarkMode ? 'text-gray-400' : 'text-gray-500'}`}>
            <ArrowRight size={24} className="mx-auto mb-2" />
            <p>{t('mappings.noMappingsYet')}</p>
            <p className="mt-1 text-sm">{t('mappings.dragFieldsToMap')}</p>
          </div>
        )}
      </div>

      {mappingHistory.length > 0 && (
        <div className="mt-3 flex justify-end">
          <button
            onClick={undoMappingChange}
            className={`px-2 py-1 text-xs rounded ${
              isDarkMode ? 'bg-gray-700 text-gray-300 hover:bg-gray-600' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            }`}
          >
            {t('mappings.undoLastChange')}
          </button>
        </div>
      )}
    </div>
  );
};

export default MappingPanel;
