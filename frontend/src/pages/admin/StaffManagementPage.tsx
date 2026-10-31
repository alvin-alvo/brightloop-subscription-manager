import React, { useState, useEffect } from 'react';
import {
  UserPlus,
  Shield,
  Mail,
  X,
  UserCheck,
  Building2,
  Clock,
  Sparkles,
} from 'lucide-react';
import {
  staffService,
  type AdminStaff,
} from '../../services/staffService';
import { StatusBadge } from '../../components/StatusBadge';
import { useToast } from '../../context/ToastContext';
import { Button } from '../../components/common/Button';
import { InputField } from '../../components/common/InputField';
import { isValidEmail } from '../../utils/validations';

export const StaffManagementPage: React.FC = () => {
  const { addToast } = useToast();

  const [staffList, setStaffList] = useState<AdminStaff[]>([]);
  const [showModal, setShowModal] = useState(false);
  const [loadError, setLoadError] = useState<string | null>(null);

  useEffect(() => {
    let active = true;
    staffService
      .getStaffMembers()
      .then((members) => {
        if (active) setStaffList(members);
      })
      .catch((err: unknown) => {
        if (active) {
          setLoadError(
            err instanceof Error ? err.message : 'Unable to load staff members.'
          );
        }
      });
    return () => {
      active = false;
    };
  }, []);

  // Invite Form State
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    role: 'staff' as 'admin' | 'staff',
    department: 'Revenue & Operations',
  });

  const [errors, setErrors] = useState<Record<string, string>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  const handleInviteStaff = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrors({});

    const newErrors: Record<string, string> = {};
    if (!formData.name.trim()) {
      newErrors.name = 'Full name is required';
    }
    if (!formData.email.trim()) {
      newErrors.email = 'Email address is required';
    } else if (!isValidEmail(formData.email)) {
      newErrors.email = 'Please enter a valid work email address';
    }

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setIsSubmitting(true);

    try {
      const newMember = await staffService.addStaffMember({
        name: formData.name.trim(),
        email: formData.email.trim().toLowerCase(),
        role: formData.role,
        department: formData.department.trim() || 'General Operations',
      });

      const members = await staffService.getStaffMembers();
      setStaffList(members);
      setShowModal(false);

      setFormData({
        name: '',
        email: '',
        role: 'staff',
        department: 'Revenue & Operations',
      });

      addToast(
        'Staff Invitation Dispatched',
        `Invitation sent to ${newMember.email} with ${newMember.role.toUpperCase()} privileges.`,
        'success'
      );
    } catch (err: unknown) {
      addToast(
        'Invitation Failed',
        err instanceof Error ? err.message : 'Unable to send the invitation.',
        'error'
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleToggleRole = async (id: string, currentRole: 'admin' | 'staff') => {
    const nextRole = currentRole === 'admin' ? 'staff' : 'admin';
    try {
      const updated = await staffService.updateStaffRole(id, nextRole);
      const members = await staffService.getStaffMembers();
      setStaffList(members);
      addToast(
        'Role Updated',
        `${updated.name} permissions modified to ${nextRole.toUpperCase()}.`,
        'info'
      );
    } catch (err: unknown) {
      addToast(
        'Role Update Failed',
        err instanceof Error ? err.message : 'Unable to update the role.',
        'error'
      );
    }
  };

  return (
    <div className="space-y-6 animate-fadeIn">
      {loadError && (
        <div className="p-3.5 bg-rose-50 border border-rose-200 rounded-lg text-xs text-stripe-ruby">
          {loadError}
        </div>
      )}
      {/* Page Header */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-2xl font-light text-stripe-ink tracking-tight">
              Staff & Internal Access Control
            </h1>
            <span className="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-indigo-50 text-stripe-primary border border-indigo-200 uppercase">
              Admin Exclusive
            </span>
          </div>
          <p className="text-sm text-stripe-ink-mute mt-1">
            Grant permission scopes, invite internal staff members, and audit organizational RBAC seats.
          </p>
        </div>

        <Button
          type="button"
          variant="primary"
          size="md"
          onClick={() => setShowModal(true)}
          leftIcon={<UserPlus className="w-4 h-4" />}
        >
          Invite Staff Member
        </Button>
      </div>

      {/* RBAC Breakdown Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="p-4 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm">
          <div className="flex items-center justify-between text-stripe-ink-mute">
            <span className="text-xs font-normal uppercase tracking-wider">Total Internal Team</span>
            <UserCheck className="w-4 h-4 text-stripe-primary" />
          </div>
          <div className="text-2xl font-light text-stripe-ink mt-1 tnum">
            {staffList.length} Accounts
          </div>
          <span className="text-xs text-stripe-ink-mute">Active in Brightloop system</span>
        </div>

        <div className="p-4 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm">
          <div className="flex items-center justify-between text-stripe-ink-mute">
            <span className="text-xs font-normal uppercase tracking-wider">System Administrators</span>
            <Shield className="w-4 h-4 text-indigo-600" />
          </div>
          <div className="text-2xl font-light text-stripe-ink mt-1 tnum">
            {staffList.filter((s) => s.role === 'admin').length}
          </div>
          <span className="text-xs text-stripe-ink-mute">Unrestricted financial & staff control</span>
        </div>

        <div className="p-4 bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm">
          <div className="flex items-center justify-between text-stripe-ink-mute">
            <span className="text-xs font-normal uppercase tracking-wider">Operational Staff</span>
            <Sparkles className="w-4 h-4 text-emerald-600" />
          </div>
          <div className="text-2xl font-light text-stripe-ink mt-1 tnum">
            {staffList.filter((s) => s.role === 'staff').length}
          </div>
          <span className="text-xs text-stripe-ink-mute">CRM, quoting, and billing operations</span>
        </div>
      </div>

      {/* Staff Members Table */}
      <div className="bg-white border border-stripe-hairline rounded-xl shadow-stripe-sm overflow-hidden">
        <div className="p-4 border-b border-stripe-hairline flex items-center justify-between">
          <h2 className="text-sm font-medium text-stripe-ink">Provisioned Internal Members</h2>
          <span className="text-xs text-stripe-ink-mute font-mono">
            {staffList.length} Active Profiles
          </span>
        </div>

        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="border-b border-stripe-hairline bg-stripe-canvas-soft text-[11px] font-medium text-stripe-ink-mute uppercase tracking-wider">
                <th className="py-3 px-5">Member</th>
                <th className="py-3 px-5">Department</th>
                <th className="py-3 px-5">Access Role</th>
                <th className="py-3 px-5">Account Status</th>
                <th className="py-3 px-5 text-right">Last Active</th>
                <th className="py-3 px-5 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stripe-hairline text-sm">
              {staffList.map((staff) => (
                <tr key={staff.id} className="hover:bg-slate-50/70 transition-colors">
                  {/* Name and Email */}
                  <td className="py-3.5 px-5">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-full bg-indigo-50 border border-indigo-100 text-stripe-primary flex items-center justify-center font-medium text-xs shrink-0 shadow-sm">
                        {staff.name
                          .split(' ')
                          .map((n) => n[0])
                          .join('')
                          .slice(0, 2)
                          .toUpperCase()}
                      </div>
                      <div>
                        <div className="font-medium text-stripe-ink">{staff.name}</div>
                        <div className="text-xs text-stripe-ink-mute flex items-center gap-1">
                          <Mail className="w-3 h-3 text-slate-400" />
                          <span>{staff.email}</span>
                        </div>
                      </div>
                    </div>
                  </td>

                  {/* Department */}
                  <td className="py-3.5 px-5 text-xs text-stripe-ink-secondary">
                    <div className="flex items-center gap-1.5">
                      <Building2 className="w-3.5 h-3.5 text-slate-400" />
                      <span>{staff.department}</span>
                    </div>
                  </td>

                  {/* Role */}
                  <td className="py-3.5 px-5">
                    <span
                      className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-medium border ${
                        staff.role === 'admin'
                          ? 'bg-indigo-50 text-stripe-primary border-indigo-200'
                          : 'bg-slate-100 text-slate-800 border-slate-200'
                      }`}
                    >
                      <Shield className="w-3 h-3" />
                      {staff.role === 'admin' ? 'Administrator' : 'Staff'}
                    </span>
                  </td>

                  {/* Status */}
                  <td className="py-3.5 px-5">
                    <StatusBadge
                      label={
                        staff.status === 'active'
                          ? 'Active'
                          : staff.status === 'invited'
                          ? 'Invite Pending'
                          : 'Suspended'
                      }
                      variant={
                        staff.status === 'active'
                          ? 'success'
                          : staff.status === 'invited'
                          ? 'warning'
                          : 'neutral'
                      }
                    />
                  </td>

                  {/* Last Active */}
                  <td className="py-3.5 px-5 text-right text-xs text-stripe-ink-mute tnum">
                    <div className="flex items-center justify-end gap-1">
                      <Clock className="w-3 h-3 text-slate-400" />
                      <span>{staff.lastActive}</span>
                    </div>
                  </td>

                  {/* Role Modification Action */}
                  <td className="py-3.5 px-5 text-right">
                    <button
                      onClick={() => handleToggleRole(staff.id, staff.role)}
                      className="text-xs font-normal text-stripe-primary hover:text-stripe-primary-deep hover:underline"
                    >
                      Make {staff.role === 'admin' ? 'Staff' : 'Admin'}
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* --- INVITE STAFF MODAL --- */}
      {showModal && (
        <div className="fixed inset-0 z-50 bg-slate-900/40 backdrop-blur-sm flex items-center justify-center p-4">
          <div className="bg-white border border-stripe-hairline rounded-2xl shadow-stripe-lg max-w-md w-full p-6 animate-fadeIn">
            <div className="flex items-center justify-between pb-3 border-b border-stripe-hairline">
              <div className="flex items-center gap-2">
                <UserPlus className="w-5 h-5 text-stripe-primary" />
                <h2 className="text-lg font-light text-stripe-ink">Invite Staff Member</h2>
              </div>
              <button
                onClick={() => setShowModal(false)}
                className="p-1 text-slate-400 hover:text-stripe-ink rounded"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <form onSubmit={handleInviteStaff} className="space-y-4 pt-4">
              <InputField
                label="Full Name"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                error={errors.name}
                placeholder="e.g., Taylor Evans"
                required
                disabled={isSubmitting}
              />

              <InputField
                label="Work Email Address"
                type="email"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                error={errors.email}
                placeholder="taylor.evans@brightloop.io"
                required
                disabled={isSubmitting}
              />

              <div className="grid grid-cols-2 gap-3">
                <div>
                  <label className="block text-xs font-medium text-stripe-ink-secondary mb-1.5">
                    Assigned Role
                  </label>
                  <select
                    value={formData.role}
                    onChange={(e) =>
                      setFormData({ ...formData, role: e.target.value as 'admin' | 'staff' })
                    }
                    disabled={isSubmitting}
                    className="w-full bg-white text-stripe-ink text-sm rounded-md px-3 py-2 border border-stripe-hairline-input focus:outline-none focus:border-stripe-primary"
                  >
                    <option value="staff">Staff (Operations)</option>
                    <option value="admin">Admin (Full Control)</option>
                  </select>
                </div>

                <div>
                  <label className="block text-xs font-medium text-stripe-ink-secondary mb-1.5">
                    Department
                  </label>
                  <select
                    value={formData.department}
                    onChange={(e) => setFormData({ ...formData, department: e.target.value })}
                    disabled={isSubmitting}
                    className="w-full bg-white text-stripe-ink text-sm rounded-md px-3 py-2 border border-stripe-hairline-input focus:outline-none focus:border-stripe-primary"
                  >
                    <option value="Revenue & Operations">Revenue & Operations</option>
                    <option value="Customer Billing Support">Billing Support</option>
                    <option value="Enterprise Sales">Enterprise Sales</option>
                    <option value="Executive Leadership">Leadership</option>
                    <option value="Finance & Accounting">Finance & Tax</option>
                  </select>
                </div>
              </div>

              <div className="p-3 bg-slate-50 rounded-lg border border-slate-200/60 text-xs text-stripe-ink-mute">
                <span className="font-semibold text-stripe-ink-secondary block mb-1">
                  RBAC Scope Information:
                </span>
                {formData.role === 'admin'
                  ? 'Administrators have complete visibility into system-wide financials, staff privileges, product pricing, and discount rules.'
                  : 'Staff members manage operational customer accounts, draft quotations, review subscriptions, and execute billing cycles.'}
              </div>

              <div className="pt-2 flex justify-end gap-2">
                <Button
                  type="button"
                  variant="ghost"
                  size="md"
                  onClick={() => setShowModal(false)}
                  disabled={isSubmitting}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  variant="primary"
                  size="md"
                  isLoading={isSubmitting}
                  loadingText="Dispatching invite..."
                >
                  Send Invitation
                </Button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default StaffManagementPage;
