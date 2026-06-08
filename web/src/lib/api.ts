const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export type ApiError = { code: string; message: string };

type Envelope<T> = { data: T };
type ErrorEnvelope = { error: ApiError };

export class ApiClient {
  private accessToken: string | null = null;

  setToken(token: string | null) {
    this.accessToken = token;
    if (typeof window !== "undefined") {
      if (token) {
        localStorage.setItem("prepio_access_token", token);
      } else {
        localStorage.removeItem("prepio_access_token");
      }
    }
  }

  loadToken() {
    if (typeof window !== "undefined") {
      this.accessToken = localStorage.getItem("prepio_access_token");
    }
    return this.accessToken;
  }

  private async request<T>(path: string, init: RequestInit = {}): Promise<T> {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...(init.headers as Record<string, string>),
    };
    if (this.accessToken) {
      headers.Authorization = `Bearer ${this.accessToken}`;
    }

    const res = await fetch(`${API_URL}${path}`, { ...init, headers });
    const body = await res.json();

    if (!res.ok) {
      const err = (body as ErrorEnvelope).error;
      throw new Error(err?.message ?? "request failed");
    }
    return (body as Envelope<T>).data;
  }

  register(email: string, username: string, password: string) {
    return this.request<AuthResponse>("/api/v1/auth/register", {
      method: "POST",
      body: JSON.stringify({ email, username, password }),
    });
  }

  login(email: string, password: string) {
    return this.request<AuthResponse>("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify({ email, password }),
    });
  }

  getProfile() {
    return this.request<Profile>("/api/v1/users/profile");
  }

  getCompanions() {
    return this.request<Companion[]>("/api/v1/companions");
  }

  completeOnboarding(targetCompanies: string[], experienceLevel: string, companionId: string) {
    return this.request<Profile>("/api/v1/users/onboarding", {
      method: "POST",
      body: JSON.stringify({
        target_companies: targetCompanies,
        experience_level: experienceLevel,
        companion_id: companionId,
      }),
    });
  }

  getDashboardHome() {
    return this.request<DashboardHome>("/api/v1/dashboard/home");
  }

  getDailyPaper() {
    return this.request<DailyPaper>("/api/v1/questions/daily");
  }

  submitAnswer(questionId: string, sessionId: string, answer: string) {
    return this.request<SubmitResponse>(`/api/v1/questions/${questionId}/submit`, {
      method: "POST",
      body: JSON.stringify({
        session_id: sessionId,
        answer,
        time_spent_seconds: 60,
      }),
    });
  }
}

export type AuthResponse = {
  access_token: string;
  refresh_token: string;
  user: { id: string; username: string; email: string };
};

export type Companion = {
  id: string;
  name: string;
  species: string;
};

export type Profile = {
  id: string;
  email: string;
  username: string;
  experience_level?: string;
  onboarding_completed: boolean;
  target_companies: string[];
  companion?: Companion;
};

export type DashboardHome = {
  streak: {
    current_streak: number;
    longest_streak: number;
    freeze_count: number;
    streak_active_today: boolean;
  };
  progress: {
    total_xp: number;
    current_level: number;
    gem_balance: number;
    xp_to_next_level: number;
  };
  companion: Companion;
  readiness: { company: string; score: number }[];
  league: { tier: string; rank: number; label: string };
  daily_quests: {
    id: string;
    title: string;
    progress: number;
    target: number;
    completed: boolean;
    reward_xp: number;
    reward_gems: number;
  }[];
  companion_message: string;
  onboarding_needed: boolean;
};

export type DailyPaper = {
  session_id: string;
  date: string;
  questions: Question[];
  minimum_to_streak: number;
};

export type Question = {
  id: string;
  body: string;
  round_type: string;
  difficulty: string;
  company_tags: string[];
  is_weekend: boolean;
};

export type SubmitResponse = {
  correct: boolean;
  feedback: string;
  xp_awarded: number;
  gems_awarded: number;
  streak_updated: boolean;
};

export const api = new ApiClient();

export const TARGET_COMPANIES = ["google", "amazon", "meta", "uber", "atlassian"] as const;
export const EXPERIENCE_LEVELS = [
  { id: "fresher", label: "Fresher" },
  { id: "junior", label: "Junior" },
  { id: "mid", label: "Mid" },
  { id: "senior", label: "Senior" },
] as const;
