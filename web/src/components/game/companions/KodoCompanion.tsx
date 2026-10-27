"use client";

type Reaction = "idle" | "correct" | "wrong";

export function KodoCompanion({
  reaction = "idle",
  className = "",
}: {
  reaction?: Reaction;
  className?: string;
}) {
  const isCorrect = reaction === "correct";
  const isWrong = reaction === "wrong";

  return (
    <svg
      viewBox="0 0 160 200"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      className={`select-none overflow-visible ${className}`}
      aria-label="Kodo the Frog in knitted scarf"
    >
      <defs>
        <radialGradient id="kodo-glow" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#4ADE80" stopOpacity="0.45" />
          <stop offset="100%" stopColor="#4ADE80" stopOpacity="0" />
        </radialGradient>
        <linearGradient id="kodo-skin" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stopColor="#86EFAC" />
          <stop offset="60%" stopColor="#4ADE80" />
          <stop offset="100%" stopColor="#22C55E" />
        </linearGradient>
        <linearGradient id="kodo-scarf" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#F97316" />
          <stop offset="100%" stopColor="#EA580C" />
        </linearGradient>
        <linearGradient id="kodo-vest" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stopColor="#1E293B" />
          <stop offset="100%" stopColor="#0F172A" />
        </linearGradient>
      </defs>

      {/* Ambient background glow */}
      <circle cx="80" cy="110" r="75" fill="url(#kodo-glow)" />

      {/* Shadow */}
      <ellipse cx="80" cy="186" rx="46" ry="9" fill="#0D0F17" fillOpacity="0.6" />

      {/* Feet & Webbed Boots */}
      <g>
        <g transform="translate(48, 168)">
          <path d="M4 14 C4 8, 12 6, 22 6 C28 6, 32 10, 32 16 L2 16 Z" fill="#15803D" />
          <rect x="0" y="14" width="34" height="5" rx="2.5" fill="#FBBF24" />
        </g>
        <g transform="translate(80, 168)">
          <path d="M4 14 C4 8, 12 6, 22 6 C28 6, 32 10, 32 16 L2 16 Z" fill="#15803D" />
          <rect x="0" y="14" width="34" height="5" rx="2.5" fill="#FBBF24" />
        </g>
      </g>

      {/* Torso & Scholar Vest */}
      <g>
        <path
          d="M48 114 C46 150, 50 166, 80 166 C110 166, 114 150, 112 114 C112 100, 48 100, 48 114 Z"
          fill="url(#kodo-vest)"
        />
        {/* Diamond argyle sweater pattern */}
        <polygon points="80,120 88,132 80,144 72,132" fill="#334155" />
        <polygon points="80,144 88,156 80,168 72,156" fill="#334155" />
      </g>

      {/* Knitted Scarf with Fringes */}
      <g>
        {/* Scarf Coil */}
        <path
          d="M44 94 C44 90, 116 90, 116 94 C118 108, 114 114, 80 114 C46 114, 42 108, 44 94 Z"
          fill="url(#kodo-scarf)"
          stroke="#C2410C"
          strokeWidth="2"
        />
        {/* Scarf Tail Draping Down */}
        <path d="M92 112 L96 148 L108 148 L104 112 Z" fill="url(#kodo-scarf)" stroke="#C2410C" strokeWidth="1.5" />
        {/* Scarf Fringes */}
        <path d="M96 148 L96 154 M100 148 L100 154 M104 148 L104 154 M108 148 L108 154" stroke="#FED7AA" strokeWidth="2" />
      </g>

      {/* Head Group (Frog Eye Domes, Glasses, Cute Frog Smile) */}
      <g className={isWrong ? "animate-pulse" : ""}>
        {/* Frog Left Eye Dome */}
        <circle cx="58" cy="38" r="18" fill="url(#kodo-skin)" />
        {/* Frog Right Eye Dome */}
        <circle cx="102" cy="38" r="18" fill="url(#kodo-skin)" />

        {/* Head Main Shape */}
        <ellipse cx="80" cy="62" rx="38" ry="32" fill="url(#kodo-skin)" />

        {/* Blush Marks */}
        <ellipse cx="52" cy="66" rx="5" ry="3" fill="#F472B6" fillOpacity="0.5" />
        <ellipse cx="108" cy="66" rx="5" ry="3" fill="#F472B6" fillOpacity="0.5" />

        {/* Nostrils */}
        <circle cx="76" cy="58" r="1.5" fill="#14532D" />
        <circle cx="84" cy="58" r="1.5" fill="#14532D" />

        {/* Frog Mouth */}
        {isCorrect ? (
          <path d="M68 68 Q80 82 92 68" stroke="#14532D" strokeWidth="2.5" strokeLinecap="round" fill="none" />
        ) : isWrong ? (
          <path d="M70 72 Q80 66 90 72" stroke="#14532D" strokeWidth="2.5" strokeLinecap="round" fill="none" />
        ) : (
          <path d="M68 68 Q80 76 92 68" stroke="#14532D" strokeWidth="2.5" strokeLinecap="round" fill="none" />
        )}

        {/* Big Frog Eyes */}
        {isCorrect ? (
          // Joyful curved frog eyes
          <g stroke="#14532D" strokeWidth="3" strokeLinecap="round" fill="none">
            <path d="M52 38 Q58 32 64 38" />
            <path d="M96 38 Q102 32 108 38" />
          </g>
        ) : (
          // Big lively glossy eyes
          <g>
            <circle cx="58" cy="38" r="10" fill="#14532D" />
            <circle cx="61" cy="35" r="3" fill="#FFFFFF" />
            <circle cx="102" cy="38" r="10" fill="#14532D" />
            <circle cx="105" cy="35" r="3" fill="#FFFFFF" />
          </g>
        )}

        {/* Round Scholar Wire Glasses */}
        <g stroke="#F59E0B" strokeWidth="2.5" fill="none">
          <circle cx="58" cy="46" r="16" />
          <circle cx="102" cy="46" r="16" />
          <path d="M74 46 L86 46" stroke="#F59E0B" strokeWidth="3" />
        </g>
      </g>

      {/* Arms / Webbed Hands */}
      {isCorrect ? (
        // Celebrating hands up
        <g>
          <path d="M48 106 C28 88, 16 62, 28 50 C36 42, 48 62, 56 92 Z" fill="#22C55E" />
          <path d="M112 106 C132 88, 144 62, 132 50 C124 42, 112 62, 104 92 Z" fill="#22C55E" />
          <g fill="#F59E0B">
            <polygon points="20,26 23,20 30,18 24,14 25,7 19,11 13,8 15,15 10,19 17,21" />
            <polygon points="140,26 143,20 150,18 144,14 145,7 139,11 133,8 135,15 130,19 137,21" />
          </g>
        </g>
      ) : isWrong ? (
        // Webbed hand adjusting glasses
        <g>
          <path d="M48 106 C36 120, 38 140, 50 144 C58 146, 58 132, 54 112 Z" fill="#22C55E" />
          <path d="M112 106 C124 95, 126 72, 114 62 C106 56, 102 68, 100 88 Z" fill="#22C55E" />
          <circle cx="114" cy="62" r="6" fill="#15803D" />
          {/* Sweat drop */}
          <path d="M124 42 C124 42, 130 50, 130 55 C130 58, 127 61, 124 61 C121 61, 118 58, 118 55 C118 50, 124 42, 124 42 Z" fill="#38BDF8" />
        </g>
      ) : (
        // Calm resting posture
        <g>
          <path d="M46 106 C36 122, 40 142, 54 144 C62 144, 60 130, 52 110 Z" fill="#22C55E" />
          <path d="M114 106 C124 122, 120 142, 106 144 C98 144, 100 130, 108 110 Z" fill="#22C55E" />
        </g>
      )}
    </svg>
  );
}
