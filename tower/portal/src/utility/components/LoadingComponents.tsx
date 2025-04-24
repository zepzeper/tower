import React from 'react';
import { Loader2 } from 'lucide-react';
import { useTheme } from '../../context/ThemeContext';

interface LoadingProps {
  size?: 'small' | 'medium' | 'large';
  text?: string;
  fullScreen?: boolean;
  className?: string;
}

export const LoadingSpinner: React.FC<LoadingProps> = ({
  size = 'medium',
  text,
  fullScreen = false,
  className = '',
}) => {
  const { theme } = useTheme();

  const sizeMap = {
    small: 'w-4 h-4',
    medium: 'w-8 h-8',
    large: 'w-12 h-12',
  };

  const spinnerSize = sizeMap[size];
  const textColor = theme === 'dark' ? 'text-gray-300' : 'text-gray-600';

  if (fullScreen) {
    return (
      <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
        <div className="bg-white dark:bg-gray-800 rounded-lg p-6 flex flex-col items-center shadow-xl">
          <Loader2 className={`${spinnerSize} text-green-500 animate-spin ${className}`} />
          {text && <p className={`mt-4 ${textColor}`}>{text}</p>}
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col items-center justify-center">
      <Loader2 className={`${spinnerSize} text-green-500 animate-spin ${className}`} />
      {text && <p className={`mt-2 ${textColor}`}>{text}</p>}
    </div>
  );
};

export const LoadingCard: React.FC<{
  lines?: number;
  className?: string;
}> = ({ lines = 5, className = '' }) => {
  const { theme } = useTheme();

  const bgColor = theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200';

  return (
    <div className={`rounded-lg border ${theme === 'dark' ? 'border-gray-700' : 'border-gray-200'} p-6 ${className}`}>
      <div className={`h-5 ${bgColor} rounded animate-pulse w-3/4 mb-4`}></div>
      {Array.from({ length: lines }).map((_, i) => (
        <div
          key={i}
          className={`h-4 ${bgColor} rounded animate-pulse mb-2 ${i % 3 === 0 ? 'w-full' : i % 3 === 1 ? 'w-5/6' : 'w-4/6'}`}
          style={{ animationDelay: `${i * 100}ms` }}
        ></div>
      ))}
    </div>
  );
};

export const LoadingTable: React.FC<{
  rows?: number;
  columns?: number;
  className?: string;
}> = ({ rows = 5, columns = 4, className = '' }) => {
  const { theme } = useTheme();

  const bgColor = theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200';
  const borderColor = theme === 'dark' ? 'border-gray-700' : 'border-gray-200';

  return (
    <div className={`overflow-hidden rounded-lg border ${borderColor} ${className}`}>
      <div className={`${theme === 'dark' ? 'bg-gray-800' : 'bg-white'}`}>
        <div className="grid" style={{ gridTemplateColumns: `repeat(${columns}, minmax(0, 1fr))` }}>
          {/* Header */}
          {Array.from({ length: columns }).map((_, i) => (
            <div
              key={`header-${i}`}
              className={`px-6 py-3 border-b ${borderColor}`}
            >
              <div className={`h-4 ${bgColor} rounded animate-pulse w-2/3`}></div>
            </div>
          ))}

          {/* Rows */}
          {Array.from({ length: rows }).map((_, rowIndex) => (
            <>
              {Array.from({ length: columns }).map((_, colIndex) => (
                <div
                  key={`cell-${rowIndex}-${colIndex}`}
                  className={`px-6 py-4 border-b ${borderColor}`}
                >
                  <div
                    className={`h-4 ${bgColor} rounded animate-pulse ${colIndex % 3 === 0 ? 'w-full' : colIndex % 3 === 1 ? 'w-3/4' : 'w-1/2'}`}
                    style={{ animationDelay: `${(rowIndex * columns + colIndex) * 50}ms` }}
                  ></div>
                </div>
              ))}
            </>
          ))}
        </div>
      </div>
    </div>
  );
};

export const LoadingGrid: React.FC<{
  items?: number;
  className?: string;
}> = ({ items = 6, className = '' }) => {
  return (
    <div className={`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 ${className}`}>
      {Array.from({ length: items }).map((_, i) => (
        <LoadingCard key={i} lines={3} />
      ))}
    </div>
  );
};

export const LoadingPage: React.FC<{
  text?: string;
}> = ({ text = 'Loading...' }) => {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <LoadingSpinner size="large" text={text} />
    </div>
  );
};
