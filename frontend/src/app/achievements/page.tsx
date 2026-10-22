'use client';

import { useEffect, useState } from 'react';
import { api, Achievement } from '@/lib/api';

export default function AchievementsPage() {
  const [achievements, setAchievements] = useState<Achievement[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get('/api/v1/achievements')
      .then(res => setAchievements(res.data.data || []))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  const conditionLabels: Record<string, string> = {
    total_xp: 'Набрать XP',
    test_count: 'Пройти тестов',
    video_count: 'Посмотреть видео',
    modules_per_week: 'Модулей за неделю',
    test_perfect: 'Идеальных тестов',
    streak_days: 'Дней подряд',
  };

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-white">🎖️ Достижения</h1>

      {loading ? (
        <div className="text-center text-slate-400 py-16">Загрузка...</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {achievements.map(ach => (
            <div key={ach.id} className="card hover:border-amber-500/50 transition-colors">
              <div className="text-4xl mb-3">{ach.icon_url || '🏅'}</div>
              <h3 className="font-bold text-white text-lg">{ach.name}</h3>
              <p className="text-slate-400 text-sm mt-1">{ach.description}</p>
              <div className="mt-3 flex items-center justify-between">
                <span className="text-xs text-slate-500">
                  {conditionLabels[ach.condition_type] || ach.condition_type}: {ach.condition_value}
                </span>
                {ach.xp_reward > 0 && (
                  <span className="xp-badge text-xs">+{ach.xp_reward} XP</span>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
