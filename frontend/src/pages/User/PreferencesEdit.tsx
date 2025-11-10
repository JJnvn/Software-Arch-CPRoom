import { useState, useEffect, FormEvent } from 'react';
import { useAuth } from '@/hooks/useAuth';
import * as auth from '@/services/auth';

export default function PreferencesEdit() {
  const { user } = useAuth();
  const [notifType, setNotifType] = useState('email');
  const [language, setLanguage] = useState('en');
  const [enabledChannels, setEnabledChannels] = useState<string[]>(['email']);
  const [isLoading, setIsLoading] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  // Notification event preferences
  const [notifyConfirmation, setNotifyConfirmation] = useState(true);
  const [notifyApproval, setNotifyApproval] = useState(true);
  const [notifyReminders, setNotifyReminders] = useState(true);
  const [notifyCancellation, setNotifyCancellation] = useState(true);

  useEffect(() => {
    (async () => {
      if (!user) {
        setIsLoading(false);
        return;
      }
      try {
        const prefs = await auth.getPreferences();
        
        // Extract preferences from nested structure
        const preferences = prefs.preferences || {};
        if (preferences.notification_type) {
          setNotifType(preferences.notification_type);
        }
        if (preferences.language) {
          setLanguage(preferences.language);
        }
        
        // Set enabled channels
        if (prefs.enabled_channels && prefs.enabled_channels.length > 0) {
          setEnabledChannels(prefs.enabled_channels);
        }
        
        // Set event preferences if available
        if (preferences.notify_confirmation !== undefined) {
          setNotifyConfirmation(preferences.notify_confirmation);
        }
        if (preferences.notify_approval !== undefined) {
          setNotifyApproval(preferences.notify_approval);
        }
        if (preferences.notify_reminders !== undefined) {
          setNotifyReminders(preferences.notify_reminders);
        }
        if (preferences.notify_cancellation !== undefined) {
          setNotifyCancellation(preferences.notify_cancellation);
        }
      } catch (error) {
        console.error('Failed to load preferences:', error);
        // Use defaults on error
      } finally {
        setIsLoading(false);
      }
    })();
  }, [user]);

  const toggleChannel = (channel: string) => {
    setEnabledChannels(prev => {
      if (prev.includes(channel)) {
        // Don't allow removing all channels
        if (prev.length === 1) return prev;
        return prev.filter(ch => ch !== channel);
      }
      return [...prev, channel];
    });
  };

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setErrorMessage(null);
    setSuccessMessage(null);

    if (!user) {
      setErrorMessage('User not logged in');
      return;
    }

    if (enabledChannels.length === 0) {
      setErrorMessage('At least one notification channel must be enabled');
      return;
    }

    setIsSubmitting(true);

    try {
      await auth.updatePreferences({
        notification_type: notifType,
        language: language,
        enabled_channels: enabledChannels,
      });
      
      setSuccessMessage('Preferences saved successfully!');
      setTimeout(() => setSuccessMessage(null), 3000);
    } catch (error: any) {
      const errMsg = error.response?.data?.error || 'Failed to save preferences';
      setErrorMessage(errMsg);
    } finally {
      setIsSubmitting(false);
    }
  }

  if (isLoading) {
    return (
      <div className="page">
        <div className="flex items-center justify-center py-12">
          <div className="spinner"></div>
          <span className="ml-3 text-gray-600">Loading preferences...</span>
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="page">
        <div className="card text-center py-8">
          <svg className="w-16 h-16 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
          <p className="text-gray-600">Please log in to manage your preferences.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="page">
      <div className="max-w-3xl mx-auto">
        <h1 className="page-title">Notification Preferences</h1>
        <p className="text-gray-600 mb-6">Manage how and when you receive notifications about your bookings</p>

        <form onSubmit={handleSubmit} className="card space-y-6">
          {/* Success Message */}
          {successMessage && (
            <div className="alert-success">
              <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
              <span>{successMessage}</span>
            </div>
          )}

          {/* Error Message */}
          {errorMessage && (
            <div className="alert-error">
              <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
              </svg>
              <span>{errorMessage}</span>
            </div>
          )}

          {/* Notification Channels */}
          <div>
            <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
              <svg className="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
              </svg>
              Delivery Channels
            </h2>
            <p className="text-sm text-gray-600 mb-3">Select how you want to receive notifications</p>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
              <div
                onClick={() => toggleChannel('email')}
                className={`p-4 border-2 rounded-lg cursor-pointer transition-all ${
                  enabledChannels.includes('email')
                    ? 'border-blue-500 bg-blue-50'
                    : 'border-gray-200 bg-white hover:border-gray-300'
                }`}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                    </svg>
                    <div>
                      <div className="font-medium">Email</div>
                      <div className="text-xs text-gray-500">Receive detailed notifications via email</div>
                    </div>
                  </div>
                  <input
                    type="checkbox"
                    checked={enabledChannels.includes('email')}
                    onChange={() => {}}
                    className="w-5 h-5"
                  />
                </div>
              </div>

              <div
                onClick={() => toggleChannel('push')}
                className={`p-4 border-2 rounded-lg cursor-pointer transition-all ${
                  enabledChannels.includes('push')
                    ? 'border-blue-500 bg-blue-50'
                    : 'border-gray-200 bg-white hover:border-gray-300'
                }`}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
                    </svg>
                    <div>
                      <div className="font-medium">Push Notifications</div>
                      <div className="text-xs text-gray-500">Get instant alerts on your device</div>
                    </div>
                  </div>
                  <input
                    type="checkbox"
                    checked={enabledChannels.includes('push')}
                    onChange={() => {}}
                    className="w-5 h-5"
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Primary Notification Type */}
          <div>
            <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
              <svg className="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              General Settings
            </h2>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Default Notification Type
                </label>
                <select
                  value={notifType}
                  onChange={(e) => {
                    setNotifType(e.target.value);
                    setErrorMessage(null);
                  }}
                  className="w-full border border-gray-300 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  disabled={isSubmitting}
                >
                  <option value="email">Email</option>
                  <option value="push">Push Notification</option>
                  <option value="all">All Channels</option>
                </select>
                <p className="text-xs text-gray-500 mt-1">Choose your preferred notification method</p>
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Language Preference
                </label>
                <select
                  value={language}
                  onChange={(e) => {
                    setLanguage(e.target.value);
                    setErrorMessage(null);
                  }}
                  className="w-full border border-gray-300 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  disabled={isSubmitting}
                >
                  <option value="en">🇺🇸 English</option>
                  <option value="th">🇹🇭 ไทย (Thai)</option>
                </select>
                <p className="text-xs text-gray-500 mt-1">Select your preferred language for notifications</p>
              </div>
            </div>
          </div>

          {/* Notification Events - Info Only */}
          <div className="border-t pt-6">
            <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
              <svg className="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Notification Events
            </h2>
            <p className="text-sm text-gray-600 mb-4">You'll receive notifications for these important events:</p>
            
            <div className="space-y-3 bg-gray-50 p-4 rounded-lg">
              <div className="flex items-center gap-3">
                <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                <span className="text-sm">Booking confirmation and creation</span>
              </div>
              <div className="flex items-center gap-3">
                <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                <span className="text-sm">Booking approval or denial (for staff-approval bookings)</span>
              </div>
              <div className="flex items-center gap-3">
                <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                <span className="text-sm">Booking reminders (30 minutes before start time)</span>
              </div>
              <div className="flex items-center gap-3">
                <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                <span className="text-sm">Booking cancellation or transfer notifications</span>
              </div>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-3 pt-4">
            <button
              type="submit"
              disabled={isSubmitting}
              className="btn-primary flex-1"
            >
              {isSubmitting ? (
                <>
                  <span className="spinner"></span>
                  <span>Saving...</span>
                </>
              ) : (
                <>
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                  <span>Save Preferences</span>
                </>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
