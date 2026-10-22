'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { api, Quiz, DIFFICULTY_LABELS, DIFFICULTY_COLORS } from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';

export default function QuizzesPage() {
  const [quizzes, setQuizzes] = useState<Quiz[]>([]);
  const [loading, setLoading] = useState(true);
  const { user } = useAuth();

  useEffect(() => {
    api.get('/api/v1/quizzes')
      .then(res => setQuizzes(res.data.data || []))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div className="text-center text-slate-400 py-16">Загрузка тестов...</div>;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-white">📝 Тесты</h1>
        <p className="text-slate-400 mt-1">
          Проверьте свои знания и получите XP. Порог прохождения: 70%.
          За 100% результат — бонусные XP!
        </p>
      </div>

      {/* Legend */}
      <div className="card">
        <div className="flex flex-wrap gap-4 text-sm">
          <div className="flex items-center gap-2">
            <span className="w-3 h-3 rounded-full bg-green-500 inline-block"></span>
            <span className="text-slate-300">Лёгкий</span>
          </div>
          <div className="flex items-center gap-2">
            <span className="w-3 h-3 rounded-full bg-yellow-500 inline-block"></span>
            <span className="text-slate-300">Средний</span>
          </div>
          <div className="flex items-center gap-2">
            <span className="w-3 h-3 rounded-full bg-red-500 inline-block"></span>
            <span className="text-slate-300">Сложный</span>
          </div>
          <div className="ml-auto text-slate-400">
            ≥70% → базовые XP | 100% → +бонусные XP
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
        {quizzes.map(quiz => (
          <div key={quiz.id} className="card flex flex-col">
            <div className="flex items-start justify-between mb-3">
              <span className={`text-xs px-2 py-1 rounded-full font-medium ${DIFFICULTY_COLORS[quiz.difficulty] || 'bg-slate-700 text-slate-400'}`}>
                {DIFFICULTY_LABELS[quiz.difficulty] || quiz.difficulty}
              </span>
              <div className="text-right">
                <div className="text-green-400 font-bold text-sm">+{quiz.xp_reward} XP</div>
                <div className="text-purple-400 text-xs">+{quiz.bonus_xp} XP бонус</div>
              </div>
            </div>

            <h3 className="text-lg font-bold text-white mb-2">{quiz.title}</h3>
            <p className="text-slate-400 text-sm flex-1 mb-4">{quiz.description}</p>

            <div className="flex items-center justify-between">
              <div className="text-xs text-slate-500">
                Порог: {quiz.pass_threshold}%
              </div>
              {user ? (
                <Link
                  href={`/quizzes/${quiz.id}`}
                  className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white text-sm font-medium rounded-lg transition-colors"
                >
                  Начать тест
                </Link>
              ) : (
                <Link
                  href="/auth"
                  className="px-4 py-2 bg-slate-600 hover:bg-slate-500 text-slate-300 text-sm font-medium rounded-lg transition-colors"
                >
                  Войти для сдачи
                </Link>
              )}
            </div>
          </div>
        ))}
      </div>

      {quizzes.length === 0 && (
        <div className="text-center text-slate-400 py-16">
          <div className="text-5xl mb-4">📝</div>
          <p>Тесты пока не добавлены</p>
        </div>
      )}
    </div>
  );
}
