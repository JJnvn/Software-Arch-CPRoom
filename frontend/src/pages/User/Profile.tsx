import { FormEvent, useEffect, useState } from "react";
import { useAuth } from "@/hooks/useAuth";
import * as auth from "@/services/auth";

export default function Profile() {
    const { user } = useAuth();
    const [name, setName] = useState(user?.name || "");
    const [email, setEmail] = useState(user?.email || "");
    const [status, setStatus] = useState<string | null>(null);

    useEffect(() => {
        setName(user?.name || "");
        setEmail(user?.email || "");
    }, [user]);

    async function onSubmit(e: FormEvent) {
        e.preventDefault();
        const raw = localStorage.getItem("auth_user");
        if (!raw) {
            console.error("no auth_user in local storage");
            return;
        }
        const user = JSON.parse(raw);
        const id = user.id;
        await auth.updateProfile(id, { name, email });
        setStatus("Profile updated");
        setTimeout(() => setStatus(null), 2000);
    }

    return (
        <div className="page">
            <h1 className="page-title">Profile</h1>
            <form onSubmit={onSubmit} className="card max-w-xl space-y-4">
                {status && (
                    <div className="text-green-700 text-sm">{status}</div>
                )}
                <div>
                    <label className="block text-sm mb-1">Name</label>
                    <input
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        className="w-full border rounded px-3 py-2"
                    />
                </div>
                <div>
                    <label className="block text-sm mb-1">Email</label>
                    <input
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        type="email"
                        className="w-full border rounded px-3 py-2"
                    />
                </div>
                <button className="px-4 py-2 bg-blue-600 text-white rounded">
                    Save Changes
                </button>
            </form>
        </div>
    );
}
