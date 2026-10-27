"use client";

type Reaction = "idle" | "correct" | "wrong";

export function PipCompanion({
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
      aria-label="Pip the Red Panda in runner jacket"
    >
      <defs>
        <radialGradient id="pip-glow" cx="50%" cy="50%" r="50%">
          <stop offset="0%" stopColor="#FF6B35" stopOpacity="0.45" />
          <stop offset="100%" stopColor="#FF6B35" stopOpacity="0" />
        </radialGradient>
        <linearGradient id="pip-fur" x1="0%" y1="0%" x2="0%" y2="100%">
          <stop offset="0%" stopColor="#F97316" />
          <stop offset="100%" stopColor="#C2410C" />
        </linearGradient>
        <linearGradient id="pip-jacket" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#EF4444" />
          <stop offset="100%" stopColor="#B91C1C" />
        </linearGradient>
        <linearGradient id="pip-tail" x1="0%" y1="0%" x2="100%" y2="100%">
          <stop offset="0%" stopColor="#EA580C" />
          <stop offset="100%" stopColor="#9A3412" />
        </linearGradient>
      </defs>

      {/* Ambient background glow */}
      <circle cx="80" cy="110" r="75" fill="url(#pip-glow)" />

      {/* Shadow */}
      <ellipse cx="80" cy="186" rx="46" ry="9" fill="#0D0F17" fillOpacity="0.6" />

      {/* Fluffy Striped Tail Curling Up */}
      <g className={isCorrect ? "animate-spin-slow" : ""}>
        <path
          d="M106 148 C135 155, 155 125, 145 95 C138 75, 118 78, 122 98 C125 112, 115 128, 98 136 Z"
          fill="url(#pip-tail)"
        />
        {/* Tail Rings / Stripes */}
        <path d="M142 102 C138 98, 128 100, 124 106" stroke="#FEF08A" strokeWidth="6" strokeLinecap="round" />
        <path d="M145 122 C140 116, 128 120, 120 126" stroke="#FEF08A" strokeWidth="6" strokeLinecap="round" />
        <path d="M136 142 C130 138, 118 140, 110 144" stroke="#FEF08A" strokeWidth="6" strokeLinecap="round" />
      </g>

      {/* Backpack behind body */}
      <rect x="52" y="94" width="56" height="52" rx="14" fill="#1E293B" stroke="#0F172A" strokeWidth="2" />
      <path d="M70 94 L70 82 Q80 78 90 82 L90 94" stroke="#334155" strokeWidth="4" fill="none" strokeLinecap="round" />

      {/* Feet & Running Kicks */}
      <g>
        {/* Left Foot */}
        <g transform="translate(48, 168)">
          <path d="M4 14 C4 8, 12 6, 22 6 C28 6, 32 10, 32 16 L2 16 Z" fill="#0F172A" />
          <rect x="0" y="14" width="34" height="5" rx="2.5" fill="#F8FAFC" />
          <path d="M8 12 Q18 8 26 12" stroke="#EF4444" strokeWidth="2" strokeLinecap="round" />
        </g>
        {/* Right Foot */}
        <g transform="translate(80, 168)">
          <path d="M4 14 C4 8, 12 6, 22 6 C28 6, 32 10, 32 16 L2 16 Z" fill="#0F172A" />
          <rect x="0" y="14" width="34" height="5" rx="2.5" fill="#F8FAFC" />
          <path d="M8 12 Q18 8 26 12" stroke="#EF4444" strokeWidth="2" strokeLinecap="round" />
        </g>
      </g>

      {/* Torso & Runner Jacket */}
      <g>
        {/* Jacket Body */}
        <path
          d="M48 110 C46 148, 50 166, 80 166 C110 166, 114 148, 112 110 C112 96, 48 96, 48 110 Z"
          fill="url(#pip-jacket)"
        />
        {/* White Sport Racing Stripe down chest */}
        <path d="M80 100 L80 166" stroke="#FFFFFF" strokeWidth="3.5" strokeLinecap="round" />
        {/* Jacket Zipper Slider */}
        <rect x="78" y="106" width="4" height="8" rx="2" fill="#FCD34D" />
        {/* Jacket Collar */}
        <path d="M60 100 Q80 110 100 100" stroke="#991B1B" strokeWidth="3" fill="none" />
        {/* Backpack Straps across shoulders */}
        <path d="M56 102 L56 156" stroke="#334155" strokeWidth="3" strokeLinecap="round" />
        <path d="M104 102 L104 156" stroke="#334155" strokeWidth="3" strokeLinecap="round" />
      </g>

      {/* Head Group */}
      <g className={isWrong ? "animate-pulse" : ""}>
        {/* Left Ear */}
        <g transform="translate(36, 34)">
          <polygon points="14,0 0,30 28,30" fill="#C2410C" />
          <polygon points="14,6 6,26 22,26" fill="#FFFFFF" />
        </g>
        {/* Right Ear */}
        <g transform="translate(96, 34)">
          <polygon points="14,0 0,30 28,30" fill="#C2410C" />
          <polygon points="14,6 6,26 22,26" fill="#FFFFFF" />
        </g>

        {/* Head Main Circle */}
        <circle cx="80" cy="68" r="36" fill="url(#pip-fur)" />

        {/* White Cheek & Brow Markings */}
        {/* Left Cheek */}
        <ellipse cx="58" cy="76" rx="12" ry="10" fill="#FFFFFF" />
        {/* Right Cheek */}
        <ellipse cx="102" cy="76" rx="12" ry="10" fill="#FFFFFF" />
        {/* White Brows */}
        <ellipse cx="64" cy="50" rx="6" ry="4" fill="#FFFFFF" />
        <ellipse cx="96" cy="50" rx="6" ry="4" fill="#FFFFFF" />

        {/* Red Panda Tear Stripes from eyes to cheeks */}
        <path d="M66 68 Q64 80 58 84" stroke="#7C2D12" strokeWidth="3" strokeLinecap="round" fill="none" />
        <path d="M94 68 Q96 80 102 84" stroke="#7C2D12" strokeWidth="3" strokeLinecap="round" fill="none" />

        {/* Snout */}
        <ellipse cx="80" cy="75" rx="12" ry="9" fill="#FFFFFF" />
        <path d="M76 72 Q80 69 84 72 L80 76 Z" fill="#18181B" />

        {/* Mouth */}
        {isCorrect ? (
          <path d="M74 80 Q80 88 86 80" stroke="#18181B" strokeWidth="2" strokeLinecap="round" fill="none" />
        ) : isWrong ? (
          <path d="M75 83 Q80 79 85 83" stroke="#18181B" strokeWidth="2" strokeLinecap="round" fill="none" />
        ) : (
          <path d="M76 80 Q80 83 84 80" stroke="#18181B" strokeWidth="2" strokeLinecap="round" fill="none" />
        )}

        {/* Eyes */}
        {isCorrect ? (
          // Joyful smiling eyes
          <g stroke="#18181B" strokeWidth="3" strokeLinecap="round" fill="none">
            <path d="M60 62 Q66 56 72 62" />
            <path d="M88 62 Q94 56 100 62" />
          </g>
        ) : (
          // Big lively eyes
          <g>
            <circle cx="66" cy="62" r="5" fill="#18181B" />
            <circle cx="68" cy="60" r="1.8" fill="#FFFFFF" />
            <circle cx="94" cy="62" r="5" fill="#18181B" />
            <circle cx="96" cy="60" r="1.8" fill="#FFFFFF" />
          </g>
        )}
      </g>

      {/* Arms & Hands */}
      {isCorrect ? (
        // Celebrating victory arms
        <g>
          {/* Left Arm High */}
          <path d="M50 106 C34 88, 22 62, 32 48 C40 40, 50 62, 58 92 Z" fill="url(#pip-jacket)" />
          <circle cx="30" cy="46" r="8" fill="#18181B" />
          {/* Right Arm High */}
          <path d="M110 106 C126 88, 138 62, 128 48 C120 40, 110 62, 102 92 Z" fill="url(#pip-jacket)" />
          <circle cx="130" cy="46" r="8" fill="#18181B" />
          {/* Golden Stars */}
          <g fill="#F59E0B">
            <polygon points="20,26 23,20 30,18 24,14 25,7 19,11 13,8 15,15 10,19 17,21" />
            <polygon points="140,26 143,20 150,18 144,14 145,7 139,11 133,8 135,15 130,19 137,21" />
          </g>
        </g>
      ) : isWrong ? (
        // Shrugging apologetic posture
        <g>
          {/* Left Hand out */}
          <path d="M50 106 C36 116, 28 130, 36 136 C44 140, 52 128, 56 112 Z" fill="url(#pip-jacket)" />
          <circle cx="34" cy="136" r="7" fill="#18181B" />
          {/* Right Arm scratching ear */}
          <path d="M110 106 C124 95, 128 72, 116 62 C108 56, 104 68, 102 88 Z" fill="url(#pip-jacket)" />
          <circle cx="116" cy="60" r="7" fill="#18181B" />
          {/* Sweat drop */}
          <path d="M126 44 C126 44, 132 52, 132 57 C132 60, 129 63, 126 63 C123 63, 120 60, 120 57 C120 52, 126 44, 126 44 Z" fill="#38BDF8" />
        </g>
      ) : (
        // Ready runner posture
        <g>
          <path d="M48 106 C38 120, 42 140, 52 142 C60 144, 60 130, 54 110 Z" fill="url(#pip-jacket)" />
          <circle cx="52" cy="144" r="7" fill="#18181B" />
          <path d="M112 106 C122 120, 118 140, 108 142 C100 144, 100 130, 106 110 Z" fill="url(#pip-jacket)" />
          <circle cx="108" cy="144" r="7" fill="#18181B" />
        </g>
      )}
    </svg>
  );
}
