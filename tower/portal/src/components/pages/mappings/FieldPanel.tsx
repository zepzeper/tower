import React from 'react';
import { Search, ChevronDown, ChevronRight } from 'lucide-react';

interface FieldDefinition {
  id: string;
  name: string;
  type: string;
  path: string;
  sample?: string;
  required?: boolean;
}

interface FieldPanelProps {
  title: string;
  type: 'source' | 'target';
  fields: FieldDefinition[];
  mappings: any[];
  search: string;
  isDarkMode: boolean;
  expanded: Record<string, boolean>;
  onSearchChange: (value: string) => void;
  onDragStart: (field: FieldDefinition, type: 'source' | 'target') => void;
  onDragEnd: () => void;
  toggleExpansion: (fieldId: string, type: 'source' | 'target') => void;
  getFieldById: (id: string, type: 'source' | 'target') => FieldDefinition | undefined;
  t: (key: string) => string;
  addMapping?: (sourceField: FieldDefinition, targetField: FieldDefinition) => void;
}

const FieldPanel: React.FC<FieldPanelProps> = ({
  title,
  type,
  fields,
  mappings,
  search,
  isDarkMode,
  expanded,
  onSearchChange,
  onDragStart,
  onDragEnd,
  toggleExpansion,
  getFieldById,
  t,
  addMapping
}) => {
  const filteredFields = fields.filter(field =>
    field.name.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className={`p-4 rounded-lg shadow-md ${isDarkMode ? 'bg-gray-800' : 'bg-white'}`}>
      <h2 className={`text-lg font-medium mb-4 ${isDarkMode ? 'text-white' : 'text-gray-900'}`}>
        {title}
      </h2>

      <div className="relative mb-4">
        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <Search size={18} className={isDarkMode ? 'text-gray-400' : 'text-gray-500'} />
        </div>
        <input
          type="text"
          value={search}
          onChange={(e) => onSearchChange(e.target.value)}
          className={`block w-full pl-10 pr-3 py-2 border rounded-md ${
            isDarkMode
              ? 'bg-gray-700 border-gray-600 text-white placeholder-gray-400'
              : 'border-gray-300 text-gray-700 placeholder-gray-500'
          }`}
          placeholder={t('mappings.searchFields')}
        />
      </div>

      <div className="overflow-y-auto max-h-96">
        {filteredFields.map((field) => {
          const isUsed = mappings.some(m => m[type === 'source' ? 'sourceField' : 'targetField'] === field.id);
          const isExpanded = expanded[field.id] || false;

          const sourceField =
            type === 'target' && isUsed
              ? getFieldById(mappings.find(m => m.targetField === field.id)?.sourceField, 'source')
              : null;

          return (
            <div
              key={field.id}
              className={`p-3 mb-2 rounded-md border ${
                isUsed
                  ? isDarkMode
                    ? 'border-green-800 bg-green-900 bg-opacity-20'
                    : 'border-green-200 bg-green-50'
                  : field.required
                    ? isDarkMode
                      ? 'border-orange-800 bg-orange-900 bg-opacity-10'
                      : 'border-orange-200 bg-orange-50'
                    : isDarkMode
                      ? 'border-gray-700 hover:border-gray-600'
                      : 'border-gray-200 hover:border-gray-300'
              }`}
              draggable
              onDragStart={() => onDragStart(field, type)}
              onDragEnd={onDragEnd}
            >
              <div className="flex justify-between items-start">
                <div>
                  <div className="font-medium text-sm flex items-center">
                    <button
                      onClick={() => toggleExpansion(field.id, type)}
                      className="mr-1 p-0.5 rounded hover:bg-gray-200 dark:hover:bg-gray-700"
                    >
                      {isExpanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
                    </button>
                    <span className={isDarkMode ? 'text-gray-200' : 'text-gray-800'}>
                      {field.name}
                    </span>

                    {type === 'target' && field.required && (
                      <span className={`ml-2 px-1.5 py-0.5 text-xs rounded ${
                        isUsed
                          ? isDarkMode ? 'bg-green-900 text-green-300' : 'bg-green-100 text-green-800'
                          : isDarkMode ? 'bg-orange-900 text-orange-300' : 'bg-orange-100 text-orange-800'
                      }`}>
                        {isUsed ? t('mappings.mapped') : t('mappings.required')}
                      </span>
                    )}

                    {type === 'target' && !field.required && isUsed && (
                      <span className={`ml-2 px-1.5 py-0.5 text-xs rounded ${
                        isDarkMode ? 'bg-green-900 text-green-300' : 'bg-green-100 text-green-800'
                      }`}>
                        {t('mappings.mapped')}
                      </span>
                    )}

                    {type === 'source' && isUsed && (
                      <span className={`ml-2 px-1.5 py-0.5 text-xs rounded ${
                        isDarkMode ? 'bg-green-900 text-green-300' : 'bg-green-100 text-green-800'
                      }`}>
                        {t('mappings.mapped')}
                      </span>
                    )}
                  </div>

                  <div className={`text-xs mt-1 ${isDarkMode ? 'text-gray-400' : 'text-gray-500'}`}>
                    {t('mappings.type')}: {field.type}
                  </div>

                  {isExpanded && (
                    <div className={`text-xs mt-0.5 ${isDarkMode ? 'text-gray-400' : 'text-gray-500'}`}>
                      {t('mappings.path')}: {field.path}
                    </div>
                  )}

                  {isExpanded && type === 'target' && sourceField && (
                    <div className={`mt-2 text-xs ${isDarkMode ? 'text-gray-300' : 'text-gray-600'}`}>
                      {t('mappings.mappedFrom')}: <span className="font-medium">{sourceField.name}</span>
                    </div>
                  )}
                </div>

                {!isUsed && type === 'target' && addMapping && (
                  <div className="ml-2">
                    <button
                      onClick={() => {
                        const matchingSource = fields.find(f =>
                          f.name.toLowerCase() === field.name.toLowerCase()
                        );
                        if (matchingSource) addMapping(matchingSource, field);
                      }}
                      className={`text-xs px-2 py-1 rounded ${
                        isDarkMode
                          ? 'bg-blue-800 text-blue-200 hover:bg-blue-700'
                          : 'bg-blue-100 text-blue-800 hover:bg-blue-200'
                      }`}
                    >
                      {t('mappings.autoMap')}
                    </button>
                  </div>
                )}
              </div>
            </div>
          );
        })}

        {filteredFields.length === 0 && (
          <div className={`text-center py-4 ${isDarkMode ? 'text-gray-400' : 'text-gray-500'}`}>
            {t('mappings.noFieldsFound')}
          </div>
        )}
      </div>
    </div>
  );
};

export default FieldPanel;
