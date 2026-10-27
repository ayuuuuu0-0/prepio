"use client";

type Reaction = "idle" | "correct" | "wrong";

export function ByteCompanion({
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
      aria-label="Byte the Capybara in hoodie"
    >
      <defs>
        <radialGradient id="byte-glow" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#7C6EF5" stopOpacity="0.45" />
          <stop offset="100%" stopColor="#7C6EF5" stopOpacity="0" />
        </radialGradient>
        <linearGradient id="byte-hoodie" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#8B7BFF" />
          <stop offset="60%" stopColor="#6C5CE7" />
          <stop offset="100%" stopColor="#5544D4" />
        </linearGradient>
        <linearGradient id="byte-fur" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stopColor="#B37C4E" />
          <stop offset="100%" stopColor="#8C5832" />
        </linearGradient>
        <linearGradient id="byte-muzzle" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stopColor="#7A4B29" />
          <stop offset="100%" stopColor="#5C351B" />
        </linearGradient>
        <linearGradient id="byte-shoe" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#2D3250" />
          <stop offset="100%" stopColor="#1B1E32" />
        </linearGradient>
      </defs>

      {/* Ambient background glow */}
      <circle cx="80" cy="110" r="75" fill="url(#byte-glow)" />

      {/* Shadow */}
      <ellipse cx="80" cy="186" rx="46" ry="9" fill="#0D0F17" fillOpacity="0.6" />

      {/* Feet & Sneakers */}
      <g className={isCorrect ? "animate-bounce" : ""}>
        {/* Left Foot / Sneaker */}
        <g transform="translate(48, 168)">
          <path d="M4 14 C4 8, 12 6, 22 6 C28 6, 32 10, 32 16 L2 16 Z" fill="url(#byte-shoe)" />
          {/* Sole */}
          <rect x="0" y="14" width="34" height="5" rx="2.5" fill="#FFFFFF" />
          {/* Neon Stripe */}
          <path d="M8 12 Q18 8 26 12" stroke="#A78BFA" strokeWidth="2" strokeLinecap="round" />
        </g>

        {/* Right Foot / Sneaker */}
        <g transform="translate(80, 168)">
          <path d="M4 14 C4 8, 12 6, 22 6 C28 6, 32 10, 32 16 L2 16 Z" fill="url(#byte-shoe)" />
          <rect x="0" y="14" width="34" height="5" rx="2.5" fill="#FFFFFF" />
          <path d="M8 12 Q18 8 26 12" stroke="#A78BFA" strokeWidth="2" strokeLinecap="round" />
        </g>
      </g>

      {/* Body & Hoodie */}
      <g>
        {/* Torso Hoodie Base */}
        <path
          d="M44 115 C42 155, 46 170, 80 170 C114 170, 118 155, 116 115 C116 95, 44 95, 44 115 Z"
          fill="url(#byte-hoodie)"
        />

        {/* Hoodie Kangaroo Pocket */}
        <path
          d="M56 138 C56 138, 64 134, 80 134 C96 134, 104 138, 104 138 C104 154, 98 162, 80 162 C62 162, 56 154, 56 138 Z"
          fill="#5946D6"
          stroke="#4D3BBA"
          strokeWidth="1.5"
        />

        {/* Hoodie Ribbed Bottom Hem */}
        <rect x="47" y="166" width="66" height="6" rx="3" fill="#4B39B7" />

        {/* Hoodie Drawstrings */}
        <path d="M72 106 L70 126" stroke="#EDE9FE" strokeWidth="2.5" strokeLinecap="round" />
        <circle cx="70" cy="127" r="2" fill="#C4B5FD" />
        <path d="M88 106 L90 126" stroke="#EDE9FE" strokeWidth="2.5" strokeLinecap="round" />
        <circle cx="90" cy="127" r="2" fill="#C4B5FD" />
      </g>

      {/* Head Group (Ears, Face, Snout, Glasses) */}
      <g className={isWrong ? "animate-pulse" : ""}>
        {/* Hoodie Hood Rim behind head */}
        <path
          d="M36 68 C34 26, 126 26, 124 68 C124 96, 36 96, 36 68 Z"
          fill="#5544D4"
          stroke="#4434B8"
          strokeWidth="3"
        />

        {/* Capybara Left Ear */}
        <g transform="translate(36, 42)">
          <ellipse cx="6" cy="8" rx="7" ry="9" fill="#8C5832" />
          <ellipse cx="6" cy="8" rx="4" ry="6" fill="#D97706" fillOpacity="0.4" />
        </g>

        {/* Capybara Right Ear */}
        <g transform="translate(112, 42)">
          <ellipse cx="6" cy="8" rx="7" ry="9" fill="#8C5832" />
          <ellipse cx="6" cy="8" rx="4" ry="6" fill="#D97706" fillOpacity="0.4" />
        </g>

        {/* Head Block / Capybara classic rectangular rounded head */}
        <rect x="42" y="36" width="76" height="66" rx="28" fill="url(#byte-fur)" />

        {/* Snout / Muzzle */}
        <path
          d="M52 64 C52 56, 108 56, 108 64 C110 88, 106 98, 80 98 C54 98, 50 88, 52 64 Z"
          fill="url(#byte-muzzle)"
        />

        {/* Nostrils */}
        <ellipse cx="72" cy="74" rx="3.5" ry="2" fill="#2E1B10" />
        <ellipse cx="88" cy="74" rx="3.5" ry="2" fill="#2E1B10" />

        {/* Mouth */}
        {isCorrect ? (
          <path d="M72 82 Q80 90 88 82" stroke="#2E1B10" strokeWidth="2.5" strokeLinecap="round" fill="none" />
        ) : isWrong ? (
          <path d="M74 86 Q80 82 86 86" stroke="#2E1B10" strokeWidth="2.5" strokeLinecap="round" fill="none" />
        ) : (
          <path d="M74 84 Q80 86 86 84" stroke="#2E1B10" strokeWidth="2" strokeLinecap="round" fill="none" />
        )}

        {/* Eyes */}
        {isCorrect ? (
          // Happy curved eyes ^^
          <g stroke="#1F150E" strokeWidth="3" strokeLinecap="round" fill="none">
            <path d="M57 56 Q63 50 69 56" />
            <path d="M91 56 Q97 50 103 56" />
          </g>
        ) : (
          // Calm / focused eyes
          <g>
            <ellipse cx="63" cy="55" rx="5" ry="3.5" fill="#1F150E" />
            <circle cx="65" cy="54" r="1.5" fill="#FFFFFF" />
            <ellipse cx="97" cy="55" rx="5" ry="3.5" fill="#1F150E" />
            <circle cx="99" cy="54" r="1.5" fill="#FFFFFF" />
          </g>
        )}

        {/* Developer Glasses */}
        <g stroke="#FBBF24" strokeWidth="2.5" fill="none">
          {/* Left Frame */}
          <rect x="52" y="46" width="22" height="18" rx="6" stroke="#FBBF24" fill="#3B82F6" fillOpacity="0.15" />
          {/* Right Frame */}
          <rect x="86" y="46" width="22" height="18" rx="6" stroke="#FBBF24" fill="#3B82F6" fillOpacity="0.15" />
          {/* Bridge */}
          <path d="M74 54 L86 54" stroke="#FBBF24" />
          {/* Temples */}
          <path d="M52 52 L42 50" stroke="#FBBF24" />
          <path d="M108 52 L118 50" stroke="#FBBF24" />
        </g>
      </g>

      {/* Arms & Hands */}
      {isCorrect ? (
        // Arms UP in celebratory cheer! \o/
        <g>
          {/* Left Arm Raised */}
          <path
            d="M48 108 C32 90, 20 62, 30 50 C38 42, 48 64, 56 94 Z"
            fill="url(#byte-hoodie)"
          />
          {/* Left Paw Fist */}
          <circle cx="28" cy="48" r="9" fill="#8C5832" />

          {/* Right Arm Raised */}
          <path
            d="M112 108 C128 90, 140 62, 130 50 C122 42, 112 64, 104 94 Z"
            fill="url(#byte-hoodie)"
          />
          {/* Right Paw Fist */}
          <circle cx="132" cy="48" r="9" fill="#8C5832" />

          {/* Victory Sparkles */}
          <g fill="#FBBF24">
            <path d="M18 36 L21 28 L29 25 L21 22 L18 14 L15 22 L7 25 L15 28 Z" />
            <path d="M142 36 L145 28 L153 25 L145 22 L142 14 L139 22 L131 25 L139 28 Z" />
          </g>
        </g>
      ) : isWrong ? (
        // Arm touching head in "oops / thinking" facepalm
        <g>
          {/* Left Arm Resting */}
          <path
            d="M45 106 C36 122, 38 144, 48 148 C56 150, 56 134, 52 110 Z"
            fill="url(#byte-hoodie)"
          />
          <circle cx="48" cy="150" r="7" fill="#8C5832" />

          {/* Right Arm Bent scratching head */}
          <path
            d="M115 106 C128 95, 132 70, 120 60 C110 52, 105 68, 105 88 Z"
            fill="url(#byte-hoodie)"
          />
          <circle cx="118" cy="58" r="8" fill="#8C5832" />

          {/* Blue Sweat Drop */}
          <path
            d="M128 42 C128 42, 134 50, 134 55 C134 58, 131 61, 128 61 C125 61, 122 58, 122 55 C122 50, 128 42, 128 42 Z"
            fill="#38BDF8"
          />
        </g>
      ) : (
        // Idle: hands resting comfortably in hoodie pocket / sides
        <g>
          {/* Left Arm in pocket */}
          <path
            d="M45 106 C36 122, 40 144, 58 144 C66 144, 62 130, 52 108 Z"
            fill="url(#byte-hoodie)"
          />
          {/* Right Arm in pocket */}
          <path
            d="M115 106 C124 122, 120 144, 102 144 C94 144, 98 130, 108 108 Z"
            fill="url(#byte-hoodie)"
          />
        </g>
      )}
    </svg>
  );
}
