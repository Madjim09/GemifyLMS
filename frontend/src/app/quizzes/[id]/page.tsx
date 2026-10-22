'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { api, Quiz, SubmitQuizOutput, DIFFICULTY_LABELS, DIFFICULTY_COLORS } from '@/lib/api';
import { useAuth } from '@/contexts/AuthContext';
import Link from 'next/link';

export default function TakeQuizPage() {
  const params = useParams();
  const router = useRouter();
  const quizID = params.id as string;
  const { user } = useAuth();

  const [quiz, setQuiz] = useState<Quiz | null>(null);
  const [answers, setAnswers] = useState<Record<string, string>>({});
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [result, setResult] = useState<SubmitQuizOutput | null>(null);
  const [error, setError] = useState('');

  useEffect(() => {
    if (!user) {
      router.push('/auth');
      return;
    }
    api.get(`/api/v1/quizzes/${quizID}`)
      .then(res => setQuiz(res.data.data))
      .catch(() => setError('Тест не найден'))
      .finally(() => setLoading(false));
  }, [quizID, user, router]);

  const handleAnswer = (questionID: string, optionID: string) => {
    setAnswers(prev => ({ ...prev, [questionID]: optionID }));
  };

  const handleSubmit = async () => {
    if (!quiz) return;
    const unanswered = quiz.questions?.filter(q => !answers[q.id]) || [];
    if (unanswered.length > 0) {
      setError(`Ответьте на все вопросы (осталось: ${unanswered.length})`);
      return;
    }
    setError('');
    setSubmitting(true);
    try {
      const res = await api.post<{ data: SubmitQuizOutput }>(`/api/v1/quizzes/${quizID}/submit`, { answers });
      setResult(res.data.data);
    } catch (err: unknown) {
      const axiosErr = err as { response?: { data?: { error?: string } } };
      setError(axiosErr?.response?.data?.error || 'Ошибка при отправке результата');
    } finally {
      setSubmitting(false);
    }
  };

  const answeredCount = Object.keys(answers).length;
  const totalQuestions = quiz?.questions?.length || 0;

  if (loading) return <div className="text-center text-slate-400 py-16">Загрузка теста...</div>;
  if (error && !quiz) return (
    <div className="text-center py-16">
      <div className="text-red-400 mb-4">{error}</div>
      <Link href="/quizzes" className="text-blue-400 hover:text-blue-300">← Назад к тестам</Link>
    </div>
  );
  if (!quiz) return null;

  if (result) {
    return (
      <div className="max-w-xl mx-auto space-y-6">
        <div className="card text-center">
          <div className="text-6xl mb-4">
            {result.is_perfect ? '💎' : result.passed ? '✅' : '❌'}
          </div>
          <h2 className="text-2xl font-bold text-white mb-2">
            {result.is_perfect ? 'Безупречный результат!' : result.passed ? 'Тест пройден!' : 'Тест не пройден'}
          </h2>

          <div className="flex justify-center gap-6 my-5">
            <div className="text-center">
              <div className="text-4xl font-bold text-white">{result.score}%</div>
              <div className="text-slate-400 text-sm">Результат</div>
            </div>
            <div className="text-center">
              <div className={`text-4xl font-bold ${result.xp_earned > 0 ? 'text-green-400' : 'text-slate-500'}`}>
                +{result.xp_earned}
              </div>
              <div className="text-slate-400 text-sm">XP получено</div>
            </div>
          </div>

          {/* Score bar */}
          <div className="mb-4">
            <div className="h-4 bg-slate-700 rounded-full overflow-hidden">
              <div
                className={`h-full rounded-full transition-all ${
                  result.score >= 70 ? 'bg-gradient-to-r from-green-500 to-emerald-400' : 'bg-red-500'
                }`}
                style={{ width: `${result.score}%` }}
              />
            </div>
            <div className="flex justify-between text-xs text-slate-500 mt-1">
              <span>0%</span>
              <span className="text-yellow-400">Порог: {quiz.pass_threshold}%</span>
              <span>100%</span>
            </div>
          </div>

          {result.is_perfect && (
            <div className="bg-purple-900/40 border border-purple-700 rounded-lg p-3 mb-4 text-sm text-purple-300">
              🎉 Бонус за 100%: <strong>+{quiz.bonus_xp} XP</strong>
            </div>
          )}

          {!result.passed && (
            <div className="bg-slate-700/50 rounded-lg p-3 mb-4 text-sm text-slate-400">
              Для прохождения нужно набрать {quiz.pass_threshold}% и выше. Попробуйте ещё раз!
            </div>
          )}

          <div className="flex gap-3 justify-center mt-4">
            <button
              onClick={() => { setResult(null); setAnswers({}); }}
              className="px-4 py-2 bg-slate-600 hover:bg-slate-500 text-white rounded-lg transition-colors"
            >
              Пройти снова
            </button>
            <Link href="/quizzes" className="px-4 py-2 bg-blue-600 hover:bg-blue-500 text-white rounded-lg transition-colors">
              Все тесты
            </Link>
            {user && (
              <Link href={`/users/${user.id}`} className="px-4 py-2 bg-purple-600 hover:bg-purple-500 text-white rounded-lg transition-colors">
                Мой профиль
              </Link>
            )}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto space-y-6">
      {/* Header */}
      <div className="card">
        <div className="flex items-start justify-between">
          <div>
            <div className="flex items-center gap-2 mb-1">
              <span className={`text-xs px-2 py-1 rounded-full font-medium ${DIFFICULTY_COLORS[quiz.difficulty]}`}>
                {DIFFICULTY_LABELS[quiz.difficulty]}
              </span>
              <span className="text-xs text-slate-400">Порог: {quiz.pass_threshold}%</span>
            </div>
            <h1 className="text-xl font-bold text-white">{quiz.title}</h1>
            <p className="text-slate-400 text-sm mt-1">{quiz.description}</p>
          </div>
          <div className="text-right ml-4">
            <div className="text-green-400 font-bold">+{quiz.xp_reward} XP</div>
            <div className="text-purple-400 text-xs">+{quiz.bonus_xp} за 100%</div>
          </div>
        </div>

        {/* Progress */}
        <div className="mt-3">
          <div className="flex justify-between text-xs text-slate-400 mb-1">
            <span>Отвечено: {answeredCount} / {totalQuestions}</span>
            <span>{Math.round((answeredCount / Math.max(totalQuestions, 1)) * 100)}%</span>
          </div>
          <div className="h-1.5 bg-slate-700 rounded-full overflow-hidden">
            <div
              className="h-full bg-blue-500 rounded-full transition-all"
              style={{ width: `${(answeredCount / Math.max(totalQuestions, 1)) * 100}%` }}
            />
          </div>
        </div>
      </div>

      {/* Questions */}
      {quiz.questions?.map((q, qIdx) => (
        <div key={q.id} className={`card border ${answers[q.id] ? 'border-blue-700/50' : 'border-slate-700'}`}>
          <div className="flex items-start gap-3 mb-4">
            <div className="w-7 h-7 rounded-full bg-slate-700 flex items-center justify-center text-sm font-bold text-slate-300 shrink-0">
              {qIdx + 1}
            </div>
            <p className="text-white font-medium leading-relaxed">{q.text}</p>
          </div>
          <div className="space-y-2 ml-10">
            {q.options.map(opt => (
              <label
                key={opt.id}
                className={`flex items-center gap-3 p-3 rounded-lg cursor-pointer transition-all ${
                  answers[q.id] === opt.id
                    ? 'bg-blue-900/50 border border-blue-600'
                    : 'bg-slate-700/50 border border-slate-600 hover:border-slate-500 hover:bg-slate-700'
                }`}
              >
                <div className={`w-4 h-4 rounded-full border-2 flex items-center justify-center shrink-0 ${
                  answers[q.id] === opt.id ? 'border-blue-400' : 'border-slate-500'
                }`}>
                  {answers[q.id] === opt.id && <div className="w-2 h-2 rounded-full bg-blue-400" />}
                </div>
                <input
                  type="radio"
                  name={q.id}
                  value={opt.id}
                  checked={answers[q.id] === opt.id}
                  onChange={() => handleAnswer(q.id, opt.id)}
                  className="sr-only"
                />
                <span className="text-slate-300 text-sm">{opt.text}</span>
              </label>
            ))}
          </div>
        </div>
      ))}

      {error && (
        <div className="bg-red-900/40 border border-red-700 rounded-lg px-4 py-2 text-red-400 text-sm">
          {error}
        </div>
      )}

      {/* Submit */}
      <div className="flex gap-3">
        <Link href="/quizzes" className="px-4 py-2.5 bg-slate-600 hover:bg-slate-500 text-white rounded-lg transition-colors">
          ← Назад
        </Link>
        <button
          onClick={handleSubmit}
          disabled={submitting || answeredCount < totalQuestions}
          className="flex-1 px-6 py-2.5 bg-blue-600 hover:bg-blue-500 disabled:opacity-50 disabled:cursor-not-allowed text-white font-semibold rounded-lg transition-colors"
        >
          {submitting ? 'Отправка...' : `Завершить тест (${answeredCount}/${totalQuestions})`}
        </button>
      </div>
    </div>
  );
}
