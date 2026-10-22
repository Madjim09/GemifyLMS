import { Skill } from '@/lib/api';

interface SkillNodeProps {
  skill: Skill;
  depth: number;
}

export function SkillNode({ skill, depth }: SkillNodeProps) {
  const colors = [
    'border-blue-500 bg-blue-900/20',
    'border-green-500 bg-green-900/20',
    'border-purple-500 bg-purple-900/20',
    'border-amber-500 bg-amber-900/20',
  ];
  const color = colors[depth % colors.length];

  return (
    <div className="w-full max-w-2xl">
      <div className={`border rounded-xl p-4 ${color}`}>
        <div className="flex items-center justify-between">
          <div>
            <h3 className="font-bold text-white">{skill.icon_url || '📚'} {skill.name}</h3>
            <p className="text-slate-400 text-sm mt-1">{skill.description}</p>
          </div>
          <div className="text-right ml-4">
            <div className="text-amber-400 text-sm font-bold">
              {skill.required_xp > 0 ? `${skill.required_xp.toLocaleString()} XP` : 'Бесплатно'}
            </div>
          </div>
        </div>
      </div>
      {skill.children && skill.children.length > 0 && (
        <div className="ml-8 mt-2 space-y-2 border-l-2 border-slate-700 pl-4">
          {skill.children.map(child => (
            <SkillNode key={child.id} skill={child} depth={depth + 1} />
          ))}
        </div>
      )}
    </div>
  );
}
