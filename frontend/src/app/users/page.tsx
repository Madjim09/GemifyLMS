'use client';

import { useEffect, useState } from 'react';
import { api, User, getXPProgress } from '@/lib/api';
import Link from 'next/link';

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get('/api/v1/users?limit=20')
      .then(res => {
        setUsers(res.data.data || []);
        setTotal(res.data.total || 0);
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-white">👥 Пользователи ({total})</h1>

      {loading ? (
        <div className="text-center text-slate-400 py-16">Загрузка...</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {users.map(user => (
            <Link key={user.id} href={`/users/${user.id}`} className="card hover:border-blue-500 transition-colors">
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white text-xl font-bold">
                  {user.name.charAt(0).toUpperCase()}
                </div>
                <div className="flex-1 min-w-0">
                  <div className="font-semibold text-white truncate">{user.name}</div>
                  <div className="text-sm text-slate-400 truncate">{user.email}</div>
                  <div className="flex gap-2 mt-1">
                    <span className="xp-badge text-xs">⭐ {user.xp.toLocaleString()} XP</span>
                    <span className="level-badge text-xs">Ур. {user.level}</span>
                  </div>
                </div>
              </div>
              <div className="mt-3">
                <div className="h-1.5 bg-slate-700 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-gradient-to-r from-blue-500 to-purple-500 rounded-full"
                    style={{ width: `${getXPProgress(user.xp, user.level)}%` }}
                  />
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
