import { ReactNode } from 'react';
import { ChevronUp, ChevronDown } from 'lucide-react';
import { cn } from '../../lib/cn';

interface TableProps {
  children: ReactNode;
  className?: string;
}

export function Table({ children, className }: TableProps) {
  return (
    <div className="overflow-x-auto border border-zinc-800 rounded-lg">
      <table className={cn('w-full', className)}>
        {children}
      </table>
    </div>
  );
}

interface TableHeadProps {
  children: ReactNode;
}

export function TableHead({ children }: TableHeadProps) {
  return (
    <thead className="border-b border-zinc-800 bg-zinc-900">
      {children}
    </thead>
  );
}

interface TableBodyProps {
  children: ReactNode;
}

export function TableBody({ children }: TableBodyProps) {
  return <tbody className="divide-y divide-zinc-800">{children}</tbody>;
}

interface TableRowProps {
  children: ReactNode;
  isHoverable?: boolean;
  onClick?: () => void;
  className?: string;
}

export function TableRow({
  children,
  isHoverable = true,
  onClick,
  className,
}: TableRowProps) {
  return (
    <tr
      onClick={onClick}
      className={cn(
        isHoverable && 'hover:bg-zinc-800 cursor-pointer',
        onClick && 'cursor-pointer',
        className
      )}
    >
      {children}
    </tr>
  );
}

interface TableCellProps {
  children: ReactNode;
  className?: string;
  colSpan?: number;
}

export function TableCell({ children, className, colSpan }: TableCellProps) {
  return (
    <td
      className={cn(
        'px-4 py-3 text-sm text-zinc-300',
        className
      )}
      colSpan={colSpan}
    >
      {children}
    </td>
  );
}

interface TableHeaderCellProps {
  children: ReactNode;
  sortable?: boolean;
  sorted?: 'asc' | 'desc' | null;
  onSort?: () => void;
  className?: string;
}

export function TableHeaderCell({
  children,
  sortable = false,
  sorted = null,
  onSort,
  className,
}: TableHeaderCellProps) {
  return (
    <th
      className={cn(
        'px-4 py-3 text-left text-xs font-semibold text-zinc-400 uppercase tracking-wider',
        sortable && 'cursor-pointer hover:text-zinc-300',
        className
      )}
      onClick={sortable ? onSort : undefined}
    >
      <div className="flex items-center gap-2">
        {children}
        {sortable && sorted && (
          sorted === 'asc' ? <ChevronUp size={16} /> : <ChevronDown size={16} />
        )}
      </div>
    </th>
  );
}

interface TableEmptyProps {
  message?: string;
}

export function TableEmpty({ message = 'No data found' }: TableEmptyProps) {
  return (
    <tr>
      <td colSpan={100} className="px-4 py-8 text-center text-zinc-400">
        {message}
      </td>
    </tr>
  );
}
