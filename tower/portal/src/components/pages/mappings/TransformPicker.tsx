import React from 'react';
import { TransformationType } from '../types';

interface TransformPickerProps {
  mappingId: string;
  currentTransform: string | null;
  transformations: TransformationType[];
  onSelect: (mappingId: string, transform: string | null) => void;
  isDarkMode: boolean;
  showAdvanced: boolean;
}

const TransformPicker: React.FC<TransformPickerProps> = ({
  mappingId,
  currentTransform,
  transformations,
  onSelect,
  isDarkMode,
  showAdvanced
}) => {
  return (
    <div className={`mt-2 p-2 rounded text-sm ${isDarkMode ? 'bg-gray-700' : 'bg-gray-100'}`}>
      <div className="grid grid-cols-2 gap-2">
        {transformations
          .filter(t => !t.isAdvanced)
          .map((transform) => (
            <button
              key={transform.id}
              onClick={() => onSelect(mappingId, transform.id === 'identity' ? null : transform.id)}
              className={`px-2 py-1 rounded text-xs ${
                (currentTransform === transform.id || (!currentTransform && transform.id === 'identity'))
                  ? isDarkMode ? 'bg-green-800 text-green-200' : 'bg-green-200 text-green-800'
                  : isDarkMode ? 'bg-gray-600 text-gray-200 hover:bg-gray-500' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }`}
            >
              {transform.name}
            </button>
        ))}
      </div>

      {showAdvanced && (
        <div className="mt-2 grid grid-cols-2 gap-2">
          {transformations
            .filter(t => t.isAdvanced)
            .map((transform) => (
              <button
                key={transform.id}
                onClick={() => onSelect(mappingId, transform.id)}
                className={`px-2 py-1 rounded text-xs ${
                  currentTransform === transform.id
                    ? isDarkMode ? 'bg-green-800 text-green-200' : 'bg-green-200 text-green-800'
                    : isDarkMode ? 'bg-gray-600 text-gray-200 hover:bg-gray-500' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                }`}
              >
                {transform.name}
              </button>
          ))}
        </div>
      )}
    </div>
  );
};

export default TransformPicker;
