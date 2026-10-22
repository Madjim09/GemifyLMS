'use client';

import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import clsx from 'clsx';
import { useAuth } from '@/contexts/AuthContext';
import { getXPProgress, getXPForNextLevel, LEVEL_THRESHOLDS } from '@/lib/api';

const links = [
  { href: '/', label: '🏠 Главная' },
  { href: '/quizzes', label: '📝 Тесты' },
  { href: '/leaderboard', label: '🏆 Рейтинг' },
  { href: '/achievements', label: '🎖️ Достижения' },
  { href: '/skills', label: '🌳 Навыки' },
  { href: '/users', label: '👥 Пользователи' },
];

export function Navigation() {
  const pathname = usePathname();
  const router = useRouter();
  const { user, logout } = useAuth();

  const handleLogout = () => {
    logout();
    router.push('/auth');
  };

  const xpProgress = user ? getXPProgress(user.xp, user.level) : 0;
  const nextLevelXP = user ? getXPForNextLevel(user.level) : 0;
  const currentLevelXP = user ? LEVEL_THRESHOLDS[user.level - 1] || 0 : 0;

  return (
    <nav className="bg-slate-800 border-b border-slate-700 sticky top-0 z-50">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex items-center gap-1 h-14">
          <div className="font-bold text-white mr-3 whitespace-nowrap text-sm">
            🎮 GamifyLMS
          </div>
          <div className="flex items-center gap-1 flex-1 overflow-x-auto">
            {links.map(link => (
              <Link
                key={link.href}
                href={link.href}
                className={clsx(
                  'px-2.5 py-1.5 rounded-lg text-xs font-medium whitespace-nowrap transition-colors',
                  pathname === link.href
                    ? 'bg-blue-600 text-white'
                    : 'text-slate-400 hover:text-white hover:bg-slate-700'
                )}
              >
                {link.label}
              </Link>
            ))}
          </div>

          {/* User info / auth */}
          <div className="flex items-center gap-3 ml-2 shrink-0">
            {user ? (
              <>
                {/* Level + XP Bar */}
                <div className="hidden md:flex items-center gap-2">
                  <span className="level-badge text-xs px-2 py-1">Ур. {user.level}</span>
                  <div className="w-32">
                    <div className="flex justify-between text-xs text-slate-400 mb-0.5">
                      <span>{user.xp.toLocaleString()} XP</span>
                      <span>{nextLevelXP > 0 ? nextLevelXP.toLocaleString() : '∞'}</span>
                    </div>
                    <div className="h-1.5 bg-slate-600 rounded-full overflow-hidden">
                      <div
                        className="h-full bg-gradient-to-r from-blue-500 to-purple-500 rounded-full transition-all"
                        style={{ width: `${xpProgress}%` }}
                      />
                    </div>
                  </div>
                </div>

                {/* Profile link */}
                <Link
                  href={`/users/${user.id}`}
                  className="flex items-center gap-1.5 text-slate-300 hover:text-white transition-colors"
                >
                  <div className="w-7 h-7 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white text-xs font-bold">
                    {user.name.charAt(0).toUpperCase()}
                  </div>
                  <span className="hidden lg:block text-xs font-medium max-w-24 truncate">{user.name}</span>
                </Link>
                <button
                  onClick={handleLogout}
                  className="text-xs text-slate-500 hover:text-slate-300 transition-colors"
                >
                  Выйти
                </button>
              </>
            ) : (
              <Link
                href="/auth"
                className="px-3 py-1.5 bg-blue-600 hover:bg-blue-500 text-white text-xs font-medium rounded-lg transition-colors whitespace-nowrap"
              >
                Войти
              </Link>
            )}
          </div>
        </div>

        {/* Mobile XP bar (shown when logged in) */}
        {user && (
          <div className="md:hidden pb-2 flex items-center gap-2">
            <span className="level-badge text-xs px-1.5 py-0.5">Ур. {user.level}</span>
            <div className="flex-1">
              <div className="h-1.5 bg-slate-600 rounded-full overflow-hidden">
                <div
                  className="h-full bg-gradient-to-r from-blue-500 to-purple-500 rounded-full transition-all"
                  style={{ width: `${xpProgress}%` }}
                />
              </div>
            </div>
            <span className="text-xs text-slate-400">{user.xp.toLocaleString()} XP</span>
          </div>
        )}
      </div>
    </nav>
  );
}
