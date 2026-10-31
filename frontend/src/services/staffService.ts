import { ApiError, apiClient } from './api';
import { mockAdminService } from './mockAdminData';
import type { AdminStaff } from './mockAdminData';

export type { AdminStaff } from './mockAdminData';

export interface InviteStaffInput {
  name: string;
  email: string;
  role: 'admin' | 'staff';
  department: string;
}

/** Whether an error means the backend is unavailable or the caller lacks a
 * usable identity — in which case we fall back to the local mock. */
function isFallbackError(error: unknown): boolean {
  if (error instanceof TypeError) return true;
  if (error instanceof ApiError && (error.status === 401 || error.status === 403)) return true;
  return false;
}

export const staffService = {
  async getStaffMembers(): Promise<AdminStaff[]> {
    try {
      return await apiClient<AdminStaff[]>('/v1/users');
    } catch (error) {
      if (isFallbackError(error)) return mockAdminService.getStaffMembers();
      throw error;
    }
  },

  async addStaffMember(input: InviteStaffInput): Promise<AdminStaff> {
    try {
      return await apiClient<AdminStaff>('/v1/users', {
        method: 'POST',
        body: JSON.stringify(input),
      });
    } catch (error) {
      if (isFallbackError(error)) {
        return mockAdminService.addStaffMember({
          ...input,
          status: 'invited',
        });
      }
      throw error;
    }
  },

  async updateStaffRole(id: string, role: 'admin' | 'staff'): Promise<AdminStaff> {
    try {
      return await apiClient<AdminStaff>(`/v1/users/${id}`, {
        method: 'PUT',
        body: JSON.stringify({ role }),
      });
    } catch (error) {
      if (isFallbackError(error)) return mockAdminService.updateStaffRole(id, role);
      throw error;
    }
  },
};