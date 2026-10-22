'use client';

import { useEffect, useState } from 'react';
import { api, LeaderboardEntry, getXPProgress } from '@/lib/api';
import { UserCard } from '@/components/UserCard';

export default function LeaderboardPage() {
  const [entries, setEntries] = useState<LeaderboardEntry[]>([]);
  const [loading, setLoading] = useState(true);
  const [limit, setLimit] = useState(10);

  useEffect(() => {
    api.get(`/api/v1/leaderboard?limit=${limit}`)
      .then(res => setEntries(res.data.data || []))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [limit]);

  const medals = ['🥇', '🥈', '🥉'];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold text-white">🏆 Рейтинг</h1>
        <select
          value={limit}
          onChange={e => setLimit(Number(e.target.value))}
          className="bg-slate-800 border border-slate-700 text-white rounded-lg px-3 py-2"
        >
          <option value={10}>Топ 10</option>
          <option value={25}>Топ 25</option>
          <option value={50}>Топ 50</option>
        </select>
      </div>

      {loading ? (
        <div className="text-center text-slate-400 py-16">Загрузка...</div>
      ) : entries.length === 0 ? (
        <div className="card text-center text-slate-400 py-16">
          Нет данных. Отправьте первое событие через webhook.
        </div>
      ) : (
        <div className="space-y-3">
          {entries.map((entry, i) => (
            <div key={entry.user.id} className={`card flex items-center gap-4 ${i < 3 ? 'border-amber-500/30' : ''}`}>
              <div className="text-3xl w-12 text-center font-bold">
                {i < 3 ? medals[i] : <span className="text-slate-500">#{entry.rank}</span>}
              </div>
              <div className="flex-1">
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-bold">
                    {entry.user.name.charAt(0).toUpperCase()}
                  </div>
                  <div>
                    <div className="font-semibold text-white">{entry.user.name}</div>
                    <div className="text-sm text-slate-400">{entry.user.email}</div>
                  </div>
                </div>
              </div>
              <div className="text-right">
                <div className="xp-badge">⭐ {entry.xp.toLocaleString()} XP</div>
                <div className="level-badge mt-1">Ур. {entry.level}</div>
              </div>
              <div className="w-32">
                <div className="h-2 bg-slate-700 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-gradient-to-r from-blue-500 to-purple-500 rounded-full transition-all"
                    style={{ width: `${getXPProgress(entry.xp, entry.level)}%` }}
                  />
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
