package database

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	quiz1ID = uuid.MustParse("a0000000-0000-0000-0000-000000000001")
	quiz2ID = uuid.MustParse("a0000000-0000-0000-0000-000000000002")
	quiz3ID = uuid.MustParse("a0000000-0000-0000-0000-000000000003")
	quiz4ID = uuid.MustParse("a0000000-0000-0000-0000-000000000004")
	quiz5ID = uuid.MustParse("a0000000-0000-0000-0000-000000000005")
)

func SeedDemoData(db *sql.DB) {
	log.Println("Seeding quiz catalog...")
	if err := seedQuizzes(db); err != nil {
		log.Printf("Quiz seed warning: %v", err)
	}

	log.Println("Seeding demo users...")
	demoPassword, _ := bcrypt.GenerateFromPassword([]byte("demo123"), bcrypt.DefaultCost)
	hash := string(demoPassword)

	users := []struct {
		id    uuid.UUID
		name  string
		email string
		xp    int64
		level int
	}{
		{uuid.MustParse("b0000000-0000-0000-0000-000000000001"), "Алексей Новиков", "novice@demo.com", 200, 1},
		{uuid.MustParse("b0000000-0000-0000-0000-000000000002"), "Мария Иванова", "active@demo.com", 3500, 4},
		{uuid.MustParse("b0000000-0000-0000-0000-000000000003"), "Дмитрий Петров", "top@demo.com", 38000, 10},
		{uuid.MustParse("b0000000-0000-0000-0000-000000000004"), "Елена Смирнова", "mid@demo.com", 7500, 5},
	}

	for _, u := range users {
		_, err := db.Exec(
			`INSERT INTO users (id, external_id, name, email, password_hash, xp, level, tests_passed, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING`,
			u.id, "demo-"+u.id.String(), u.name, u.email, hash, u.xp, u.level,
			map[uuid.UUID]int{
				uuid.MustParse("b0000000-0000-0000-0000-000000000001"): 0,
				uuid.MustParse("b0000000-0000-0000-0000-000000000002"): 3,
				uuid.MustParse("b0000000-0000-0000-0000-000000000003"): 15,
				uuid.MustParse("b0000000-0000-0000-0000-000000000004"): 5,
			}[u.id],
		)
		if err != nil {
			log.Printf("Demo user seed warning (%s): %v", u.email, err)
		}
	}

	// Achievements for demo users
	user2 := uuid.MustParse("b0000000-0000-0000-0000-000000000002")
	user3 := uuid.MustParse("b0000000-0000-0000-0000-000000000003")
	user4 := uuid.MustParse("b0000000-0000-0000-0000-000000000004")

	awardAchievements(db, user2, []string{
		"Первопроходец", "Первые шаги", "Опытный ученик", "Тест пройден", "Серийный тестировщик",
	})
	awardAchievements(db, user3, []string{
		"Первопроходец", "Первые шаги", "Опытный ученик", "Мастер знаний", "Легенда платформы",
		"Тест пройден", "Серийный тестировщик", "Мастер тестов", "Безупречный результат",
		"Киноман", "Огненная серия",
	})
	awardAchievements(db, user4, []string{
		"Первопроходец", "Первые шаги", "Опытный ученик",
		"Тест пройден", "Серийный тестировщик",
	})

	// Skills for demo users
	seedUserSkills(db, user2, map[string]string{
		"MS Office": "mastered", "Коммуникации": "in_progress", "Soft Skills": "available",
		"Работа в команде": "available", "Agile методологии": "available",
	})
	seedUserSkills(db, user3, map[string]string{
		"Soft Skills": "mastered", "Hard Skills": "mastered", "Инструментальные навыки": "mastered",
		"Коммуникации": "mastered", "Лидерство": "mastered", "Работа в команде": "mastered",
		"Эмоциональный интеллект": "mastered", "Управление проектами": "mastered",
		"Финансовая грамотность": "mastered", "Анализ данных": "in_progress",
		"Стратегическое мышление": "available", "MS Office": "mastered",
		"Agile методологии": "mastered", "Работа с CRM": "in_progress",
		"Основы программирования": "available",
	})
	seedUserSkills(db, user4, map[string]string{
		"MS Office": "mastered", "Коммуникации": "in_progress",
		"Управление проектами": "available", "Agile методологии": "available",
		"Работа в команде": "in_progress",
	})

	// Quiz results for demo users
	seedQuizResult(db, user2, quiz1ID, 80, 100, false)
	seedQuizResult(db, user2, quiz2ID, 85, 100, false)
	seedQuizResult(db, user2, quiz3ID, 70, 100, false)
	seedQuizResult(db, user3, quiz1ID, 100, 150, true)
	seedQuizResult(db, user3, quiz2ID, 100, 150, true)
	seedQuizResult(db, user3, quiz3ID, 100, 150, true)
	seedQuizResult(db, user3, quiz4ID, 90, 200, false)
	seedQuizResult(db, user3, quiz5ID, 80, 150, false)
	seedQuizResult(db, user4, quiz1ID, 75, 100, false)
	seedQuizResult(db, user4, quiz3ID, 80, 100, false)

	log.Println("Demo data seeding complete")
}

func awardAchievements(db *sql.DB, userID uuid.UUID, names []string) {
	for _, name := range names {
		var achID uuid.UUID
		err := db.QueryRow(`SELECT id FROM achievements WHERE name = $1`, name).Scan(&achID)
		if err != nil {
			log.Printf("Achievement not found: %s", name)
			continue
		}
		_, err = db.Exec(
			`INSERT INTO user_achievements (id, user_id, achievement_id, earned_at)
			VALUES ($1, $2, $3, $4) ON CONFLICT (user_id, achievement_id) DO NOTHING`,
			uuid.New(), userID, achID, time.Now(),
		)
		if err != nil {
			log.Printf("Award achievement warning (%s): %v", name, err)
		}
	}
}

func seedUserSkills(db *sql.DB, userID uuid.UUID, skills map[string]string) {
	for skillName, status := range skills {
		var skillID uuid.UUID
		err := db.QueryRow(`SELECT id FROM skills WHERE name = $1`, skillName).Scan(&skillID)
		if err != nil {
			log.Printf("Skill not found: %s", skillName)
			continue
		}
		progress := 0
		if status == "mastered" {
			progress = 100
		} else if status == "in_progress" {
			progress = 50
		}
		_, err = db.Exec(
			`INSERT INTO user_skills (id, user_id, skill_id, status, progress, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
			ON CONFLICT (user_id, skill_id) DO NOTHING`,
			uuid.New(), userID, skillID, status, progress,
		)
		if err != nil {
			log.Printf("Seed user skill warning (%s): %v", skillName, err)
		}
	}
}

func seedQuizResult(db *sql.DB, userID, quizID uuid.UUID, score int, xpEarned int64, isPerfect bool) {
	_, err := db.Exec(
		`INSERT INTO user_quiz_results (id, user_id, quiz_id, score, xp_earned, is_perfect, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())`,
		uuid.New(), userID, quizID, score, xpEarned, isPerfect,
	)
	if err != nil {
		log.Printf("Seed quiz result warning: %v", err)
	}
}

func seedQuizzes(db *sql.DB) error {
	type quizDef struct {
		id          uuid.UUID
		title       string
		description string
		difficulty  string
		xpReward    int64
		bonusXP     int64
	}

	quizDefs := []quizDef{
		{quiz1ID, "Основы Agile и Scrum", "Проверьте знание принципов гибкой разработки и фреймворка Scrum", "easy", 100, 50},
		{quiz2ID, "Деловые коммуникации", "Тест на знание эффективных коммуникаций в деловой среде", "easy", 100, 50},
		{quiz3ID, "Управление временем", "Проверьте знание методов и инструментов тайм-менеджмента", "medium", 150, 75},
		{quiz4ID, "Основы лидерства", "Тест на знание стилей и принципов эффективного лидерства", "hard", 200, 100},
		{quiz5ID, "Финансовая грамотность", "Проверьте понимание базовых финансовых показателей и концепций", "medium", 150, 75},
	}

	for _, qd := range quizDefs {
		_, err := db.Exec(
			`INSERT INTO quizzes (id, title, description, difficulty, xp_reward, bonus_xp, pass_threshold)
			VALUES ($1, $2, $3, $4, $5, $6, 70) ON CONFLICT (id) DO NOTHING`,
			qd.id, qd.title, qd.description, qd.difficulty, qd.xpReward, qd.bonusXP,
		)
		if err != nil {
			return err
		}
	}

	if err := seedQuizQuestions(db); err != nil {
		return err
	}

	return nil
}

func seedQuizQuestions(db *sql.DB) error {
	type opt struct {
		id        uuid.UUID
		text      string
		isCorrect bool
		order     int
	}
	type question struct {
		id      uuid.UUID
		quizID  uuid.UUID
		text    string
		order   int
		options []opt
	}

	questions := []question{
		// Quiz 1: Agile
		{uuid.MustParse("c0000000-0000-0000-0001-000000000001"), quiz1ID,
			"Что такое спринт в Scrum?", 1, []opt{
				{uuid.MustParse("d0000000-0001-0001-0001-000000000001"), "Длительный период разработки без изменений", false, 1},
				{uuid.MustParse("d0000000-0001-0001-0001-000000000002"), "Короткий итерационный цикл разработки длиной 1–4 недели", true, 2},
				{uuid.MustParse("d0000000-0001-0001-0001-000000000003"), "Команда разработчиков", false, 3},
				{uuid.MustParse("d0000000-0001-0001-0001-000000000004"), "Список всех задач проекта", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0001-000000000002"), quiz1ID,
			"Кто отвечает за приоритизацию бэклога продукта?", 2, []opt{
				{uuid.MustParse("d0000000-0001-0002-0001-000000000001"), "Scrum-мастер", false, 1},
				{uuid.MustParse("d0000000-0001-0002-0001-000000000002"), "Команда разработки", false, 2},
				{uuid.MustParse("d0000000-0001-0002-0001-000000000003"), "Владелец продукта (Product Owner)", true, 3},
				{uuid.MustParse("d0000000-0001-0002-0001-000000000004"), "Технический директор", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0001-000000000003"), quiz1ID,
			"Что такое Daily Scrum?", 3, []opt{
				{uuid.MustParse("d0000000-0001-0003-0001-000000000001"), "Ежедневное 15-минутное совещание команды", true, 1},
				{uuid.MustParse("d0000000-0001-0003-0001-000000000002"), "Ежедневный отчёт для руководства", false, 2},
				{uuid.MustParse("d0000000-0001-0003-0001-000000000003"), "Ежедневное планирование на весь день", false, 3},
				{uuid.MustParse("d0000000-0001-0003-0001-000000000004"), "Ежедневная ретроспектива", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0001-000000000004"), quiz1ID,
			"Какой артефакт Scrum описывает все функции, требующие реализации?", 4, []opt{
				{uuid.MustParse("d0000000-0001-0004-0001-000000000001"), "Журнал спринта", false, 1},
				{uuid.MustParse("d0000000-0001-0004-0001-000000000002"), "Инкремент продукта", false, 2},
				{uuid.MustParse("d0000000-0001-0004-0001-000000000003"), "Бэклог продукта", true, 3},
				{uuid.MustParse("d0000000-0001-0004-0001-000000000004"), "Диаграмма сгорания задач", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0001-000000000005"), quiz1ID,
			"Какова основная цель ретроспективы спринта?", 5, []opt{
				{uuid.MustParse("d0000000-0001-0005-0001-000000000001"), "Показ результатов стейкхолдерам", false, 1},
				{uuid.MustParse("d0000000-0001-0005-0001-000000000002"), "Планирование следующего спринта", false, 2},
				{uuid.MustParse("d0000000-0001-0005-0001-000000000003"), "Улучшение процессов работы команды", true, 3},
				{uuid.MustParse("d0000000-0001-0005-0001-000000000004"), "Оценка производительности сотрудников", false, 4},
			}},

		// Quiz 2: Коммуникации
		{uuid.MustParse("c0000000-0000-0000-0002-000000000001"), quiz2ID,
			"Что такое активное слушание?", 1, []opt{
				{uuid.MustParse("d0000000-0002-0001-0001-000000000001"), "Записывание всего сказанного собеседником", false, 1},
				{uuid.MustParse("d0000000-0002-0001-0001-000000000002"), "Молчание во время разговора", false, 2},
				{uuid.MustParse("d0000000-0002-0001-0001-000000000003"), "Полное внимание, понимание и обратная связь", true, 3},
				{uuid.MustParse("d0000000-0002-0001-0001-000000000004"), "Перебивание для уточнений", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0002-000000000002"), quiz2ID,
			"Как лучше всего строить конструктивную обратную связь?", 2, []opt{
				{uuid.MustParse("d0000000-0002-0002-0001-000000000001"), "Сразу указать на ошибки без предисловий", false, 1},
				{uuid.MustParse("d0000000-0002-0002-0001-000000000002"), "Сначала отметить сильные стороны, затем — зоны улучшения", true, 2},
				{uuid.MustParse("d0000000-0002-0002-0001-000000000003"), "Отложить разговор на неопределённое время", false, 3},
				{uuid.MustParse("d0000000-0002-0002-0001-000000000004"), "Написать список всех ошибок без комментариев", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0002-000000000003"), quiz2ID,
			"Что такое невербальная коммуникация?", 3, []opt{
				{uuid.MustParse("d0000000-0002-0003-0001-000000000001"), "Письменное деловое общение", false, 1},
				{uuid.MustParse("d0000000-0002-0003-0001-000000000002"), "Общение через мессенджеры", false, 2},
				{uuid.MustParse("d0000000-0002-0003-0001-000000000003"), "Язык тела, мимика, жесты и интонация", true, 3},
				{uuid.MustParse("d0000000-0002-0003-0001-000000000004"), "Голосовые сообщения", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0002-000000000004"), quiz2ID,
			"Какой стиль общения наиболее эффективен в деловой среде?", 4, []opt{
				{uuid.MustParse("d0000000-0002-0004-0001-000000000001"), "Агрессивный", false, 1},
				{uuid.MustParse("d0000000-0002-0004-0001-000000000002"), "Пассивный", false, 2},
				{uuid.MustParse("d0000000-0002-0004-0001-000000000003"), "Ассертивный — уверенный и уважительный", true, 3},
				{uuid.MustParse("d0000000-0002-0004-0001-000000000004"), "Манипулятивный", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0002-000000000005"), quiz2ID,
			"Что важно при написании деловых электронных писем?", 5, []opt{
				{uuid.MustParse("d0000000-0002-0005-0001-000000000001"), "Максимально длинные и детальные письма", false, 1},
				{uuid.MustParse("d0000000-0002-0005-0001-000000000002"), "Краткость, ясность и конкретная тема письма", true, 2},
				{uuid.MustParse("d0000000-0002-0005-0001-000000000003"), "Эмоциональный и личностный стиль", false, 3},
				{uuid.MustParse("d0000000-0002-0005-0001-000000000004"), "Отсутствие приветствия и подписи", false, 4},
			}},

		// Quiz 3: Тайм-менеджмент
		{uuid.MustParse("c0000000-0000-0000-0003-000000000001"), quiz3ID,
			"Что такое матрица Эйзенхауэра?", 1, []opt{
				{uuid.MustParse("d0000000-0003-0001-0001-000000000001"), "Метод командообразования", false, 1},
				{uuid.MustParse("d0000000-0003-0001-0001-000000000002"), "Инструмент расстановки приоритетов по срочности и важности", true, 2},
				{uuid.MustParse("d0000000-0003-0001-0001-000000000003"), "Техника принятия коллективных решений", false, 3},
				{uuid.MustParse("d0000000-0003-0001-0001-000000000004"), "Метод финансового планирования", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0003-000000000002"), quiz3ID,
			"Что означает техника Помодоро?", 2, []opt{
				{uuid.MustParse("d0000000-0003-0002-0001-000000000001"), "50 минут работы, 20 минут отдыха", false, 1},
				{uuid.MustParse("d0000000-0003-0002-0001-000000000002"), "25 минут работы, 5 минут отдыха", true, 2},
				{uuid.MustParse("d0000000-0003-0002-0001-000000000003"), "Работа без перерывов до завершения задачи", false, 3},
				{uuid.MustParse("d0000000-0003-0002-0001-000000000004"), "30 минут работы, 30 минут отдыха", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0003-000000000003"), quiz3ID,
			"Как называется закон, согласно которому работа заполняет всё отведённое на неё время?", 3, []opt{
				{uuid.MustParse("d0000000-0003-0003-0001-000000000001"), "Закон Мёрфи", false, 1},
				{uuid.MustParse("d0000000-0003-0003-0001-000000000002"), "Принцип Парето", false, 2},
				{uuid.MustParse("d0000000-0003-0003-0001-000000000003"), "Закон Паркинсона", true, 3},
				{uuid.MustParse("d0000000-0003-0003-0001-000000000004"), "Принцип Питера", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0003-000000000004"), quiz3ID,
			"Что означает принцип Парето применительно к управлению временем?", 4, []opt{
				{uuid.MustParse("d0000000-0003-0004-0001-000000000001"), "50% усилий даёт 50% результата", false, 1},
				{uuid.MustParse("d0000000-0003-0004-0001-000000000002"), "20% усилий даёт 80% результата", true, 2},
				{uuid.MustParse("d0000000-0003-0004-0001-000000000003"), "80% усилий — на планирование, 20% — на исполнение", false, 3},
				{uuid.MustParse("d0000000-0003-0004-0001-000000000004"), "10% задач занимают 90% времени", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0003-000000000005"), quiz3ID,
			"Что такое «поедание лягушки» в тайм-менеджменте?", 5, []opt{
				{uuid.MustParse("d0000000-0003-0005-0001-000000000001"), "Откладывание самых важных задач на потом", false, 1},
				{uuid.MustParse("d0000000-0003-0005-0001-000000000002"), "Выполнение самых неприятных задач первыми с утра", true, 2},
				{uuid.MustParse("d0000000-0003-0005-0001-000000000003"), "Делегирование сложных задач подчинённым", false, 3},
				{uuid.MustParse("d0000000-0003-0005-0001-000000000004"), "Работа исключительно в ночное время", false, 4},
			}},

		// Quiz 4: Лидерство
		{uuid.MustParse("c0000000-0000-0000-0004-000000000001"), quiz4ID,
			"Что такое эмоциональный интеллект лидера?", 1, []opt{
				{uuid.MustParse("d0000000-0004-0001-0001-000000000001"), "Высокий коэффициент интеллекта (IQ)", false, 1},
				{uuid.MustParse("d0000000-0004-0001-0001-000000000002"), "Способность распознавать и управлять эмоциями — своими и чужими", true, 2},
				{uuid.MustParse("d0000000-0004-0001-0001-000000000003"), "Умение быть всегда строгим и требовательным", false, 3},
				{uuid.MustParse("d0000000-0004-0001-0001-000000000004"), "Высокая техническая компетентность", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0004-000000000002"), quiz4ID,
			"Что означает ситуационное лидерство?", 2, []opt{
				{uuid.MustParse("d0000000-0004-0002-0001-000000000001"), "Применение одного стиля руководства во всех ситуациях", false, 1},
				{uuid.MustParse("d0000000-0004-0002-0001-000000000002"), "Адаптация стиля руководства к уровню зрелости сотрудника", true, 2},
				{uuid.MustParse("d0000000-0004-0002-0001-000000000003"), "Жёсткий контроль всех рабочих процессов", false, 3},
				{uuid.MustParse("d0000000-0004-0002-0001-000000000004"), "Полное делегирование всех полномочий", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0004-000000000003"), quiz4ID,
			"Что является ключевым в трансформационном лидерстве?", 3, []opt{
				{uuid.MustParse("d0000000-0004-0003-0001-000000000001"), "Система вознаграждений за результат", false, 1},
				{uuid.MustParse("d0000000-0004-0003-0001-000000000002"), "Строгий контроль и управление по отклонениям", false, 2},
				{uuid.MustParse("d0000000-0004-0003-0001-000000000003"), "Вдохновение и мотивация через общее видение и ценности", true, 3},
				{uuid.MustParse("d0000000-0004-0003-0001-000000000004"), "Чёткие инструкции для каждой задачи", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0004-000000000004"), quiz4ID,
			"Какой стиль лидерства предполагает минимальное вмешательство руководителя?", 4, []opt{
				{uuid.MustParse("d0000000-0004-0004-0001-000000000001"), "Авторитарный", false, 1},
				{uuid.MustParse("d0000000-0004-0004-0001-000000000002"), "Демократический", false, 2},
				{uuid.MustParse("d0000000-0004-0004-0001-000000000003"), "Транзакционный", false, 3},
				{uuid.MustParse("d0000000-0004-0004-0001-000000000004"), "Делегирующий (Laissez-faire)", true, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0004-000000000005"), quiz4ID,
			"Что описывает пирамида потребностей Маслоу применительно к мотивации сотрудников?", 5, []opt{
				{uuid.MustParse("d0000000-0004-0005-0001-000000000001"), "Иерархию управленческих уровней", false, 1},
				{uuid.MustParse("d0000000-0004-0005-0001-000000000002"), "Иерархию потребностей от базовых до потребности в самореализации", true, 2},
				{uuid.MustParse("d0000000-0004-0005-0001-000000000003"), "Этапы карьерного роста сотрудника", false, 3},
				{uuid.MustParse("d0000000-0004-0005-0001-000000000004"), "Модель принятия решений в условиях неопределённости", false, 4},
			}},

		// Quiz 5: Финансы
		{uuid.MustParse("c0000000-0000-0000-0005-000000000001"), quiz5ID,
			"Что означает аббревиатура ROI?", 1, []opt{
				{uuid.MustParse("d0000000-0005-0001-0001-000000000001"), "Возврат на инвестиции (Return on Investment)", true, 1},
				{uuid.MustParse("d0000000-0005-0001-0001-000000000002"), "Рентабельность активов", false, 2},
				{uuid.MustParse("d0000000-0005-0001-0001-000000000003"), "Норма прибыли", false, 3},
				{uuid.MustParse("d0000000-0005-0001-0001-000000000004"), "Коэффициент ликвидности", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0005-000000000002"), quiz5ID,
			"Что такое EBITDA?", 2, []opt{
				{uuid.MustParse("d0000000-0005-0002-0001-000000000001"), "Прибыль до вычета процентов, налогов, износа и амортизации", true, 1},
				{uuid.MustParse("d0000000-0005-0002-0001-000000000002"), "Чистая прибыль компании после всех вычетов", false, 2},
				{uuid.MustParse("d0000000-0005-0002-0001-000000000003"), "Выручка от реализации продукции", false, 3},
				{uuid.MustParse("d0000000-0005-0002-0001-000000000004"), "Собственный капитал компании", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0005-000000000003"), quiz5ID,
			"Что такое точка безубыточности?", 3, []opt{
				{uuid.MustParse("d0000000-0005-0003-0001-000000000001"), "Точка максимальной прибыли", false, 1},
				{uuid.MustParse("d0000000-0005-0003-0001-000000000002"), "Объём продаж, при котором выручка равна совокупным затратам", true, 2},
				{uuid.MustParse("d0000000-0005-0003-0001-000000000003"), "Минимально допустимая цена товара", false, 3},
				{uuid.MustParse("d0000000-0005-0003-0001-000000000004"), "Максимальный объём производства при данных ресурсах", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0005-000000000004"), quiz5ID,
			"Что такое диверсификация инвестиций?", 4, []opt{
				{uuid.MustParse("d0000000-0005-0004-0001-000000000001"), "Концентрация всех средств в одном высокодоходном активе", false, 1},
				{uuid.MustParse("d0000000-0005-0004-0001-000000000002"), "Распределение инвестиций по разным активам для снижения риска", true, 2},
				{uuid.MustParse("d0000000-0005-0004-0001-000000000003"), "Продажа всех активов и перевод в наличные", false, 3},
				{uuid.MustParse("d0000000-0005-0004-0001-000000000004"), "Получение кредита под залог активов", false, 4},
			}},
		{uuid.MustParse("c0000000-0000-0000-0005-000000000005"), quiz5ID,
			"Что показывает NPV (чистая приведённая стоимость) проекта?", 5, []opt{
				{uuid.MustParse("d0000000-0005-0005-0001-000000000001"), "Общую выручку проекта без учёта затрат", false, 1},
				{uuid.MustParse("d0000000-0005-0005-0001-000000000002"), "Разность между дисконтированными доходами и расходами проекта", true, 2},
				{uuid.MustParse("d0000000-0005-0005-0001-000000000003"), "Текущую ставку рефинансирования", false, 3},
				{uuid.MustParse("d0000000-0005-0005-0001-000000000004"), "Размер дивидендов на одну акцию", false, 4},
			}},
	}

	for _, q := range questions {
		_, err := db.Exec(
			`INSERT INTO quiz_questions (id, quiz_id, text, "order") VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING`,
			q.id, q.quizID, q.text, q.order,
		)
		if err != nil {
			return err
		}
		for _, o := range q.options {
			_, err := db.Exec(
				`INSERT INTO quiz_options (id, question_id, text, is_correct, "order") VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING`,
				o.id, q.id, o.text, o.isCorrect, o.order,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
