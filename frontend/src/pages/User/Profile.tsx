import { FormEvent, useEffect, useState } from "react";
import { useAuth } from "@/hooks/useAuth";
import * as auth from "@/services/auth";

export default function Profile() {
    const { user, refreshUser } = useAuth();
    const [name, setName] = useState(user?.name || "");
    const [email, setEmail] = useState(user?.email || "");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [successMessage, setSuccessMessage] = useState<string | null>(null);
    const [errorMessage, setErrorMessage] = useState<string | null>(null);

    useEffect(() => {
        setName(user?.name || "");
        setEmail(user?.email || "");
    }, [user]);

    async function onSubmit(e: FormEvent) {
        e.preventDefault();
        setErrorMessage(null);
        setSuccessMessage(null);

        // Validation
        if (!name.trim()) {
            setErrorMessage("Name is required");
            return;
        }

        if (!email.trim() || !/^\S+@\S+\.\S+$/.test(email)) {
            setErrorMessage("Valid email is required");
            return;
        }

        if (password && password.length < 6) {
            setErrorMessage("Password must be at least 6 characters");
            return;
        }

        if (password !== confirmPassword) {
            setErrorMessage("Passwords do not match");
            return;
        }

        setIsSubmitting(true);

        try {
            const payload: any = { name, email };
            if (password) {
                payload.password = password;
            }

            await auth.updateProfile(payload);
            
            // Refresh user context to reflect changes
            if (refreshUser) {
                await refreshUser();
            }

            setSuccessMessage("Profile updated successfully!");
            setPassword("");
            setConfirmPassword("");
            
            setTimeout(() => setSuccessMessage(null), 3000);
        } catch (error: any) {
            const errMsg = error.response?.data?.error || "Failed to update profile";
            setErrorMessage(errMsg);
        } finally {
            setIsSubmitting(false);
        }
    }

    return (
        <div className="page">
            <div className="max-w-2xl mx-auto">
                <h1 className="page-title">My Profile</h1>
                <p className="text-gray-600 mb-6">Update your account information and password</p>

                <form onSubmit={onSubmit} className="card space-y-6">
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

                    {/* Account Information Section */}
                    <div>
                        <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
                            <svg className="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                            </svg>
                            Account Information
                        </h2>
                        
                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Name *
                                </label>
                                <input
                                    type="text"
                                    value={name}
                                    onChange={(e) => {
                                        setName(e.target.value);
                                        setErrorMessage(null);
                                    }}
                                    className="w-full border border-gray-300 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                    placeholder="Enter your name"
                                    required
                                    disabled={isSubmitting}
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Email *
                                </label>
                                <input
                                    type="email"
                                    value={email}
                                    onChange={(e) => {
                                        setEmail(e.target.value);
                                        setErrorMessage(null);
                                    }}
                                    className="w-full border border-gray-300 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                    placeholder="your.email@example.com"
                                    required
                                    disabled={isSubmitting}
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Role
                                </label>
                                <div className="flex items-center gap-2">
                                    <span className={`badge-${user?.role === 'staff' ? 'warning' : 'info'}`}>
                                        {user?.role || 'user'}
                                    </span>
                                    <span className="text-sm text-gray-500">(Cannot be changed)</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    {/* Password Change Section */}
                    <div className="border-t pt-6">
                        <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
                            <svg className="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                            </svg>
                            Change Password
                        </h2>
                        <p className="text-sm text-gray-600 mb-4">Leave blank if you don't want to change your password</p>

                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    New Password
                                </label>
                                <input
                                    type="password"
                                    value={password}
                                    onChange={(e) => {
                                        setPassword(e.target.value);
                                        setErrorMessage(null);
                                    }}
                                    className="w-full border border-gray-300 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                    placeholder="Enter new password (min 6 characters)"
                                    disabled={isSubmitting}
                                />
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">
                                    Confirm New Password
                                </label>
                                <input
                                    type="password"
                                    value={confirmPassword}
                                    onChange={(e) => {
                                        setConfirmPassword(e.target.value);
                                        setErrorMessage(null);
                                    }}
                                    className="w-full border border-gray-300 rounded-lg px-4 py-2.5 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                    placeholder="Confirm new password"
                                    disabled={isSubmitting}
                                />
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
                                    <span>Save Changes</span>
                                </>
                            )}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}
