package database

import (
	"database/sql"
	"log"
)

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createUsersTable,
		alterUsersAddPasswordHash,
		alterUsersAddTestsPassed,
		createXPTransactionsTable,
		createAchievementsTable,
		createUserAchievementsTable,
		createSkillsTable,
		alterSkillsAddCategory,
		createUserSkillsTable,
		createEventRecordsTable,
		createQuizzesTable,
		createQuizQuestionsTable,
		createQuizOptionsTable,
		createUserQuizResultsTable,
		insertDefaultAchievements,
		insertRootSkills,
		insertSoftSkills,
		insertHardSkills,
		insertInstrumentalSkills,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			log.Printf("Migration warning: %v", err)
		}
	}

	log.Println("Migrations completed")
	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	external_id VARCHAR(255) UNIQUE NOT NULL,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	password_hash VARCHAR(255) DEFAULT '',
	xp BIGINT DEFAULT 0,
	level INTEGER DEFAULT 1,
	tests_passed INTEGER DEFAULT 0,
	avatar_url TEXT DEFAULT '',
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_xp ON users(xp DESC);
CREATE INDEX IF NOT EXISTS idx_users_external_id ON users(external_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
`

const alterUsersAddPasswordHash = `
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255) DEFAULT '';
`

const alterUsersAddTestsPassed = `
ALTER TABLE users ADD COLUMN IF NOT EXISTS tests_passed INTEGER DEFAULT 0;
`

const createXPTransactionsTable = `
CREATE TABLE IF NOT EXISTS xp_transactions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	amount BIGINT NOT NULL,
	reason VARCHAR(255) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_xp_transactions_user_id ON xp_transactions(user_id);
`

const createAchievementsTable = `
CREATE TABLE IF NOT EXISTS achievements (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) UNIQUE NOT NULL,
	description TEXT NOT NULL,
	icon_url TEXT DEFAULT '',
	condition_type VARCHAR(100) NOT NULL,
	condition_value INTEGER NOT NULL DEFAULT 1,
	xp_reward BIGINT DEFAULT 0,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
`

const createUserAchievementsTable = `
CREATE TABLE IF NOT EXISTS user_achievements (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	achievement_id UUID NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
	earned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	UNIQUE(user_id, achievement_id)
);
CREATE INDEX IF NOT EXISTS idx_user_achievements_user_id ON user_achievements(user_id);
`

const createSkillsTable = `
CREATE TABLE IF NOT EXISTS skills (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) UNIQUE NOT NULL,
	description TEXT DEFAULT '',
	category VARCHAR(100) DEFAULT '',
	parent_id UUID REFERENCES skills(id) ON DELETE SET NULL,
	required_xp BIGINT DEFAULT 0,
	icon_url TEXT DEFAULT '',
	"order" INTEGER DEFAULT 0,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
`

const alterSkillsAddCategory = `
ALTER TABLE skills ADD COLUMN IF NOT EXISTS category VARCHAR(100) DEFAULT '';
`

const createUserSkillsTable = `
CREATE TABLE IF NOT EXISTS user_skills (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
	status VARCHAR(50) DEFAULT 'locked',
	progress INTEGER DEFAULT 0,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	UNIQUE(user_id, skill_id)
);
`

const createEventRecordsTable = `
CREATE TABLE IF NOT EXISTS event_records (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	event_type VARCHAR(100) NOT NULL,
	object_id TEXT DEFAULT '',
	xp_awarded BIGINT DEFAULT 0,
	raw_payload JSONB,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_event_records_user_id ON event_records(user_id);
CREATE INDEX IF NOT EXISTS idx_event_records_event_type ON event_records(event_type);
CREATE INDEX IF NOT EXISTS idx_event_records_created_at ON event_records(created_at);
`

const createQuizzesTable = `
CREATE TABLE IF NOT EXISTS quizzes (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	title VARCHAR(255) NOT NULL,
	description TEXT DEFAULT '',
	difficulty VARCHAR(50) DEFAULT 'medium',
	xp_reward BIGINT DEFAULT 100,
	bonus_xp BIGINT DEFAULT 50,
	pass_threshold INTEGER DEFAULT 70,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
`

const createQuizQuestionsTable = `
CREATE TABLE IF NOT EXISTS quiz_questions (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
	text TEXT NOT NULL,
	"order" INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_quiz_questions_quiz_id ON quiz_questions(quiz_id);
`

const createQuizOptionsTable = `
CREATE TABLE IF NOT EXISTS quiz_options (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	question_id UUID NOT NULL REFERENCES quiz_questions(id) ON DELETE CASCADE,
	text TEXT NOT NULL,
	is_correct BOOLEAN DEFAULT FALSE,
	"order" INTEGER DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_quiz_options_question_id ON quiz_options(question_id);
`

const createUserQuizResultsTable = `
CREATE TABLE IF NOT EXISTS user_quiz_results (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
	score INTEGER NOT NULL,
	xp_earned BIGINT DEFAULT 0,
	is_perfect BOOLEAN DEFAULT FALSE,
	completed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_user_quiz_results_user_id ON user_quiz_results(user_id);
`

const insertDefaultAchievements = `
INSERT INTO achievements (name, description, icon_url, condition_type, condition_value, xp_reward)
VALUES
	('Первопроходец', 'Сделайте первый шаг — наберите 1 XP', '🌱', 'total_xp', 1, 10),
	('Первые шаги', 'Наберите 500 XP', '🚶', 'total_xp', 500, 50),
	('Опытный ученик', 'Наберите 3000 XP', '📚', 'total_xp', 3000, 150),
	('Мастер знаний', 'Наберите 10000 XP', '🏆', 'total_xp', 10000, 300),
	('Легенда платформы', 'Наберите 30000 XP', '⚡', 'total_xp', 30000, 500),
	('Тест пройден', 'Успешно сдайте первый тест', '✅', 'test_count', 1, 25),
	('Серийный тестировщик', 'Успешно сдайте 5 тестов', '📊', 'test_count', 5, 100),
	('Мастер тестов', 'Успешно сдайте 10 тестов', '🎯', 'test_count', 10, 200),
	('Безупречный результат', 'Сдайте тест на 100%%', '💎', 'test_perfect', 1, 150),
	('Недельный марафон', 'Завершите 5 модулей за неделю', '🏃', 'modules_per_week', 5, 150),
	('Киноман', 'Посмотрите 10 учебных видео', '🎬', 'video_count', 10, 75),
	('Огненная серия', 'Заходите на платформу 7 дней подряд', '🔥', 'streak_days', 7, 100)
ON CONFLICT (name) DO NOTHING;
`

const insertRootSkills = `
INSERT INTO skills (name, description, category, required_xp, icon_url, "order")
VALUES
	('Soft Skills', 'Межличностные и коммуникативные навыки', 'root', 0, '🤝', 1),
	('Hard Skills', 'Профессиональные и аналитические навыки', 'root', 0, '💼', 2),
	('Инструментальные навыки', 'Инструменты, технологии и методологии', 'root', 0, '🔧', 3)
ON CONFLICT (name) DO NOTHING;
`

const insertSoftSkills = `
INSERT INTO skills (name, description, category, parent_id, required_xp, icon_url, "order")
SELECT 'Коммуникации', 'Навыки деловой и межличностной коммуникации', 'soft', id, 500, '💬', 1
FROM skills WHERE name = 'Soft Skills'
ON CONFLICT (name) DO NOTHING;

INSERT INTO skills (name, description, category, parent_id, required_xp, icon_url, "order")
SELECT 'Работа в команде', 'Эффективное взаимодействие в коллективе', 'soft', id, 500, '👥', 2
FROM skills WHERE name = 'Soft Skills'
ON CONFLICT (name) DO NOTHING;

INSERT INTO skills (name, description, category, parent_id, required_xp, icon_url, "order")
SELECT 'Эмоциональный интеллект', 'Управление эмоциями и эмпатия', 'soft', id, 1500, '❤️', 3
FROM skills WHERE name = 'Soft Skills'
ON CONFLICT (name) DO NOTHING;

INSERT INTO skills (name, description, category, parent_id, required_xp, icon_url, "order")
SELECT 'Лидерство', 'Развитие лидерских качеств и мотивация команды', 'soft', id, 5000, '👑', 4
FROM skills WHERE name = 'Soft Skills'
ON CONFLICT (name) DO NOTHING;
`

const insertHardSkills = `
INSERT INTO skills (name, description, category, parent_id, required_xp, icon_url, "order")
SELECT 'Управление проектами', 'Планирование и контроль реализации проектов', 'hard', id, 1500, '📋', 1
FROM skills WHERE name = 'Hard Skills'
ON CONFLICT (name) DO NOTHING;

INSERT INTO skills (name, description, category, parent_id, required_xp, icon_url, "order")
SELECT 'Финансовая грамотность', 'Основы финансового анализа и планирования', 'hard', id, 3000, '💰', 2
FROM skills WHERE name = 'Hard Skills'
ON CONFLICT (name) DO NOTHING;

INSERT INTO skills (name, description, category, parent_id, required_xp, icon_url, "order")
SELECT 'Анализ данных', 'Работа с данными и принятие решений на основе данных', 'hard', id, 5000, '📊', 3
FROM skills WHERE name = 'Hard Skills'
ON CONFLICT (name) DO NOTHING;

INSERT INTO skills (name, description, category, parent_id, required_xp, icon_url, "order")
SELECT 'Стратегическое мышление', 'Долгосрочное планирование и стратегический анализ', 'hard', id, 9000, '🎯', 4
FROM skills WHERE name = 'Hard Skills'
ON CONFLICT (name) DO NOTHING;
`

const insertInstrumentalSkills = `
INSERT INTO skills (name, description, category, parent_id, required_xp, icon_url, "order")
SELECT 'MS Office', 'Работа с пакетом Microsoft Office', 'instrumental', id, 0, '📄', 1
FROM skills WHERE name = 'Инструментальные навыки'
ON CONFLICT (name) DO NOTHING;

INSERT INTO skills (name, description, category, parent_id, required_xp, icon_url, "order")
SELECT 'Agile методологии', 'Гибкие методологии разработки и управления', 'instrumental', id, 1500, '🔄', 2
FROM skills WHERE name = 'Инструментальные навыки'
ON CONFLICT (name) DO NOTHING;

INSERT INTO skills (name, description, category, parent_id, required_xp, icon_url, "order")
SELECT 'Работа с CRM', 'Управление взаимоотношениями с клиентами', 'instrumental', id, 3000, '🗂️', 3
FROM skills WHERE name = 'Инструментальные навыки'
ON CONFLICT (name) DO NOTHING;

INSERT INTO skills (name, description, category, parent_id, required_xp, icon_url, "order")
SELECT 'Основы программирования', 'Логическое мышление и базовые навыки программирования', 'instrumental', id, 5000, '💻', 4
FROM skills WHERE name = 'Инструментальные навыки'
ON CONFLICT (name) DO NOTHING;
`
