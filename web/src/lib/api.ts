import type {
  Answer,
  AnswerResult,
  AttemptData,
  CompletionData,
  PathData,
} from "@/lib/lesson/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export type ApiError = { code: string; message: string };

type Envelope<T> = { data: T };
type ErrorEnvelope = { error: ApiError };

/** ApiRequestError carries the HTTP status and the server's machine-readable code. */
export class ApiRequestError extends Error {
  constructor(
    message: string,
    readonly status: number,
    readonly code: string,
  ) {
    super(message);
    this.name = "ApiRequestError";
  }
}

export class ApiClient {
  private accessToken: string | null = null;
  private refreshPromise: Promise<boolean> | null = null;

  /** setAuthTokens stores the access token in memory; the refresh token lives in an httpOnly cookie set by the server. */
  setAuthTokens(accessToken: string | null) {
    this.accessToken = accessToken;
    if (typeof window !== "undefined") {
      localStorage.removeItem("prepio_access_token");
    }
  }

  /** logout revokes the session server-side (clearing the refresh cookie) and forgets the access token. */
  async logout(): Promise<void> {
    try {
      if (this.accessToken) {
        await this.request<unknown>("/api/v1/auth/logout", { method: "POST" }, false);
      }
    } catch {
      // Best effort: the local session is dropped regardless.
    } finally {
      this.setAuthTokens(null);
    }
  }

  /** ensureSession bootstraps the access token from the refresh cookie on page load. */
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
          credentials: "include",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ refresh_token: "" }),
        });
        if (!res.ok) return false;
        const body = parseBody(await res.text()) as Envelope<AuthResponse> | null;
        const token = body?.data?.access_token;
        if (!token) return false;
        this.accessToken = token;
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

    if (res.status === 401 && retry) {
      const refreshed = await this.refreshAccessToken();
      if (refreshed) return this.request<T>(path, init, false);
      this.setAuthTokens(null);
      throw new ApiRequestError("session expired — please log in again", 401, "unauthorized");
    }

    let text = "";
    try {
      text = await res.text();
    } catch {
      // An unreadable body is treated like an empty one.
    }
    const body = parseBody(text);

    if (!res.ok) {
      const err = (body as Partial<ErrorEnvelope> | null)?.error;
      const fallback = res.statusText
        ? `request failed (${res.status} ${res.statusText})`
        : `request failed (${res.status})`;
      throw new ApiRequestError(err?.message ?? fallback, res.status, err?.code ?? "unknown");
    }
    if (body === null || typeof body !== "object" || !("data" in body)) {
      throw new ApiRequestError("unexpected response from server", res.status, "bad_response");
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

  /** getTopics returns the topic catalog (what a learner can choose to focus on). */
  getTopics() {
    return this.request<TopicInfo[]>("/api/v1/topics");
  }

  completeOnboarding(experienceLevel: string, companionId: string, focusTopics: string[]) {
    return this.request<Profile>("/api/v1/users/onboarding", {
      method: "POST",
      body: JSON.stringify({
        experience_level: experienceLevel,
        companion_id: companionId,
        focus_topics: focusTopics,
      }),
    });
  }

  getDashboardHome() {
    return this.request<DashboardHome>("/api/v1/dashboard/home");
  }

  /** getLeague returns this week's league: tier, standings with public cards, last result. */
  getLeague() {
    return this.request<League>("/api/v1/league");
  }

  /** getPath returns worlds and nodes with status, previews, and unlock hints; focus topics' worlds come first. */
  getPath(focus: string[] = []) {
    const q = focus.length > 0 ? `?focus=${encodeURIComponent(focus.join(","))}` : "";
    return this.request<PathData>(`/api/v1/path${q}`);
  }

  /** startAttempt starts a lesson attempt, or resumes the one in progress. */
  startAttempt(lessonId: string) {
    return this.request<AttemptData>(`/api/v1/lessons/${encodeURIComponent(lessonId)}/attempts`, {
      method: "POST",
    });
  }

  /** answerStep sends one try; the server grades it. A repeated try returns the stored result. */
  answerStep(attemptId: string, stepId: string, tryNo: number, answer: Answer) {
    return this.request<AnswerResult>(
      `/api/v1/attempts/${encodeURIComponent(attemptId)}/steps/${encodeURIComponent(stepId)}/answer`,
      { method: "POST", body: JSON.stringify({ try: tryNo, answer }) },
    );
  }

  /** completeAttempt finishes the attempt; the response includes Progress rewards. */
  completeAttempt(attemptId: string) {
    return this.request<CompletionData>(`/api/v1/attempts/${encodeURIComponent(attemptId)}/complete`, {
      method: "POST",
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
  companion?: Companion;
  focus_topics?: string[];
};

export type TopicInfo = { slug: string; name: string; description: string };

/** TopicCard is the learner readiness in one topic; mastery is null until a skill is started. */
export type TopicCard = {
  slug: string;
  name: string;
  description: string;
  mastery: number | null;
  skills_started: number;
  skills_total: number;
  focused: boolean;
};

export type NextLesson = {
  lesson_id: string;
  title: string;
  node_label: string;
  world_name: string;
  est_minutes: number;
  xp_preview: number;
  in_progress: boolean;
  kind: "lesson" | "boss";
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
  topics: TopicCard[];
  focus_topics: string[];
  next_lesson: NextLesson | null;
  /** league is this week's league at a glance; rank and XP are only set once joined. */
  league: {
    tier_index: number;
    tier_slug: string;
    tier_name: string;
    joined: boolean;
    rank: number;
    cohort_size: number;
    weekly_xp: number;
    zone: LeagueZone;
    ends_at: string;
  };
  companion_message: string;
  onboarding_needed: boolean;
};

export type LeagueTier = { index: number; slug: string; name: string };

/** LeagueZone is what a rank would earn if the week ended now (decided by the server). */
export type LeagueZone = "promote" | "stay" | "demote";

export type LeagueStanding = {
  rank: number;
  user_id: string;
  weekly_xp: number;
  zone: LeagueZone;
  is_me: boolean;
  username: string;
  companion_name: string;
  companion_species: string;
};

export type LeagueResult = {
  week_start: string;
  rank: number;
  cohort_size: number;
  outcome: "promoted" | "stayed" | "demoted";
  from_tier: LeagueTier;
  to_tier: LeagueTier;
};

/** League is this week's league. Until the learner earns XP this week, joined is false and standings are empty. */
export type League = {
  week_start: string;
  ends_at: string;
  joined: boolean;
  tier: LeagueTier;
  tiers: LeagueTier[];
  my_rank: number;
  standings: LeagueStanding[];
  last_result: LeagueResult | null;
};

/** parseBody parses a JSON body; empty or non-JSON bodies (e.g. an HTML proxy error page) yield null. */
function parseBody(text: string): unknown {
  if (!text) return null;
  try {
    return JSON.parse(text);
  } catch {
    return null;
  }
}

export const api = new ApiClient();

export const EXPERIENCE_LEVELS = [
  { id: "fresher", label: "Fresher" },
  { id: "junior", label: "Junior" },
  { id: "mid", label: "Mid" },
  { id: "senior", label: "Senior" },
] as const;
