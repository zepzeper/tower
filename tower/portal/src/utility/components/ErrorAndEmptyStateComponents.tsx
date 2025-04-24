import React from 'react';
import { AlertCircle, Info, SearchX, FileX, RefreshCw, Database } from 'lucide-react';
import { useTheme } from '../../context/ThemeContext';

interface ErrorStateProps {
  title?: string;
  message?: string;
  onRetry?: () => void;
  className?: string;
}

export const ErrorState: React.FC<ErrorStateProps> = ({
  title = 'Something went wrong',
  message = 'An error occurred while fetching data. Please try again.',
  onRetry,
  className = '',
}) => {
  const { theme } = useTheme();

  const bgColor = theme === 'dark' ? 'bg-red-900/30' : 'bg-red-50';
  const borderColor = theme === 'dark' ? 'border-red-800' : 'border-red-200';
  const textColor = theme === 'dark' ? 'text-red-300' : 'text-red-600';

  return (
    <div className={`p-6 rounded-lg border ${borderColor} ${bgColor} ${className}`}>
      <div className="flex items-start">
        <AlertCircle className={`${textColor} h-6 w-6 mr-3 flex-shrink-0 mt-0.5`} />
        <div>
          <h3 className={`font-medium ${textColor}`}>{title}</h3>
          <p className={`mt-1 ${textColor} text-sm`}>{message}</p>

          {onRetry && (
            <button
              onClick={onRetry}
              className="mt-4 px-4 py-2 rounded-md bg-gray-600 text-white hover:bg-gray-700 inline-flex items-center"
            >
              <RefreshCw size={16} className="mr-2" />
              Retry
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

interface EmptyStateProps {
  title?: string;
  message?: string;
  icon?: React.ReactNode;
  action?: {
    label: string;
    onClick: () => void;
  };
  className?: string;
}

export const EmptyState: React.FC<EmptyStateProps> = ({
  title = 'No results found',
  message = 'We couldn\'t find any results matching your criteria.',
  icon = <SearchX size={48} />,
  action,
  className = '',
}) => {
  const { theme } = useTheme();

  const textColor = theme === 'dark' ? 'text-gray-300' : 'text-gray-600';
  const iconColor = theme === 'dark' ? 'text-gray-400' : 'text-gray-500';

  return (
    <div className={`p-8 text-center ${className}`}>
      <div className="flex justify-center">
        <div className={`${iconColor} mb-4`}>
          {icon}
        </div>
      </div>
      <h3 className={`text-lg font-medium ${theme === 'dark' ? 'text-white' : 'text-gray-900'}`}>
        {title}
      </h3>
      <p className={`mt-2 ${textColor}`}>{message}</p>

      {action && (
        <button
          onClick={action.onClick}
          className="mt-4 px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 inline-flex items-center"
        >
          {action.label}
        </button>
      )}
    </div>
  );
};

export const NoDataState: React.FC<EmptyStateProps> = (props) => {
  return (
    <EmptyState
      title={props.title || 'No data available'}
      message={props.message || 'There is no data available to display at this time.'}
      icon={props.icon || <Database size={48} />}
      action={props.action}
      className={props.className}
    />
  );
};

export const NoResultsState: React.FC<EmptyStateProps> = (props) => {
  return (
    <EmptyState
      title={props.title || 'No results found'}
      message={props.message || 'We couldn\'t find any results matching your criteria. Try changing your search or filters.'}
      icon={props.icon || <SearchX size={48} />}
      action={props.action}
      className={props.className}
    />
  );
};

export const FileNotFoundState: React.FC<EmptyStateProps> = (props) => {
  return (
    <EmptyState
      title={props.title || 'File not found'}
      message={props.message || 'The requested file could not be found or may have been deleted.'}
      icon={props.icon || <FileX size={48} />}
      action={props.action}
      className={props.className}
    />
  );
};

export const InfoMessage: React.FC<{
  title?: string;
  message: string;
  className?: string;
}> = ({
  title,
  message,
  className = '',
}) => {
    const { theme } = useTheme();

    const bgColor = theme === 'dark' ? 'bg-blue-900/30' : 'bg-blue-50';
    const borderColor = theme === 'dark' ? 'border-blue-800' : 'border-blue-200';
    const textColor = theme === 'dark' ? 'text-blue-300' : 'text-blue-600';

    return (
      <div className={`p-4 rounded-md flex items-start ${bgColor} border ${borderColor} ${className}`}>
        <Info className={`${textColor} h-5 w-5 mr-3 flex-shrink-0 mt-0.5`} />
        <div>
          {title && <h4 className={`font-medium ${textColor}`}>{title}</h4>}
          <p className={`${textColor} ${title ? 'mt-1' : ''} text-sm`}>{message}</p>
        </div>
      </div>
    );
  };
