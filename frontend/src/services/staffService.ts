import { apiClient } from './api';

export interface AdminStaff {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'staff';
  department: string;
  status: 'active' | 'invited' | 'suspended';
  lastActive: string;
  createdAt: string;
}

export interface InviteStaffInput {
  name: string;
  email: string;
  role: 'admin' | 'staff';
  department: string;
}

/** Staff administration against the real user service. */
export const staffService = {
  async getStaffMembers(): Promise<AdminStaff[]> {
    return apiClient<AdminStaff[]>('/v1/users');
  },

  async addStaffMember(input: InviteStaffInput): Promise<AdminStaff> {
    return apiClient<AdminStaff>('/v1/users', {
      method: 'POST',
      body: JSON.stringify(input),
    });
  },

  async updateStaffRole(id: string, role: 'admin' | 'staff'): Promise<AdminStaff> {
    return apiClient<AdminStaff>(`/v1/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ role }),
    });
  },
};