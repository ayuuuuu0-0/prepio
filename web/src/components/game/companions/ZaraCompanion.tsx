"use client";

type Reaction = "idle" | "correct" | "wrong";

export function ZaraCompanion({
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
      aria-label="Zara the Orange Cat in bomber jacket and headphones"
    >
      <defs>
        <radialGradient id="zara-glow" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#EA580C" stopOpacity="0.45" />
          <stop offset="100%" stopColor="#EA580C" stopOpacity="0" />
        </radialGradient>
        <linearGradient id="zara-fur" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stopColor="#FB923C" />
          <stop offset="50%" stopColor="#F97316" />
          <stop offset="100%" stopColor="#EA580C" />
        </linearGradient>
        <linearGradient id="zara-bomber" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#334155" />
          <stop offset="60%" stopColor="#1E293B" />
          <stop offset="100%" stopColor="#0F172A" />
        </linearGradient>
        <linearGradient id="zara-headphones" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#06B6D4" />
          <stop offset="100%" stopColor="#0891B2" />
        </linearGradient>
      </defs>

      {/* Ambient background glow */}
      <circle cx="80" cy="110" r="75" fill="url(#zara-glow)" />

      {/* Shadow */}
      <ellipse cx="80" cy="186" rx="46" ry="9" fill="#0D0F17" fillOpacity="0.6" />

      {/* Orange Cat Tail Swishing */}
      <path
        d="M108 144 C128 148, 148 130, 142 108 C138 92, 126 94, 128 106 C130 118, 118 132, 104 136 Z"
        fill="url(#zara-fur)"
      />
      {/* Orange Cat Tail Ring */}
      <path d="M136 108 C132 114, 124 118, 120 120" stroke="#C2410C" strokeWidth="4" strokeLinecap="round" />

      {/* Feet & Boots */}
      <g>
        {/* Left Foot */}
        <g transform="translate(48, 168)">
          <path d="M4 14 C4 8, 12 6, 22 6 C28 6, 32 10, 32 16 L2 16 Z" fill="#0F172A" />
          <rect x="0" y="14" width="34" height="5" rx="2.5" fill="#F97316" />
        </g>
        {/* Right Foot */}
        <g transform="translate(80, 168)">
          <path d="M4 14 C4 8, 12 6, 22 6 C28 6, 32 10, 32 16 L2 16 Z" fill="#0F172A" />
          <rect x="0" y="14" width="34" height="5" rx="2.5" fill="#F97316" />
        </g>
      </g>

      {/* Bomber Jacket Body */}
      <g>
        <path
          d="M44 112 C42 152, 46 168, 80 168 C114 168, 118 152, 116 112 C116 96, 44 96, 44 112 Z"
          fill="url(#zara-bomber)"
        />
        {/* Orange Ribbed Collar */}
        <path d="M58 98 Q80 110 102 98" stroke="#F97316" strokeWidth="4" fill="none" strokeLinecap="round" />
        {/* Orange Ribbed Waistband */}
        <rect x="47" y="164" width="66" height="6" rx="3" fill="#EA580C" />
        {/* Front Metal Zipper */}
        <path d="M80 106 L80 164" stroke="#94A3B8" strokeWidth="2.5" strokeDasharray="3 2" />
        {/* Developer Badge Patch on Chest */}
        <rect x="54" y="118" width="14" height="10" rx="3" fill="#10B981" />
        <path d="M57 123 L61 123 M63 123 L65 123" stroke="#FFFFFF" strokeWidth="1.5" strokeLinecap="round" />
        {/* Cyber Patch on Right Chest */}
        <circle cx="102" cy="123" r="5" fill="#8B5CF6" />
      </g>

      {/* Head Group (Ears, Whiskers, Face) */}
      <g className={isWrong ? "animate-pulse" : ""}>
        {/* Cat Left Pointed Ear */}
        <g transform="translate(40, 26)">
          <polygon points="12,0 0,36 26,36" fill="#EA580C" />
          <polygon points="12,8 5,32 20,32" fill="#F472B6" />
        </g>
        {/* Cat Right Pointed Ear */}
        <g transform="translate(94, 26)">
          <polygon points="14,0 0,36 26,36" fill="#EA580C" />
          <polygon points="14,8 6,32 21,32" fill="#F472B6" />
        </g>

        {/* Head Main Shape */}
        <ellipse cx="80" cy="68" rx="36" ry="32" fill="url(#zara-fur)" />

        {/* Orange Cat Tabby Forehead Stripes */}
        <path d="M80 44 L80 52" stroke="#C2410C" strokeWidth="3" strokeLinecap="round" />
        <path d="M72 46 L74 52" stroke="#C2410C" strokeWidth="2.5" strokeLinecap="round" />
        <path d="M88 46 L86 52" stroke="#C2410C" strokeWidth="2.5" strokeLinecap="round" />

        {/* White Muzzle */}
        <ellipse cx="80" cy="76" rx="14" ry="10" fill="#FFFBEB" />
        {/* Pink Nose */}
        <polygon points="80,75 76,71 84,71" fill="#F43F5E" />

        {/* Whiskers */}
        <g stroke="#78350F" strokeWidth="1.8" strokeLinecap="round">
          {/* Left */}
          <path d="M62 74 L42 70" />
          <path d="M62 77 L40 78" />
          {/* Right */}
          <path d="M98 74 L118 70" />
          <path d="M98 77 L120 78" />
        </g>

        {/* Mouth */}
        {isCorrect ? (
          <path d="M74 80 Q80 88 86 80" stroke="#78350F" strokeWidth="2.2" strokeLinecap="round" fill="none" />
        ) : isWrong ? (
          <path d="M75 83 Q80 79 85 83" stroke="#78350F" strokeWidth="2.2" strokeLinecap="round" fill="none" />
        ) : (
          <path d="M74 78 Q77 82 80 78 Q83 82 86 78" stroke="#78350F" strokeWidth="2" strokeLinecap="round" fill="none" />
        )}

        {/* Cat Eyes */}
        {isCorrect ? (
          // Joyful cat curve ^^
          <g stroke="#451A03" strokeWidth="3" strokeLinecap="round" fill="none">
            <path d="M58 62 Q64 56 70 62" />
            <path d="M90 62 Q96 56 102 62" />
          </g>
        ) : (
          // Almond feline eyes with green glow
          <g>
            <ellipse cx="64" cy="62" rx="6" ry="7" fill="#10B981" />
            <ellipse cx="64" cy="62" rx="2" ry="6" fill="#064E3B" />
            <circle cx="66" cy="59" r="1.8" fill="#FFFFFF" />

            <ellipse cx="96" cy="62" rx="6" ry="7" fill="#10B981" />
            <ellipse cx="96" cy="62" rx="2" ry="6" fill="#064E3B" />
            <circle cx="98" cy="59" r="1.8" fill="#FFFFFF" />
          </g>
        )}

        {/* Studio Over-Ear Headphones */}
        <g>
          {/* Padded Headband Arch */}
          <path
            d="M40 68 C38 28, 122 28, 120 68"
            stroke="#0F172A"
            strokeWidth="7"
            strokeLinecap="round"
            fill="none"
          />
          <path
            d="M48 48 C60 38, 100 38, 112 48"
            stroke="#06B6D4"
            strokeWidth="3"
            strokeLinecap="round"
            fill="none"
          />
          {/* Left Earcup */}
          <rect x="34" y="56" width="12" height="24" rx="6" fill="url(#zara-headphones)" stroke="#0E7490" strokeWidth="2" />
          {/* Right Earcup */}
          <rect x="114" y="56" width="12" height="24" rx="6" fill="url(#zara-headphones)" stroke="#0E7490" strokeWidth="2" />
        </g>
      </g>

      {/* Arms & Hands */}
      {isCorrect ? (
        // DJ Celebration Arms
        <g>
          <path d="M46 108 C30 90, 18 64, 28 50 C36 42, 46 64, 54 94 Z" fill="url(#zara-bomber)" />
          <circle cx="26" cy="48" r="8" fill="#EA580C" />
          <path d="M114 108 C130 90, 142 64, 132 50 C124 42, 114 64, 106 94 Z" fill="url(#zara-bomber)" />
          <circle cx="134" cy="48" r="8" fill="#EA580C" />
          {/* Neon music notes / sparkles */}
          <g fill="#06B6D4">
            <path d="M22 28 L28 22 L28 36 L24 38 Z" />
            <circle cx="23" cy="38" r="3" />
            <path d="M136 28 L142 22 L142 36 L138 38 Z" />
            <circle cx="137" cy="38" r="3" />
          </g>
        </g>
      ) : isWrong ? (
        // Hand to cheek in slight bewilderment
        <g>
          <path d="M46 108 C34 122, 38 142, 48 144 C56 146, 58 132, 54 112 Z" fill="url(#zara-bomber)" />
          <circle cx="48" cy="144" r="7" fill="#EA580C" />
          <path d="M114 108 C126 95, 128 72, 118 62 C110 56, 106 68, 104 88 Z" fill="url(#zara-bomber)" />
          <circle cx="118" cy="62" r="7" fill="#EA580C" />
          {/* Sweat Drop */}
          <path d="M126 44 C126 44, 132 52, 132 57 C132 60, 129 63, 126 63 C123 63, 120 60, 120 57 C120 52, 126 44, 126 44 Z" fill="#38BDF8" />
        </g>
      ) : (
        // Chill DJ posture with hands resting
        <g>
          <path d="M45 108 C36 122, 40 142, 54 144 C62 144, 60 130, 52 110 Z" fill="url(#zara-bomber)" />
          <circle cx="54" cy="144" r="7" fill="#EA580C" />
          <path d="M115 108 C124 122, 120 142, 106 144 C98 144, 100 130, 108 110 Z" fill="url(#zara-bomber)" />
          <circle cx="106" cy="144" r="7" fill="#EA580C" />
        </g>
      )}
    </svg>
  );
}
