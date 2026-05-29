import { ReactNode } from 'react';
import * as SelectPrimitive from '@radix-ui/react-select';
import { ChevronDown, Check } from 'lucide-react';
import { cn } from '../../lib/cn';

interface SelectProps {
  value: string;
  onValueChange: (value: string) => void;
  label?: string;
  placeholder?: string;
  error?: string;
  disabled?: boolean;
  children: ReactNode;
  className?: string;
}

export function Select({
  value,
  onValueChange,
  label,
  placeholder,
  error,
  disabled,
  children,
  className,
}: SelectProps) {
  return (
    <div className="w-full">
      {label && (
        <label className="block text-sm font-medium text-zinc-300 mb-2">
          {label}
        </label>
      )}
      <SelectPrimitive.Root value={value} onValueChange={onValueChange} disabled={disabled}>
        <SelectPrimitive.Trigger
          className={cn(
            'inline-flex items-center justify-between w-full h-10 rounded-lg border bg-zinc-950 px-3 py-2 text-sm text-white placeholder-zinc-500 outline-none transition disabled:opacity-50 disabled:cursor-not-allowed',
            error
              ? 'border-red-600 focus:border-red-500'
              : 'border-zinc-800 focus:border-indigo-500',
            className
          )}
        >
          <SelectPrimitive.Value placeholder={placeholder} />
          <SelectPrimitive.Icon asChild>
            <ChevronDown size={16} className="opacity-50" />
          </SelectPrimitive.Icon>
        </SelectPrimitive.Trigger>
        <SelectPrimitive.Portal>
          <SelectPrimitive.Content
            position="popper"
            className="z-50 w-[--radix-select-trigger-width] rounded-lg border border-zinc-800 bg-zinc-900 p-1 shadow-lg"
          >
            <SelectPrimitive.Viewport>
              {children}
            </SelectPrimitive.Viewport>
          </SelectPrimitive.Content>
        </SelectPrimitive.Portal>
      </SelectPrimitive.Root>
      {error && <p className="mt-1 text-xs text-red-500">{error}</p>}
    </div>
  );
}

interface SelectItemProps {
  value: string;
  children: ReactNode;
}

export function SelectItem({ value, children }: SelectItemProps) {
  return (
    <SelectPrimitive.Item
      value={value}
      className="relative flex items-center h-10 px-3 py-2 text-sm text-zinc-300 hover:bg-zinc-800 cursor-pointer rounded-md data-[state=checked]:bg-indigo-600 data-[state=checked]:text-white"
    >
      <SelectPrimitive.ItemText>{children}</SelectPrimitive.ItemText>
      <SelectPrimitive.ItemIndicator asChild>
        <Check size={16} className="ml-auto" />
      </SelectPrimitive.ItemIndicator>
    </SelectPrimitive.Item>
  );
}

interface SelectGroupProps {
  label?: string;
  children: ReactNode;
}

export function SelectGroup({ label, children }: SelectGroupProps) {
  return (
    <SelectPrimitive.Group>
      {label && (
        <SelectPrimitive.Label className="px-3 py-2 text-xs font-semibold text-zinc-400 uppercase tracking-wider">
          {label}
        </SelectPrimitive.Label>
      )}
      {children}
    </SelectPrimitive.Group>
  );
}
