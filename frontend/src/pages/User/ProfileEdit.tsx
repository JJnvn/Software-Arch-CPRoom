import { useEffect, useState } from 'react';
import * as auth from '@/services/auth';

export default function ProfileEdit() {
  const [profile, setProfile] = useState<any>(null);
  const [editing, setEditing] = useState(false);
  const [formData, setFormData] = useState({ name: '', email: '', password: '' });
  const [status, setStatus] = useState('');

  useEffect(() => {
    loadProfile();
  }, []);

  const loadProfile = async () => {
    try {
      const data = await auth.getProfile();
      setProfile(data);
      setFormData({ name: data.name || '', email: data.email || '', password: '' });
    } catch {
      setStatus('Failed to load profile');
    }
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    setStatus('Updating...');
    try {
      const payload: any = {};
      if (formData.name !== profile.name) payload.name = formData.name;
      if (formData.email !== profile.email) payload.email = formData.email;
      if (formData.password) payload.password = formData.password;

      await auth.updateProfile(payload);
      await loadProfile();
      setFormData(prev => ({ ...prev, password: '' }));
      setEditing(false);
      setStatus('Profile updated successfully!');
      setTimeout(() => setStatus(''), 3000);
    } catch (error: any) {
      setStatus(error.response?.data?.error || 'Failed to update profile');
    }
  };

  if (!profile) return <div className="page">Loading...</div>;

  return (
    <div className="page">
      <h1 className="page-title">My Profile</h1>
      
      {!editing ? (
        <div className="card space-y-3">
          <div className="space-y-2">
            <div><strong>Name:</strong> {profile.name || '—'}</div>
            <div><strong>Email:</strong> {profile.email || '—'}</div>
            <div><strong>Role:</strong> {profile.role || 'user'}</div>
            <div><strong>ID:</strong> <span className="text-xs text-gray-500">{profile.id}</span></div>
          </div>
          <button onClick={() => setEditing(true)} className="btn-primary">
            Edit Profile
          </button>
        </div>
      ) : (
        <form onSubmit={handleUpdate} className="card space-y-3">
          <div>
            <label className="block text-sm font-medium mb-1">Name</label>
            <input
              type="text"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="input"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">Email</label>
            <input
              type="email"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              className="input"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-1">New Password</label>
            <input
              type="password"
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              className="input"
              placeholder="Leave blank to keep current password"
            />
            <p className="text-xs text-gray-500 mt-1">Only fill this if you want to change your password</p>
          </div>
          {status && (
            <div className={`text-sm p-2 rounded ${status.includes('success') ? 'bg-green-50 text-green-700' : status.includes('Updating') ? 'bg-blue-50 text-blue-700' : 'bg-red-50 text-red-700'}`}>
              {status}
            </div>
          )}
          <div className="flex gap-2">
            <button type="submit" className="btn-primary">Save Changes</button>
            <button type="button" onClick={() => {
              setEditing(false);
              setFormData({ name: profile.name || '', email: profile.email || '', password: '' });
              setStatus('');
            }} className="btn-secondary">
              Cancel
            </button>
          </div>
        </form>
      )}
    </div>
  );
}
