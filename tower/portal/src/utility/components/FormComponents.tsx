import React, { ReactNode } from 'react';
import { Eye, EyeOff, Search, ChevronDown, Check, X } from 'lucide-react';
import { useTheme } from '../../context/ThemeContext';

// ----------------
// Input Components
// ----------------

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  icon?: ReactNode;
  className?: string;
}

export const Input: React.FC<InputProps> = ({
  label,
  error,
  icon,
  className = '',
  ...props
}) => {
  const { theme } = useTheme();

  const inputClasses = `w-full px-3 ${icon ? 'pl-10' : 'pl-3'} py-2 rounded-md focus:outline-none focus:ring-2 ${theme === 'dark'
    ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
    : 'bg-white border-gray-300 text-gray-900 focus:ring-green-500 focus:border-green-500'
    } ${error ? 'border-red-500 focus:ring-red-500 focus:border-red-500' : ''}`;

  const labelClasses = `block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'
    }`;

  const errorClasses = `mt-1 text-sm ${theme === 'dark' ? 'text-red-400' : 'text-red-600'
    }`;

  return (
    <div className={className}>
      {label && <label className={labelClasses}>{label}</label>}
      <div className="relative">
        {icon && (
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            {icon}
          </div>
        )}
        <input className={inputClasses} {...props} />
      </div>
      {error && <p className={errorClasses}>{error}</p>}
    </div>
  );
};

interface PasswordInputProps extends Omit<InputProps, 'type'> {
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

export const PasswordInput: React.FC<PasswordInputProps> = ({
  label,
  error,
  className = '',
  value,
  onChange,
  ...props
}) => {
  const { theme } = useTheme();
  const [showPassword, setShowPassword] = React.useState(false);

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  const buttonClasses = `absolute inset-y-0 right-0 pr-3 flex items-center focus:outline-none ${theme === 'dark' ? 'text-gray-400 hover:text-gray-300' : 'text-gray-500 hover:text-gray-700'
    }`;

  return (
    <div className={className}>
      {label && (
        <label className={`block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'
          }`}>
          {label}
        </label>
      )}
      <div className="relative">
        <input
          type={showPassword ? 'text' : 'password'}
          value={value}
          onChange={onChange}
          className={`w-full px-3 py-2 pr-10 rounded-md focus:outline-none focus:ring-2 ${theme === 'dark'
            ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
            : 'bg-white border-gray-300 text-gray-900 focus:ring-green-500 focus:border-green-500'
            } ${error ? 'border-red-500 focus:ring-red-500 focus:border-red-500' : ''}`}
          {...props}
        />
        <button
          type="button"
          className={buttonClasses}
          onClick={togglePasswordVisibility}
          tabIndex={-1}
        >
          {showPassword ? <EyeOff size={18} /> : <Eye size={18} />}
        </button>
      </div>
      {error && (
        <p className={`mt-1 text-sm ${theme === 'dark' ? 'text-red-400' : 'text-red-600'
          }`}>
          {error}
        </p>
      )}
    </div>
  );
};

interface SearchInputProps extends Omit<InputProps, 'type' | 'icon'> {
  onSearch?: (value: string) => void;
}

export const SearchInput: React.FC<SearchInputProps> = ({
  label,
  error,
  className = '',
  value,
  onChange,
  onSearch,
  ...props
}) => {
  const { theme } = useTheme();

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && onSearch && typeof value === 'string') {
      onSearch(value);
    }
  };

  const iconColor = theme === 'dark' ? 'text-gray-400' : 'text-gray-500';

  return (
    <Input
      type="text"
      label={label}
      error={error}
      className={className}
      value={value}
      onChange={onChange}
      onKeyDown={handleKeyDown}
      icon={<Search className={`h-5 w-5 ${iconColor}`} />}
      {...props}
    />
  );
};

// ----------------
// Select Component
// ----------------

interface SelectOption {
  value: string;
  label: string;
}

interface SelectProps extends Omit<React.SelectHTMLAttributes<HTMLSelectElement>, 'onChange'> {
  label?: string;
  error?: string;
  options: SelectOption[];
  onChange: (value: string) => void;
  className?: string;
}

export const Select: React.FC<SelectProps> = ({
  label,
  error,
  options,
  onChange,
  value,
  className = '',
  ...props
}) => {
  const { theme } = useTheme();

  const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    onChange(e.target.value);
  };

  const selectClasses = `w-full px-3 py-2 pr-10 rounded-md appearance-none focus:outline-none focus:ring-2 ${theme === 'dark'
    ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
    : 'bg-white border-gray-300 text-gray-900 focus:ring-green-500 focus:border-green-500'
    } ${error ? 'border-red-500 focus:ring-red-500 focus:border-red-500' : ''}`;

  const labelClasses = `block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'
    }`;

  const errorClasses = `mt-1 text-sm ${theme === 'dark' ? 'text-red-400' : 'text-red-600'
    }`;

  return (
    <div className={className}>
      {label && <label className={labelClasses}>{label}</label>}
      <div className="relative">
        <select
          className={selectClasses}
          value={value}
          onChange={handleChange}
          {...props}
        >
          {options.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
        <div className="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
          <ChevronDown className={theme === 'dark' ? 'text-gray-400' : 'text-gray-500'} size={16} />
        </div>
      </div>
      {error && <p className={errorClasses}>{error}</p>}
    </div>
  );
};

// ----------------
// Checkbox and Toggle Components
// ----------------

interface CheckboxProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'type' | 'onChange'> {
  label: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
  className?: string;
}

export const Checkbox: React.FC<CheckboxProps> = ({
  label,
  checked,
  onChange,
  className = '',
  ...props
}) => {
  const { theme } = useTheme();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    onChange(e.target.checked);
  };

  const checkboxClasses = `h-4 w-4 rounded focus:ring-2 focus:ring-offset-2 ${theme === 'dark'
    ? 'text-green-500 focus:ring-green-600 border-gray-700 bg-gray-700 focus:ring-offset-gray-900'
    : 'text-green-600 focus:ring-green-500 border-gray-300 focus:ring-offset-white'
    }`;

  const labelClasses = `ml-2 block text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'
    }`;

  return (
    <div className={`flex items-center ${className}`}>
      <input
        type="checkbox"
        checked={checked}
        onChange={handleChange}
        className={checkboxClasses}
        {...props}
      />
      <label className={labelClasses}>{label}</label>
    </div>
  );
};

interface ToggleProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'type' | 'onChange'> {
  label?: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
  className?: string;
}

export const Toggle: React.FC<ToggleProps> = ({
  label,
  checked,
  onChange,
  className = '',
  ...props
}) => {
  const { theme } = useTheme();

  const handleChange = () => {
    onChange(!checked);
  };

  const toggleClasses = `relative inline-flex flex-shrink-0 h-6 w-11 border-2 border-transparent rounded-full cursor-pointer transition-colors ease-in-out duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 ${checked
    ? 'bg-green-500 focus:ring-green-500'
    : theme === 'dark'
      ? 'bg-gray-700 focus:ring-gray-500'
      : 'bg-gray-200 focus:ring-gray-400'
    } ${theme === 'dark' ? 'focus:ring-offset-gray-900' : 'focus:ring-offset-white'}`;

  const toggleHandleClasses = `pointer-events-none inline-block h-5 w-5 rounded-full bg-white shadow transform ring-0 transition ease-in-out duration-200 ${checked ? 'translate-x-5' : 'translate-x-0'
    }`;

  const labelClasses = `${label ? 'ml-3' : ''} text-sm ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'
    }`;

  return (
    <div className={`flex items-center ${className}`}>
      <button
        type="button"
        className={toggleClasses}
        onClick={handleChange}
        aria-pressed={checked}
        {...props}
      >
        <span className="sr-only">Toggle</span>
        <span className={toggleHandleClasses}></span>
      </button>
      {label && <span className={labelClasses}>{label}</span>}
    </div>
  );
};

// ----------------
// Text Area Component
// ----------------

interface TextAreaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  error?: string;
  className?: string;
}

export const TextArea: React.FC<TextAreaProps> = ({
  label,
  error,
  className = '',
  ...props
}) => {
  const { theme } = useTheme();

  const textareaClasses = `w-full px-3 py-2 rounded-md focus:outline-none focus:ring-2 ${theme === 'dark'
    ? 'bg-gray-700 border-gray-600 text-white focus:ring-green-500 focus:border-green-500'
    : 'bg-white border-gray-300 text-gray-900 focus:ring-green-500 focus:border-green-500'
    } ${error ? 'border-red-500 focus:ring-red-500 focus:border-red-500' : ''}`;

  const labelClasses = `block text-sm font-medium mb-1 ${theme === 'dark' ? 'text-gray-300' : 'text-gray-700'
    }`;

  const errorClasses = `mt-1 text-sm ${theme === 'dark' ? 'text-red-400' : 'text-red-600'
    }`;

  return (
    <div className={className}>
      {label && <label className={labelClasses}>{label}</label>}
      <textarea className={textareaClasses} {...props} />
      {error && <p className={errorClasses}>{error}</p>}
    </div>
  );
};

// ----------------
// Badge Component
// ----------------

type BadgeVariant = 'default' | 'success' | 'warning' | 'error' | 'info';

interface BadgeProps {
  children: React.ReactNode;
  variant?: BadgeVariant;
  className?: string;
}

export const Badge: React.FC<BadgeProps> = ({
  children,
  variant = 'default',
  className = '',
}) => {
  const { theme } = useTheme();

  const variantClasses = {
    default: theme === 'dark'
      ? 'bg-gray-700 text-gray-300'
      : 'bg-gray-100 text-gray-800',
    success: theme === 'dark'
      ? 'bg-green-900 bg-opacity-30 text-green-300'
      : 'bg-green-100 text-green-800',
    warning: theme === 'dark'
      ? 'bg-yellow-900 bg-opacity-30 text-yellow-300'
      : 'bg-yellow-100 text-yellow-800',
    error: theme === 'dark'
      ? 'bg-red-900 bg-opacity-30 text-red-300'
      : 'bg-red-100 text-red-800',
    info: theme === 'dark'
      ? 'bg-blue-900 bg-opacity-30 text-blue-300'
      : 'bg-blue-100 text-blue-800',
  };

  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${variantClasses[variant]} ${className}`}>
      {children}
    </span>
  );
};

// ----------------
// Button Components
// ----------------

type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'danger' | 'success' | 'ghost';
type ButtonSize = 'small' | 'medium' | 'large';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode;
  variant?: ButtonVariant;
  size?: ButtonSize;
  isLoading?: boolean;
  leftIcon?: React.ReactNode;
  rightIcon?: React.ReactNode;
  className?: string;
}
