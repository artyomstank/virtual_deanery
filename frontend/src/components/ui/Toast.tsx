import { X } from 'lucide-react';
import { useNotificationStore } from '../../store/notification.store';

export function ToastContainer() {
  const toasts = useNotificationStore((state) => state.toasts);
  const removeToast = useNotificationStore((state) => state.removeToast);

  return (
    <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-3 pointer-events-none">
      {toasts.map((toast) => (
        <div
          key={toast.id}
          className={`pointer-events-auto flex items-center gap-3 rounded-lg px-4 py-3 shadow-lg animate-in slide-in-from-right-full duration-300 ${
            toast.type === 'success'
              ? 'bg-emerald-900 text-emerald-100 border border-emerald-700'
              : toast.type === 'error'
              ? 'bg-red-900 text-red-100 border border-red-700'
              : 'bg-amber-900 text-amber-100 border border-amber-700'
          }`}
        >
          <span className="text-sm">{toast.message}</span>
          <button
            onClick={() => removeToast(toast.id)}
            className="ml-2 inline-flex flex-shrink-0"
          >
            <X size={16} />
          </button>
        </div>
      ))}
    </div>
  );
}
