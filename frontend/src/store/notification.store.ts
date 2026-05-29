import { create } from 'zustand';

export type ToastType = 'success' | 'error' | 'warning';

export interface Toast {
  id: string;
  message: string;
  type: ToastType;
  duration?: number;
}

interface NotificationState {
  toasts: Toast[];
  addToast: (message: string, type: ToastType, duration?: number) => void;
  removeToast: (id: string) => void;
  success: (message: string) => void;
  error: (message: string) => void;
  warning: (message: string) => void;
}

export const useNotificationStore = create<NotificationState>((set) => ({
  toasts: [],

  addToast: (message, type, duration = 3000) => {
    const id = Math.random().toString(36).substr(2, 9);
    set((state) => ({
      toasts: [...state.toasts, { id, message, type, duration }],
    }));

    if (duration > 0) {
      setTimeout(() => {
        set((state) => ({
          toasts: state.toasts.filter((t) => t.id !== id),
        }));
      }, duration);
    }
  },

  removeToast: (id) => {
    set((state) => ({
      toasts: state.toasts.filter((t) => t.id !== id),
    }));
  },

  success: (message) => useNotificationStore.getState().addToast(message, 'success'),

  error: (message) => useNotificationStore.getState().addToast(message, 'error'),

  warning: (message) => useNotificationStore.getState().addToast(message, 'warning'),
}));
