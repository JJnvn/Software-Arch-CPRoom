import { FormEvent, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

export default function Register() {
  const { register, loading } = useAuth();
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    try {
      await register(name, email, password);
      navigate('/', { replace: true });
    } catch (e: any) {
      setError('Registration failed');
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 p-6">
      <form onSubmit={onSubmit} className="bg-white rounded-lg shadow p-6 w-full max-w-md space-y-4">
        <h1 className="text-2xl font-semibold">Register</h1>
        {error && <div className="text-red-600 text-sm">{error}</div>}
        <div>
          <label className="block text-sm mb-1">Name</label>
          <input value={name} onChange={e => setName(e.target.value)} className="w-full border rounded px-3 py-2" placeholder="Jane Doe" required />
        </div>
        <div>
          <label className="block text-sm mb-1">Email</label>
          <input value={email} onChange={e => setEmail(e.target.value)} type="email" className="w-full border rounded px-3 py-2" placeholder="you@company.com" required />
        </div>
        <div>
          <label className="block text-sm mb-1">Password</label>
          <input value={password} onChange={e => setPassword(e.target.value)} type="password" className="w-full border rounded px-3 py-2" placeholder="Create a strong password" required />
        </div>
        <button disabled={loading} className="w-full bg-blue-600 text-white rounded py-2">{loading ? 'Creating account…' : 'Create Account'}</button>
        <p className="text-sm text-gray-600">Have an account? <Link to="/login" className="text-blue-600 hover:underline">Login</Link></p>
      </form>
    </div>
  );
}

