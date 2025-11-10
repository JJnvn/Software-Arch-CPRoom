import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';

export default function OAuthCallback() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const token = searchParams.get('token');
    const errorParam = searchParams.get('error');

    if (errorParam) {
      setError('Authentication failed. Redirecting to login...');
      setTimeout(() => navigate('/login?error=' + errorParam), 2000);
      return;
    }

    if (token) {
      // Store token
      localStorage.setItem('AUTH_TOKEN', token);
      
      // Reload to trigger auth context update
      window.location.href = '/';
    } else {
      setError('No authentication token received. Redirecting to login...');
      setTimeout(() => navigate('/login'), 2000);
    }
  }, [searchParams, navigate]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 p-6">
      <div className="bg-white rounded-lg shadow-lg p-8 w-full max-w-md text-center">
        {error ? (
          <>
            <svg className="w-16 h-16 mx-auto text-red-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <h1 className="text-2xl font-bold text-gray-900 mb-2">Authentication Failed</h1>
            <p className="text-gray-600">{error}</p>
          </>
        ) : (
          <>
            <div className="spinner mx-auto mb-4"></div>
            <h1 className="text-2xl font-bold text-gray-900 mb-2">Completing Sign In</h1>
            <p className="text-gray-600">Please wait while we set up your account...</p>
          </>
        )}
      </div>
    </div>
  );
}
