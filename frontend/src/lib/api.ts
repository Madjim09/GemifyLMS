import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export const api = axios.create({
  baseURL: API_URL,
  headers: { 'Content-Type': 'application/json' },
});

// Attach JWT token to every request if present
api.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  }
  return config;
});

export interface User {
  id: string;
  external_id: string;
  name: string;
  email: string;
  xp: number;
  level: number;
  tests_passed: number;
  avatar_url: string;
  created_at: string;
  updated_at: string;
}

export interface AuthOutput {
  user: User;
  token: string;
}

export interface Achievement {
  id: string;
  name: string;
  description: string;
  icon_url: string;
  condition_type: string;
  condition_value: number;
  xp_reward: number;
}

export interface UserAchievement {
  id: string;
  user_id: string;
  achievement_id: string;
  achievement: Achievement;
  earned_at: string;
}

export interface Skill {
  id: string;
  name: string;
  description: string;
  category: string;
  required_xp: number;
  icon_url: string;
  order: number;
  children?: Skill[];
}

export interface UserSkill {
  id: string;
  user_id: string;
  skill_id: string;
  skill: Skill;
  status: 'locked' | 'available' | 'in_progress' | 'mastered';
  progress: number;
}

export interface LeaderboardEntry {
  rank: number;
  user: User;
  xp: number;
  level: number;
}

export interface QuizOption {
  id: string;
  question_id: string;
  text: string;
  is_correct?: boolean;
  order: number;
}

export interface QuizQuestion {
  id: string;
  quiz_id: string;
  text: string;
  order: number;
  options: QuizOption[];
}

export interface Quiz {
  id: string;
  title: string;
  description: string;
  difficulty: 'easy' | 'medium' | 'hard';
  xp_reward: number;
  bonus_xp: number;
  pass_threshold: number;
  questions?: QuizQuestion[];
  created_at: string;
}

export interface UserQuizResult {
  id: string;
  user_id: string;
  quiz_id: string;
  quiz?: Quiz;
  score: number;
  xp_earned: number;
  is_perfect: boolean;
  completed_at: string;
}

export interface SubmitQuizOutput {
  score: number;
  passed: boolean;
  is_perfect: boolean;
  xp_earned: number;
  result: UserQuizResult;
}

export const LEVEL_THRESHOLDS = [0, 500, 1500, 3000, 5500, 9000, 14000, 21000, 30000, 42000];

export function getXPProgress(xp: number, level: number): number {
  const current = LEVEL_THRESHOLDS[level - 1] || 0;
  const next = LEVEL_THRESHOLDS[level] || current;
  if (next === current) return 100;
  return Math.min(100, ((xp - current) / (next - current)) * 100);
}

export function getXPForNextLevel(level: number): number {
  return LEVEL_THRESHOLDS[level] || -1;
}

export const DIFFICULTY_LABELS: Record<string, string> = {
  easy: 'Лёгкий',
  medium: 'Средний',
  hard: 'Сложный',
};

export const DIFFICULTY_COLORS: Record<string, string> = {
  easy: 'bg-green-900 text-green-400',
  medium: 'bg-yellow-900 text-yellow-400',
  hard: 'bg-red-900 text-red-400',
};
