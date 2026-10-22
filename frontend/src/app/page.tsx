import Link from "next/link";

export default function Home() {
  return (
    <div className="space-y-8">
      <div className="text-center py-16">
        <h1 className="text-5xl font-bold text-white mb-4">
          🎮 Gamification Platform
        </h1>
        <p className="text-xl text-slate-400 max-w-2xl mx-auto">
          Платформа геймификации корпоративного обучения. Зарабатывай XP, открывай достижения, соревнуйся с коллегами.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Link href="/leaderboard" className="card hover:border-blue-500 transition-colors group">
          <div className="text-4xl mb-4">🏆</div>
          <h2 className="text-xl font-bold text-white group-hover:text-blue-400 transition-colors">Рейтинг</h2>
          <p className="text-slate-400 mt-2">Топ сотрудников по очкам опыта</p>
        </Link>

        <Link href="/achievements" className="card hover:border-amber-500 transition-colors group">
          <div className="text-4xl mb-4">🎖️</div>
          <h2 className="text-xl font-bold text-white group-hover:text-amber-400 transition-colors">Достижения</h2>
          <p className="text-slate-400 mt-2">Коллекция цифровых бейджей</p>
        </Link>

        <Link href="/skills" className="card hover:border-green-500 transition-colors group">
          <div className="text-4xl mb-4">🌳</div>
          <h2 className="text-xl font-bold text-white group-hover:text-green-400 transition-colors">Дерево навыков</h2>
          <p className="text-slate-400 mt-2">Визуализация прогресса обучения</p>
        </Link>
      </div>

      {/* <div className="card">
        <h2 className="text-xl font-bold text-white mb-4">Быстрый тест (xAPI Event)</h2>
        <p className="text-slate-400 mb-4">Отправь тестовое событие через webhook:</p>
        <pre className="bg-slate-900 rounded-lg p-4 text-sm text-green-400 overflow-x-auto">
{`curl -X POST http://localhost:8080/webhook/event \\
  -H "Content-Type: application/json" \\
  -d '{
    "actor": {
      "name": "Иван Иванов",
      "mbox": "mailto:ivan@company.com"
    },
    "verb": {
      "id": "http://adlnet.gov/expapi/verbs/passed",
      "display": {"ru-RU": "прошёл"}
    },
    "object": {
      "id": "https://lms.company.com/courses/1",
      "objectType": "Activity",
      "definition": {
        "name": {"ru-RU": "Основы управления проектами"}
      }
    }
  }'`}
        </pre>
      </div> */}
    </div>
  );
}
