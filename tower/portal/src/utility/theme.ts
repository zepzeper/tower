import { Theme } from '../context/ThemeContext';

// Utility function to get theme-specific styles
export const getThemeStyles = (theme: Theme) => {
  // Common style objects for reuse across components
  const styles = {
    // Page backgrounds
    pageBackground: theme === 'dark' ? 'bg-gray-900 text-white' : 'bg-gray-50 text-gray-900',

    // Card and section styles
    card: theme === 'dark' ? 'bg-gray-800' : 'bg-white',
    cardHeader: theme === 'dark' ? 'border-gray-700' : 'border-gray-200',

    // Form element styles
    input: theme === 'dark'
      ? 'bg-gray-700 border-gray-600 text-white'
      : 'bg-white border-gray-300 text-gray-900',

    // Button styles
    primaryButton: theme === 'dark'
      ? 'bg-green-600 hover:bg-green-700 text-white'
      : 'bg-green-500 hover:bg-green-600 text-white',
    secondaryButton: theme === 'dark'
      ? 'bg-gray-700 hover:bg-gray-600 text-white'
      : 'bg-gray-200 hover:bg-gray-300 text-gray-800',
    dangerButton: 'bg-red-600 text-white hover:bg-red-700',

    // Table styles
    tableHeader: theme === 'dark' ? 'bg-gray-700' : 'bg-gray-50',
    tableRow: theme === 'dark' ? 'bg-gray-800 hover:bg-gray-750' : 'bg-white hover:bg-gray-50',
    tableDivider: theme === 'dark' ? 'divide-gray-700' : 'divide-gray-200',

    // Status badges
    active: theme === 'dark'
      ? 'bg-green-900 bg-opacity-30 text-green-300'
      : 'bg-green-100 text-green-800',
    inactive: theme === 'dark'
      ? 'bg-gray-700 text-gray-300'
      : 'bg-gray-100 text-gray-800',

    // Connection type badges
    inbound: theme === 'dark'
      ? 'bg-blue-900 bg-opacity-30 text-blue-300'
      : 'bg-blue-100 text-blue-800',
    outbound: theme === 'dark'
      ? 'bg-green-900 bg-opacity-30 text-green-300'
      : 'bg-green-100 text-green-800',
    bidirectional: theme === 'dark'
      ? 'bg-purple-900 bg-opacity-30 text-purple-300'
      : 'bg-purple-100 text-purple-800',

    // Alert/message styles
    success: theme === 'dark'
      ? 'bg-green-900/30 border-green-800 text-green-300'
      : 'bg-green-50 border-green-200 text-green-800',
    error: theme === 'dark'
      ? 'bg-red-900/30 border-red-800 text-red-300'
      : 'bg-red-50 border-red-200 text-red-600',

    // Text colors
    text: {
      primary: theme === 'dark' ? 'text-white' : 'text-gray-900',
      secondary: theme === 'dark' ? 'text-gray-400' : 'text-gray-500',
      muted: theme === 'dark' ? 'text-gray-500' : 'text-gray-400'
    },

    // Border colors
    border: theme === 'dark' ? 'border-gray-700' : 'border-gray-200',

    // Tab navigation
    tab: {
      active: theme === 'dark'
        ? 'border-green-500 text-green-400'
        : 'border-green-600 text-green-600',
      inactive: theme === 'dark'
        ? 'border-transparent text-gray-400 hover:text-gray-300'
        : 'border-transparent text-gray-500 hover:text-gray-700'
    }
  };

  return styles;
};
