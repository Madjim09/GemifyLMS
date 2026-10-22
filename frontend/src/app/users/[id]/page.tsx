'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import {
  api, User, UserAchievement, UserSkill, UserQuizResult,
  getXPProgress, getXPForNextLevel, DIFFICULTY_LABELS, DIFFICULTY_COLORS
} from '@/lib/api';

const STATUS_LABELS = {
  locked: '🔒 Закрыт',
  available: '✨ Доступен',
  in_progress: '📖 В процессе',
  mastered: '✅ Освоен',
};

const STATUS_COLORS = {
  locked: 'bg-slate-700 text-slate-400',
  available: 'bg-blue-900 text-blue-400',
  in_progress: 'bg-yellow-900 text-yellow-400',
  mastered: 'bg-green-900 text-green-400',
};

export default function UserProfilePage() {
  const params = useParams();
  const id = params.id as string;
  const [user, setUser] = useState<User | null>(null);
  const [achievements, setAchievements] = useState<UserAchievement[]>([]);
  const [skills, setSkills] = useState<UserSkill[]>([]);
  const [quizResults, setQuizResults] = useState<UserQuizResult[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<'achievements' | 'skills' | 'tests'>('achievements');

  useEffect(() => {
    Promise.all([
      api.get(`/api/v1/users/${id}`),
      api.get(`/api/v1/users/${id}/achievements`),
      api.get(`/api/v1/users/${id}/skills`),
      api.get(`/api/v1/users/${id}/quiz-results`),
    ]).then(([userRes, achRes, skillRes, quizRes]) => {
      setUser(userRes.data.data);
      setAchievements(achRes.data.data || []);
      setSkills(skillRes.data.data || []);
      setQuizResults(quizRes.data.data || []);
    }).catch(console.error)
      .finally(() => setLoading(false));
  }, [id]);

  if (loading) return <div className="text-center text-slate-400 py-16">Загрузка профиля...</div>;
  if (!user) return <div className="text-center text-slate-400 py-16">Пользователь не найден</div>;

  const progress = getXPProgress(user.xp, user.level);
  const nextLevelXP = getXPForNextLevel(user.level);

  const earnedSkills = skills.filter(s => s.status === 'mastered' || s.status === 'in_progress');
  const availableSkills = skills.filter(s => s.status === 'available');
  const lockedSkills = skills.filter(s => s.status === 'locked');
  const bestScore = quizResults.length > 0 ? Math.max(...quizResults.map(r => r.score)) : 0;
  const perfectCount = quizResults.filter(r => r.is_perfect).length;

  return (
    <div className="space-y-6">
      {/* Profile Header */}
      <div className="card">
        <div className="flex items-start gap-5">
          <div className="w-20 h-20 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white text-3xl font-bold shrink-0">
            {user.name.charAt(0).toUpperCase()}
          </div>
          <div className="flex-1">
            <h1 className="text-2xl font-bold text-white">{user.name}</h1>
            <p className="text-slate-400 text-sm">{user.email}</p>
            <div className="flex flex-wrap gap-2 mt-2">
              <span className="xp-badge">⭐ {user.xp.toLocaleString()} XP</span>
              <span className="level-badge">Уровень {user.level}</span>
              <span className="bg-slate-700 text-slate-300 text-xs px-2 py-0.5 rounded-full">
                📝 {user.tests_passed} тестов
              </span>
              <span className="bg-slate-700 text-slate-300 text-xs px-2 py-0.5 rounded-full">
                🎖️ {achievements.length} достижений
              </span>
            </div>
          </div>
        </div>

        {/* XP Progress bar */}
        <div className="mt-4">
          <div className="flex justify-between text-sm text-slate-400 mb-1">
            <span>Прогресс до уровня {user.level + 1}</span>
            <span>
              {nextLevelXP > 0
                ? `${user.xp.toLocaleString()} / ${nextLevelXP.toLocaleString()} XP`
                : 'Максимальный уровень'}
            </span>
          </div>
          <div className="h-3 bg-slate-700 rounded-full overflow-hidden">
            <div
              className="h-full bg-gradient-to-r from-blue-500 to-purple-500 rounded-full transition-all"
              style={{ width: `${progress}%` }}
            />
          </div>
          <div className="text-right text-xs text-slate-500 mt-1">{progress.toFixed(1)}%</div>
        </div>
      </div>

      {/* Stats row */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div className="card text-center">
          <div className="text-2xl font-bold text-white">{achievements.length}</div>
          <div className="text-slate-400 text-xs">Достижений</div>
        </div>
        <div className="card text-center">
          <div className="text-2xl font-bold text-green-400">{earnedSkills.length}</div>
          <div className="text-slate-400 text-xs">Навыков освоено</div>
        </div>
        <div className="card text-center">
          <div className="text-2xl font-bold text-blue-400">{user.tests_passed}</div>
          <div className="text-slate-400 text-xs">Тестов пройдено</div>
        </div>
        <div className="card text-center">
          <div className="text-2xl font-bold text-purple-400">{perfectCount}</div>
          <div className="text-slate-400 text-xs">Идеальных тестов</div>
        </div>
      </div>

      {/* Tabs */}
      <div className="flex gap-1 bg-slate-800 rounded-xl p-1">
        {(['achievements', 'skills', 'tests'] as const).map(tab => (
          <button
            key={tab}
            onClick={() => setActiveTab(tab)}
            className={`flex-1 py-2 text-sm font-medium rounded-lg transition-colors ${
              activeTab === tab
                ? 'bg-blue-600 text-white'
                : 'text-slate-400 hover:text-white'
            }`}
          >
            {tab === 'achievements' && `🎖️ Достижения (${achievements.length})`}
            {tab === 'skills' && `🌳 Навыки`}
            {tab === 'tests' && `📝 Тесты (${quizResults.length})`}
          </button>
        ))}
      </div>

      {/* Achievements tab */}
      {activeTab === 'achievements' && (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {achievements.length === 0 ? (
            <div className="col-span-2 text-center text-slate-400 py-8">
              <div className="text-4xl mb-2">🎖️</div>
              <p>Достижений пока нет</p>
            </div>
          ) : (
            achievements.map(ua => (
              <div key={ua.id} className="card flex items-center gap-4">
                <div className="text-3xl">{ua.achievement.icon_url || '🏅'}</div>
                <div className="flex-1">
                  <div className="font-semibold text-white">{ua.achievement.name}</div>
                  <div className="text-xs text-slate-400">{ua.achievement.description}</div>
                  <div className="flex items-center gap-2 mt-1">
                    <span className="text-xs text-green-400">+{ua.achievement.xp_reward} XP</span>
                    <span className="text-xs text-slate-500">
                      {new Date(ua.earned_at).toLocaleDateString('ru-RU')}
                    </span>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      )}

      {/* Skills tab */}
      {activeTab === 'skills' && (
        <div className="space-y-4">
          {/* Earned / In Progress */}
          {earnedSkills.length > 0 && (
            <div>
              <h3 className="text-sm font-semibold text-slate-400 uppercase tracking-wide mb-2">
                Открытые навыки ({earnedSkills.length})
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                {earnedSkills.map(us => (
                  <div key={us.skill_id || us.skill?.id} className="flex items-center justify-between bg-slate-800 rounded-lg px-4 py-2.5">
                    <div className="flex items-center gap-2">
                      <span className="text-lg">{us.skill?.icon_url || '⭐'}</span>
                      <span className="text-white text-sm">{us.skill?.name}</span>
                    </div>
                    <span className={`text-xs px-2 py-0.5 rounded-full ${STATUS_COLORS[us.status]}`}>
                      {STATUS_LABELS[us.status]}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Available */}
          {availableSkills.length > 0 && (
            <div>
              <h3 className="text-sm font-semibold text-slate-400 uppercase tracking-wide mb-2">
                Доступны для изучения ({availableSkills.length})
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                {availableSkills.map(us => (
                  <div key={us.skill_id || us.skill?.id} className="flex items-center justify-between bg-slate-800 rounded-lg px-4 py-2.5">
                    <div className="flex items-center gap-2">
                      <span className="text-lg">{us.skill?.icon_url || '✨'}</span>
                      <span className="text-white text-sm">{us.skill?.name}</span>
                    </div>
                    <span className={`text-xs px-2 py-0.5 rounded-full ${STATUS_COLORS[us.status]}`}>
                      {STATUS_LABELS[us.status]}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Locked */}
          {lockedSkills.length > 0 && (
            <div>
              <h3 className="text-sm font-semibold text-slate-400 uppercase tracking-wide mb-2">
                Заблокированы ({lockedSkills.length})
              </h3>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                {lockedSkills.map(us => (
                  <div key={us.skill_id || us.skill?.id} className="flex items-center justify-between bg-slate-800/50 rounded-lg px-4 py-2.5 opacity-60">
                    <div className="flex items-center gap-2">
                      <span className="text-lg">{us.skill?.icon_url || '🔒'}</span>
                      <div>
                        <span className="text-slate-400 text-sm">{us.skill?.name}</span>
                        <div className="text-xs text-slate-500">
                          Требуется: {us.skill?.required_xp?.toLocaleString()} XP
                        </div>
                      </div>
                    </div>
                    <span className={`text-xs px-2 py-0.5 rounded-full ${STATUS_COLORS.locked}`}>
                      {STATUS_LABELS.locked}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* Tests tab */}
      {activeTab === 'tests' && (
        <div className="space-y-3">
          {quizResults.length === 0 ? (
            <div className="text-center text-slate-400 py-8">
              <div className="text-4xl mb-2">📝</div>
              <p>Тесты ещё не пройдены</p>
              <Link href="/quizzes" className="text-blue-400 hover:text-blue-300 mt-2 inline-block">
                Перейти к тестам →
              </Link>
            </div>
          ) : (
            quizResults.map(r => (
              <div key={r.id} className="card flex items-center justify-between">
                <div>
                  <div className="flex items-center gap-2">
                    <span className="font-medium text-white">{r.quiz?.title || 'Тест'}</span>
                    {r.is_perfect && <span className="text-xs bg-purple-900 text-purple-300 px-1.5 py-0.5 rounded">💎 100%</span>}
                    {r.quiz && (
                      <span className={`text-xs px-1.5 py-0.5 rounded ${DIFFICULTY_COLORS[r.quiz.difficulty]}`}>
                        {DIFFICULTY_LABELS[r.quiz.difficulty]}
                      </span>
                    )}
                  </div>
                  <div className="text-xs text-slate-400 mt-0.5">
                    {new Date(r.completed_at).toLocaleDateString('ru-RU', {
                      day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit'
                    })}
                  </div>
                </div>
                <div className="text-right">
                  <div className={`text-xl font-bold ${r.score >= 70 ? 'text-green-400' : 'text-red-400'}`}>
                    {r.score}%
                  </div>
                  <div className="text-xs text-slate-400">+{r.xp_earned} XP</div>
                </div>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
}
