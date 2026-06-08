import { ReactNode } from "react";

/** SpeechBubble renders companion dialogue in a dark HUD-style bubble. */
export function SpeechBubble({
  children,
  className = "",
  speakerName,
}: {
  children: ReactNode;
  className?: string;
  speakerName?: string;
}) {
  return (
    <div
      className={`speech-bubble relative rounded-2xl border border-white/8 px-4 py-3 text-sm leading-relaxed ${className}`}
      style={{
        background: "#1A1D27",
        boxShadow: "0 4px 20px rgba(0,0,0,0.4), inset 0 1px 0 rgba(255,255,255,0.05)",
      }}
    >
      {speakerName && (
        <p
          className="font-mono mb-1 text-xs font-semibold uppercase tracking-widest"
          style={{ color: "#7C6EF5" }}
        >
          {speakerName}
        </p>
      )}
      <p className="font-body font-medium" style={{ color: "#C8CCDA" }}>
        {children}
      </p>
    </div>
  );
}
