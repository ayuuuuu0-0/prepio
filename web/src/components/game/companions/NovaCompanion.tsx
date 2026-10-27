"use client";

type Reaction = "idle" | "correct" | "wrong";

export function NovaCompanion({
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
      aria-label="Nova the Hedgehog in cyber vest"
    >
      <defs>
        <radialGradient id="nova-glow" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#7B5CFF" stopOpacity="0.45" />
          <stop offset="100%" stopColor="#7B5CFF" stopOpacity="0" />
        </radialGradient>
        <linearGradient id="nova-quills" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stopColor="#A78BFA" />
          <stop offset="100%" stopColor="#7C3AED" />
        </linearGradient>
        <linearGradient id="nova-vest" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#0F172A" />
          <stop offset="100%" stopColor="#1E1B4B" />
        </linearGradient>
      </defs>

      {/* Ambient background glow */}
      <circle cx="80" cy="110" r="75" fill="url(#nova-glow)" />

      {/* Shadow */}
      <ellipse cx="80" cy="186" rx="46" ry="9" fill="#0D0F17" fillOpacity="0.6" />

      {/* Cyber Hedgehog Spines / Crest */}
      <g>
        <polygon points="46,30 36,12 56,24" fill="#6D28D9" />
        <polygon points="64,22 60,4 76,18" fill="#7C3AED" />
        <polygon points="86,18 100,4 96,22" fill="#7C3AED" />
        <polygon points="104,24 124,12 114,30" fill="#6D28D9" />
      </g>

      {/* Feet & High-Tech Sneakers */}
      <g>
        <g transform="translate(48, 168)">
          <path d="M4 14 C4 8, 12 6, 22 6 C28 6, 32 10, 32 16 L2 16 Z" fill="#0F172A" />
          <rect x="0" y="14" width="34" height="5" rx="2.5" fill="#22D3EE" />
        </g>
        <g transform="translate(80, 168)">
          <path d="M4 14 C4 8, 12 6, 22 6 C28 6, 32 10, 32 16 L2 16 Z" fill="#0F172A" />
          <rect x="0" y="14" width="34" height="5" rx="2.5" fill="#22D3EE" />
        </g>
      </g>

      {/* Cyber Vest & Body */}
      <g>
        <path
          d="M44 112 C42 152, 46 168, 80 168 C114 168, 118 152, 116 112 C116 96, 44 96, 44 112 Z"
          fill="url(#nova-vest)"
        />
        {/* Glowing Circuit Lines */}
        <path d="M64 112 L64 148 L76 156" stroke="#22D3EE" strokeWidth="2.5" strokeLinecap="round" fill="none" />
        <circle cx="76" cy="156" r="3" fill="#22D3EE" />
        <path d="M96 112 L96 148 L84 156" stroke="#F43F5E" strokeWidth="2.5" strokeLinecap="round" fill="none" />
        <circle cx="84" cy="156" r="3" fill="#F43F5E" />
      </g>

      {/* Head Group (Cute Snout, Cyber Visor) */}
      <g className={isWrong ? "animate-pulse" : ""}>
        {/* Head Shape */}
        <circle cx="80" cy="66" r="36" fill="url(#nova-quills)" />

        {/* Hedgehog Cute Snout */}
        <ellipse cx="80" cy="74" rx="11" ry="8" fill="#5B21B6" />
        <circle cx="80" cy="73" r="2.5" fill="#0F172A" />

        {/* Cyber Visor / Goggles */}
        <rect x="46" y="52" width="68" height="18" rx="8" fill="#0F172A" stroke="#22D3EE" strokeWidth="2" />
        {/* Visor Glare & Neon Grid */}
        <path d="M52 61 L108 61" stroke="#22D3EE" strokeWidth="2" strokeDasharray="6 3" />
        {isCorrect ? (
          <path d="M58 61 L64 56 L70 61 M90 61 L96 56 L102 61" stroke="#34D399" strokeWidth="2.5" strokeLinecap="round" />
        ) : isWrong ? (
          <path d="M60 58 L68 64 M68 58 L60 64 M92 58 L100 64 M100 58 L92 64" stroke="#F43F5E" strokeWidth="2.5" strokeLinecap="round" />
        ) : (
          <circle cx="64" cy="61" r="3" fill="#22D3EE" />
        )}
      </g>

      {/* Arms */}
      {isCorrect ? (
        <g>
          <path d="M46 108 C30 90, 18 64, 28 50 C36 42, 46 64, 54 94 Z" fill="#5B21B6" />
          <circle cx="26" cy="48" r="8" fill="#7C3AED" />
          <path d="M114 108 C130 90, 142 64, 132 50 C124 42, 114 64, 106 94 Z" fill="#5B21B6" />
          <circle cx="134" cy="48" r="8" fill="#7C3AED" />
        </g>
      ) : isWrong ? (
        <g>
          <path d="M46 108 C34 122, 38 142, 48 144 C56 146, 58 132, 54 112 Z" fill="#5B21B6" />
          <circle cx="48" cy="144" r="7" fill="#7C3AED" />
          <path d="M114 108 C126 95, 128 72, 118 62 C110 56, 106 68, 104 88 Z" fill="#5B21B6" />
          <circle cx="118" cy="62" r="7" fill="#7C3AED" />
          <path d="M126 44 C126 44, 132 52, 132 57 C132 60, 129 63, 126 63 C123 63, 120 60, 120 57 C120 52, 126 44, 126 44 Z" fill="#38BDF8" />
        </g>
      ) : (
        <g>
          <path d="M45 108 C36 122, 40 142, 54 144 C62 144, 60 130, 52 110 Z" fill="#5B21B6" />
          <circle cx="54" cy="144" r="7" fill="#7C3AED" />
          <path d="M115 108 C124 122, 120 142, 106 144 C98 144, 100 130, 108 110 Z" fill="#5B21B6" />
          <circle cx="106" cy="144" r="7" fill="#7C3AED" />
        </g>
      )}
    </svg>
  );
}
