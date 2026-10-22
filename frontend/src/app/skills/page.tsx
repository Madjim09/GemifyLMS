'use client';

import { useEffect, useState } from 'react';
import { api, Skill } from '@/lib/api';
import { SkillNode } from '@/components/SkillNode';

export default function SkillsPage() {
  const [skills, setSkills] = useState<Skill[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.get('/api/v1/skills')
      .then(res => setSkills(res.data.data || []))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-white">🌳 Дерево навыков</h1>
      <p className="text-slate-400">Открывайте новые навыки, накапливая опыт</p>

      {loading ? (
        <div className="text-center text-slate-400 py-16">Загрузка...</div>
      ) : (
        <div className="flex flex-col items-center gap-4">
          {skills.map(skill => (
            <SkillNode key={skill.id} skill={skill} depth={0} />
          ))}
        </div>
      )}
    </div>
  );
}
