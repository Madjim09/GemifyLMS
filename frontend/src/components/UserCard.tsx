import { User, getXPProgress } from '@/lib/api';

interface UserCardProps {
  user: User;
}

export function UserCard({ user }: UserCardProps) {
  const progress = getXPProgress(user.xp, user.level);

  return (
    <div className="flex items-center gap-3">
      <div className="w-10 h-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-bold">
        {user.name.charAt(0).toUpperCase()}
      </div>
      <div className="flex-1 min-w-0">
        <div className="font-semibold text-white truncate">{user.name}</div>
        <div className="h-1.5 bg-slate-700 rounded-full overflow-hidden mt-1">
          <div
            className="h-full bg-gradient-to-r from-blue-500 to-purple-500 rounded-full"
            style={{ width: `${progress}%` }}
          />
        </div>
      </div>
    </div>
  );
}
