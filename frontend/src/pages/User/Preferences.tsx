import { FormEvent, useEffect, useState } from 'react';
import * as auth from '@/services/auth';

export default function Preferences() {
  const [notificationType, setNotificationType] = useState('email');
  const [language, setLanguage] = useState('en');
  const [status, setStatus] = useState<string | null>(null);

  useEffect(() => {
    // Could fetch preferences here if available
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    await auth.updatePreferences({ notificationType, language });
    setStatus('Preferences saved');
    setTimeout(() => setStatus(null), 2000);
  }

  return (
    <div className="page">
      <h1 className="page-title">Preferences</h1>
      <form onSubmit={onSubmit} className="card max-w-xl space-y-4">
        {status && <div className="text-green-700 text-sm">{status}</div>}
        <div>
          <label className="block text-sm mb-1">Notification Type</label>
          <select value={notificationType} onChange={e => setNotificationType(e.target.value)} className="w-full border rounded px-3 py-2">
            <option value="email">Email</option>
            <option value="sms">SMS</option>
            <option value="push">Push</option>
          </select>
        </div>
        <div>
          <label className="block text-sm mb-1">Language</label>
          <select value={language} onChange={e => setLanguage(e.target.value)} className="w-full border rounded px-3 py-2">
            <option value="en">English</option>
            <option value="es">Spanish</option>
            <option value="fr">French</option>
          </select>
        </div>
        <button className="px-4 py-2 bg-blue-600 text-white rounded">Save Preferences</button>
      </form>
    </div>
  );
}

