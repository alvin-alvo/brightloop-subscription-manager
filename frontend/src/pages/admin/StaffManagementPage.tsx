import React from 'react';
import { UserPlus, Shield, Mail } from 'lucide-react';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';

export const StaffManagementPage: React.FC = () => {
  const { addToast } = useToast();

  const staffMembers = [
    {
      id: 'stf_01',
      name: 'Alex Vance',
      email: 'alex.admin@brightloop.io',
      role: 'Admin',
      department: 'Executive Leadership',
      status: 'active' as const,
      lastActive: 'Just now',
    },
    {
      id: 'stf_02',
      name: 'Sarah Miller',
      email: 'sarah.miller@brightloop.io',
      role: 'Staff',
      department: 'Revenue & Operations',
      status: 'active' as const,
      lastActive: '12m ago',
    },
    {
      id: 'stf_03',
      name: 'Marcus Brody',
      email: 'marcus.brody@brightloop.io',
      role: 'Staff',
      department: 'Customer Billing Support',
      status: 'active' as const,
      lastActive: '2h ago',
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
            Staff & Internal Access Control
          </h1>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Grant permission scopes, invite team members, and audit access logs.
          </p>
        </div>
        <button
          onClick={() =>
            addToast('Invite Sent', 'Invitation email dispatched to new team member.', 'success')
          }
          className="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium text-white bg-stripe-primary hover:bg-stripe-primary-deep rounded-full shadow-stripe-sm transition-colors"
        >
          <UserPlus className="w-4 h-4" />
          Invite Staff Member
        </button>
      </div>

      <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-stripe-hairline bg-stripe-canvas-soft text-[11px] font-medium text-stripe-ink-mute uppercase tracking-wider">
                <th className="py-3 px-5">Member</th>
                <th className="py-3 px-5">Department</th>
                <th className="py-3 px-5">Access Role</th>
                <th className="py-3 px-5">Status</th>
                <th className="py-3 px-5 text-right">Last Active</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stripe-hairline text-sm">
              {staffMembers.map((staff) => (
                <tr key={staff.id} className="hover:bg-slate-50/70 transition-colors">
                  <td className="py-3 px-5">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-full bg-indigo-50 border border-indigo-100 text-stripe-primary flex items-center justify-center font-medium text-xs">
                        {staff.name
                          .split(' ')
                          .map((n) => n[0])
                          .join('')}
                      </div>
                      <div>
                        <div className="font-medium text-stripe-ink">{staff.name}</div>
                        <div className="text-xs text-stripe-ink-mute flex items-center gap-1">
                          <Mail className="w-3 h-3" />
                          {staff.email}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="py-3 px-5 text-stripe-ink-secondary text-xs">
                    {staff.department}
                  </td>
                  <td className="py-3 px-5">
                    <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium bg-slate-100 text-stripe-ink-secondary">
                      <Shield className="w-3 h-3 text-stripe-primary" />
                      {staff.role}
                    </span>
                  </td>
                  <td className="py-3 px-5">
                    <StatusBadge label="Active" variant="success" />
                  </td>
                  <td className="py-3 px-5 text-right text-xs text-stripe-ink-mute tnum">
                    {staff.lastActive}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default StaffManagementPage;

