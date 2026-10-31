import { ApiError, apiClient } from './api';

export interface AppNotification {
  id: string;
  type: string;
  title: string;
  body: string;
  read: boolean;
  createdAt: string;
  displayTime: string;
}

export interface NotificationList {
  items: AppNotification[];
  unreadCount: number;
}

/** Local mock used when the backend is unreachable or auth is unavailable. */
const MOCK_NOTIFICATIONS: AppNotification[] = [
  {
    id: 'notif_demo_3',
    type: 'in_app',
    title: 'Recurring billing run completed',
    body: 'Nightly billing processed 3 invoices and collected $1,920.00.',
    read: false,
    createdAt: new Date(Date.now() - 2 * 3600_000).toISOString(),
    displayTime: '2h ago',
  },
  {
    id: 'notif_demo_2',
    type: 'in_app',
    title: 'Subscription renewal imminent',
    body: 'Enterprise Annual renews in 15 days for Acme Enterprises Inc.',
    read: false,
    createdAt: new Date(Date.now() - 5 * 3600_000).toISOString(),
    displayTime: '5h ago',
  },
  {
    id: 'notif_demo_1',
    type: 'in_app',
    title: 'Welcome to Brightloop',
    body: 'Your Brightloop workspace is ready. Explore the finance dashboard.',
    read: true,
    createdAt: new Date(Date.now() - 26 * 3600_000).toISOString(),
    displayTime: 'Yesterday',
  },
];

function isFallbackError(error: unknown): boolean {
  if (error instanceof TypeError) return true;
  if (error instanceof ApiError && (error.status === 401 || error.status === 403)) return true;
  return false;
}

function localUnread(): number {
  return MOCK_NOTIFICATIONS.filter((n) => !n.read).length;
}

export const notificationsService = {
  async getNotifications(limit = 20): Promise<NotificationList> {
    try {
      return await apiClient<NotificationList>(`/v1/notifications?limit=${limit}`);
    } catch (error) {
      if (isFallbackError(error)) {
        return { items: MOCK_NOTIFICATIONS, unreadCount: localUnread() };
      }
      throw error;
    }
  },

  async markRead(id: string): Promise<void> {
    try {
      await apiClient<unknown>(`/v1/notifications/${id}/read`, { method: 'POST' });
    } catch (error) {
      if (isFallbackError(error)) {
        const target = MOCK_NOTIFICATIONS.find((n) => n.id === id);
        if (target) target.read = true;
        return;
      }
      throw error;
    }
  },

  async markAllRead(): Promise<void> {
    try {
      await apiClient<unknown>('/v1/notifications/read-all', { method: 'POST' });
    } catch (error) {
      if (isFallbackError(error)) {
        MOCK_NOTIFICATIONS.forEach((n) => {
          n.read = true;
        });
        return;
      }
      throw error;
    }
  },
};