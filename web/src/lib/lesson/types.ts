/** Wire types for the lesson runtime API (GET /path, attempts, answers, completion). */

export type StepType = "intro" | "mcq" | "true_false" | "fill_blank" | "arrange" | "prose";

export type Beat = { text: string; emphasis?: string; visual?: string; duration_ms: number };

/** ClientStep is what the server sends: never any answer, explanation, or rubric. */
export type ClientStep = {
  id: string;
  type: StepType;
  position: number;
  intro?: { beats: Beat[]; media_url?: string };
  mcq?: { prompt: string; options: string[] };
  true_false?: { statement: string };
  fill_blank?: { code: string; blanks: number; bank: string[] };
  arrange?: { prompt: string; items: string[] };
  prose?: { prompt: string; min_chars: number };
};

export type LessonInfo = {
  id: string;
  slug: string;
  title: string;
  kind: "lesson" | "boss";
  difficulty: "easy" | "medium" | "hard";
  est_minutes: number;
};

export type StepProgress = { step_id: string; tries: number; done: boolean };

export type AttemptData = {
  attempt_id: string;
  resumed: boolean;
  lesson: LessonInfo;
  steps: ClientStep[];
  progress: StepProgress[];
};

/** Answer is the learner's submission; which field is used depends on the step type. */
export type Answer = {
  choice?: number;
  value?: boolean;
  blanks?: string[];
  order?: number[];
  text?: string;
};

export type AnswerResult = {
  correct: boolean;
  try: number;
  explanation?: string;
  why_not?: string;
  correct_answer?: Answer;
  feedback?: string;
  score?: number;
  strengths?: string[];
  gaps?: string[];
  step_done: boolean;
  all_done: boolean;
  replayed: boolean;
};

export type MasteryChange = {
  skill_slug: string;
  skill_name: string;
  topic_slug?: string;
  topic_name?: string;
  before: number;
  after: number;
  delta: number;
  accuracy: number;
};

export type Rewards = {
  attempt_id: string;
  first_completion: boolean;
  xp_awarded: number;
  gems_awarded: number;
  mastery_changes: MasteryChange[];
};

export type UnlockedNode = { id: string; slug: string; label: string; lesson_id: string };

export type CompletionData = {
  attempt_id: string;
  lesson: LessonInfo;
  takeaways: string[];
  graded_steps: number;
  first_try_correct: number;
  total_tries: number;
  accuracy: number;
  node_id: string;
  unlocked_nodes: UnlockedNode[];
  rewards: Rewards | null;
  rewards_pending: boolean;
};

export type NodeStatus = "locked" | "current" | "available" | "done";

export type PathNode = {
  id: string;
  slug: string;
  label: string;
  node_type: "lesson" | "boss";
  status: NodeStatus;
  unlock_hint?: string;
  in_progress: boolean;
  lesson_id: string;
  lesson_slug: string;
  title: string;
  takeaways: string[];
  kind: "lesson" | "boss";
  difficulty: "easy" | "medium" | "hard";
  est_minutes: number;
  xp_preview: number;
};

export type PathWorld = {
  id: string;
  slug: string;
  name: string;
  description: string;
  theme: string;
  /** topic is the topic slug the world builds; focused worlds are the learner's focus topics. */
  topic?: string;
  focused: boolean;
  nodes: PathNode[];
};

export type PathData = { worlds: PathWorld[] };
