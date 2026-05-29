import { ReactNode } from 'react';
import { cn } from '../../lib/cn';

interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  isLoading?: boolean;
  children: ReactNode;
}

const variants = {
  primary: 'bg-indigo-600 hover:bg-indigo-700 text-white',
  secondary: 'bg-zinc-800 hover:bg-zinc-700 text-zinc-100',
  danger: 'bg-red-600 hover:bg-red-700 text-white',
  ghost: 'hover:bg-zinc-800 text-zinc-300 hover:text-white',
};

const sizes = {
  sm: 'px-3 py-1.5 text-sm',
  md: 'px-4 py-2 text-sm',
  lg: 'px-6 py-3 text-base',
};

export function Button({
  variant = 'primary',
  size = 'md',
  isLoading = false,
  disabled,
  children,
  className,
  ...props
}: ButtonProps) {
  return (
    <button
      disabled={disabled || isLoading}
      className={cn(
        'inline-flex items-center justify-center rounded-lg font-medium transition disabled:opacity-50 disabled:cursor-not-allowed',
        variants[variant],
        sizes[size],
        className
      )}
      {...props}
    >
      {isLoading ? (
        <div className="inline-flex items-center gap-2">
          <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
          {children}
        </div>
      ) : (
        children
      )}
    </button>
  );
}

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  error?: string;
  label?: string;
}

export function Input({ error, label, className, ...props }: InputProps) {
  return (
    <div className="w-full">
      {label && (
        <label className="block text-sm font-medium text-zinc-300 mb-2">
          {label}
        </label>
      )}
      <input
        className={cn(
          'h-10 w-full rounded-lg border bg-zinc-950 px-3 py-2 text-sm text-white placeholder-zinc-500 outline-none transition',
          error
            ? 'border-red-600 focus:border-red-500'
            : 'border-zinc-800 focus:border-indigo-500',
          className
        )}
        {...props}
      />
      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}
    </div>
  );
}

interface BadgeProps {
  variant?: 'success' | 'error' | 'warning' | 'info';
  children: ReactNode;
}

const badgeVariants = {
  success: 'bg-emerald-900 text-emerald-100 border border-emerald-700',
  error: 'bg-red-900 text-red-100 border border-red-700',
  warning: 'bg-amber-900 text-amber-100 border border-amber-700',
  info: 'bg-blue-900 text-blue-100 border border-blue-700',
};

export function Badge({ variant = 'info', children }: BadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
        badgeVariants[variant]
      )}
    >
      {children}
    </span>
  );
}
