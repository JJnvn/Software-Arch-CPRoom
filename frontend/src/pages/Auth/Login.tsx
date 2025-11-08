import { FormEvent, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

export default function Login() {
  const { login, loading } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();
  const location = useLocation() as any;

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    try {
      await login(email, password);
      const to = location.state?.from?.pathname || '/';
      navigate(to, { replace: true });
    } catch (e: any) {
      setError('Login failed');
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 p-6">
      <form onSubmit={onSubmit} className="bg-white rounded-lg shadow p-6 w-full max-w-md space-y-4">
        <h1 className="text-2xl font-semibold">Login</h1>
        {error && <div className="text-red-600 text-sm">{error}</div>}
        <div>
          <label className="block text-sm mb-1">Email</label>
          <input value={email} onChange={e => setEmail(e.target.value)} type="email" className="w-full border rounded px-3 py-2" placeholder="you@company.com" required />
        </div>
        <div>
          <label className="block text-sm mb-1">Password</label>
          <input value={password} onChange={e => setPassword(e.target.value)} type="password" className="w-full border rounded px-3 py-2" placeholder="••••••••" required />
        </div>
        <button disabled={loading} className="w-full bg-blue-600 text-white rounded py-2">{loading ? 'Signing in…' : 'Sign In'}</button>
        <p className="text-sm text-gray-600">No account? <Link to="/register" className="text-blue-600 hover:underline">Register</Link></p>
      </form>
    </div>
  );
}

