const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";
const REFRESH_KEY = "prepio_refresh_token"; // mobile fallback only

export type ApiError = { code: string; message: string };

type Envelope<T> = { data: T };
type ErrorEnvelope = { error: ApiError };

export class ApiClient {
  private accessToken: string | null = null;
  private refreshPromise: Promise<boolean> | null = null;

  /** setAuthTokens stores access token in memory; refresh token lives in httpOnly cookie (set by server). */
  setAuthTokens(accessToken: string | null, _refreshToken: string | null = null) {
    this.accessToken = accessToken;
    if (typeof window !== "undefined") {
      sessionStorage.removeItem(REFRESH_KEY);
      localStorage.removeItem("prepio_access_token");
    }
  }

  /** setToken clears or sets access token only (legacy). */
  setToken(token: string | null) {
    this.accessToken = token;
    if (!token && typeof window !== "undefined") {
      sessionStorage.removeItem(REFRESH_KEY);
      localStorage.removeItem("prepio_access_token");
    }
  }

  /** loadToken returns the in-memory access token. */
  loadToken() {
    return this.accessToken;
  }

  /** ensureSession bootstraps access token via refresh token on page load. */
  async ensureSession(): Promise<boolean> {
    if (this.accessToken) return true;
    if (typeof window === "undefined") return false;

    const legacy = localStorage.getItem("prepio_access_token");
    if (legacy) {
      this.accessToken = legacy;
      localStorage.removeItem("prepio_access_token");
      return true;
    }

    return this.refreshAccessToken();
  }

  private async refreshAccessToken(): Promise<boolean> {
    if (this.refreshPromise) return this.refreshPromise;

    this.refreshPromise = (async () => {
      try {
        const res = await fetch(`${API_URL}/api/v1/auth/refresh`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({ refresh_token: "" }),
        });
        const body = await res.json();
        if (!res.ok) return false;

        const data = (body as Envelope<AuthResponse>).data;
        this.accessToken = data.access_token;
        return true;
      } catch {
        return false;
      } finally {
        this.refreshPromise = null;
      }
    })();

    return this.refreshPromise;
  }

  private async request<T>(path: string, init: RequestInit = {}, retry = true): Promise<T> {
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...(init.headers as Record<string, string>),
    };
    if (this.accessToken) {
      headers.Authorization = `Bearer ${this.accessToken}`;
    }

    const res = await fetch(`${API_URL}${path}`, { ...init, headers, credentials: "include" });
    const body = await res.json();

    if (res.status === 401 && retry) {
      const refreshed = await this.refreshAccessToken();
      if (refreshed) return this.request<T>(path, init, false);
      this.setAuthTokens(null, null);
      throw new Error("session expired — please log in again");
    }

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

  getQuestionHistory(sessionId: string) {
    return this.request<HistoryEntry[]>(`/api/v1/questions/history?session_id=${sessionId}`);
  }

  getJourney() {
    return this.request<JourneyData>("/api/v1/journey");
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
  league: { tier: string; rank: number; label: string; available: boolean };
  daily_quests: {
    id: string;
    title: string;
    progress: number;
    target: number;
    completed: boolean;
    reward_xp: number;
    reward_gems: number;
    coming_soon: boolean;
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

export type HistoryEntry = {
  question_id: string;
  session_id: string;
  correct: boolean;
  score: number;
  submitted_at: string;
};

export type SubmitResponse = {
  correct: boolean;
  score: number;
  feedback: string;
  xp_awarded: number;
  gems_awarded: number;
  streak_updated: boolean;
  readiness_delta: number;
  strengths: string[];
  gaps: string[];
};

export type JourneyData = {
  world: { id: string; slug: string; name: string; description: string; theme: string };
  nodes: {
    id: string;
    label: string;
    node_type: string;
    status: string;
    question_id?: string;
    sort_order: number;
  }[];
  session_id: string;
};

export const api = new ApiClient();

export const TARGET_COMPANIES = ["google", "amazon", "meta", "uber", "atlassian"] as const;
export const EXPERIENCE_LEVELS = [
  { id: "fresher", label: "Fresher" },
  { id: "junior", label: "Junior" },
  { id: "mid", label: "Mid" },
  { id: "senior", label: "Senior" },
] as const;
