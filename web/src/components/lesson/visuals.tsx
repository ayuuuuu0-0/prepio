/**
 * Named visuals an intro beat can ask for (`visual: bars` in the lesson YAML). Each is a small
 * looping illustration drawn with SVG and CSS, so it works offline and respects reduced motion
 * (the global rule in globals.css stops every animation).
 *
 * Adding a visual means adding it here AND to the allowed list in the lesson validator
 * (services/question/internal/lesson/validate.go), so unknown names fail CI instead of silently
 * rendering nothing.
 */
import type { ReactNode } from "react";

export const VISUAL_NAMES = ["bars", "flow", "pulse"] as const;
export type VisualName = (typeof VISUAL_NAMES)[number];

export function isVisualName(name: string | undefined): name is VisualName {
  return !!name && (VISUAL_NAMES as readonly string[]).includes(name);
}

/** Two bars racing to fill: a slow one and a fast one. */
function Bars() {
  return (
    <div className="w-full max-w-xs space-y-4" role="img" aria-label="A slow bar and a fast bar filling">
      {[
        { label: "Database", color: "#FF6B35", cls: "visual-bar-slow" },
        { label: "Memory", color: "#34D399", cls: "visual-bar-fast" },
      ].map((bar) => (
        <div key={bar.label}>
          <p className="font-mono mb-1 text-[11px] font-bold uppercase tracking-[0.18em]" style={{ color: "#8B92A8" }}>
            {bar.label}
          </p>
          <div className="h-4 overflow-hidden rounded-full" style={{ background: "#242836" }}>
            <div
              className={`h-full rounded-full ${bar.cls}`}
              style={{ background: bar.color, boxShadow: `0 0 14px ${bar.color}88` }}
            />
          </div>
        </div>
      ))}
    </div>
  );
}

function Node({ x, label, color }: { x: number; label: string; color: string }) {
  return (
    <g>
      <rect x={x - 38} y={34} width={76} height={44} rx={12} fill="#1A1D27" stroke={color} strokeWidth={2} />
      <text x={x} y={61} textAnchor="middle" fill="#E8EAED" fontSize={13} fontWeight={700} className="font-display">
        {label}
      </text>
    </g>
  );
}

/** A request travelling App to Cache and back, with the database path drawn faint behind it. */
function Flow() {
  return (
    <svg viewBox="0 0 320 110" className="w-full max-w-sm" role="img" aria-label="A request goes from the app to the cache and back">
      <line x1={78} y1={56} x2={122} y2={56} stroke="#7C6EF5" strokeWidth={2} />
      <line x1={198} y1={56} x2={242} y2={56} stroke="#2E3347" strokeWidth={2} strokeDasharray="5 5" />
      <Node x={40} label="App" color="#7C6EF5" />
      <Node x={160} label="Cache" color="#34D399" />
      <Node x={280} label="Database" color="#4A5068" />
      <circle cy={56} r={6} fill="#B6ACFF" className="visual-flow-dot" />
      <text x={100} y={26} textAnchor="middle" fill="#34D399" fontSize={11} fontWeight={700} className="font-mono">
        hit
      </text>
    </svg>
  );
}

/** Rings pulsing out of a point: a value that keeps being served after the source moved on. */
function Pulse() {
  return (
    <svg viewBox="0 0 160 160" className="h-40 w-40" role="img" aria-label="Rings pulsing outward">
      {[0, 1, 2].map((i) => (
        <circle
          key={i}
          cx={80}
          cy={80}
          r={22}
          fill="none"
          stroke="#7C6EF5"
          strokeWidth={3}
          className="visual-ring"
          style={{ animationDelay: `${i * 0.8}s` }}
        />
      ))}
      <circle cx={80} cy={80} r={20} fill="#7C6EF5" />
      <text x={80} y={86} textAnchor="middle" fill="#fff" fontSize={18} fontWeight={800} className="font-display">
        ?
      </text>
    </svg>
  );
}

const registry: Record<VisualName, () => ReactNode> = {
  bars: Bars,
  flow: Flow,
  pulse: Pulse,
};

/** IntroVisual renders a named visual, or nothing for an unknown name. */
export function IntroVisual({ name }: { name: string | undefined }) {
  if (!isVisualName(name)) return null;
  const Visual = registry[name];
  return <Visual />;
}
