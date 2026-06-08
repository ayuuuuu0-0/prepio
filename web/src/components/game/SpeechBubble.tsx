import { ReactNode } from "react";

/** SpeechBubble displays companion dialogue in a friendly game style. */
export function SpeechBubble({ children, className = "" }: { children: ReactNode; className?: string }) {
  return (
    <div
      className={`speech-bubble relative rounded-3xl bg-white px-5 py-4 text-sm font-semibold leading-relaxed text-[#3C3C3C] shadow-md ${className}`}
    >
      {children}
    </div>
  );
}
