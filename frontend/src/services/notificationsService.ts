import { apiClient } from './api';

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

/** In-app notifications against the real notification service. */
export const notificationsService = {
  async getNotifications(limit = 20): Promise<NotificationList> {
    return apiClient<NotificationList>(`/v1/notifications?limit=${limit}`);
  },

  async markRead(id: string): Promise<void> {
    await apiClient<unknown>(`/v1/notifications/${id}/read`, { method: 'POST' });
  },

  async markAllRead(): Promise<void> {
    await apiClient<unknown>('/v1/notifications/read-all', { method: 'POST' });
  },
};