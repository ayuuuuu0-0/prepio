/** ReadinessRing shows company-specific readiness as an animated SVG ring. */
export function ReadinessRing({
  company,
  score,
  color,
  delay = 0,
}: {
  company: string;
  score: number;
  color: string;
  delay?: number;
}) {
  const radius = 44;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (score / 100) * circumference;

  return (
    <div className="flex flex-col items-center gap-2">
      <div className="relative h-24 w-24">
        <div
          className="absolute inset-2 rounded-full blur-md opacity-20"
          style={{ background: color }}
          aria-hidden
        />
        <svg className="relative h-24 w-24 -rotate-90" viewBox="0 0 100 100">
          <circle cx="50" cy="50" r={radius} fill="none" stroke="#2E3347" strokeWidth="7" />
          <circle
            cx="50"
            cy="50"
            r={radius}
            fill="none"
            stroke={color}
            strokeWidth="7"
            strokeLinecap="round"
            strokeDasharray={circumference}
            strokeDashoffset={offset}
            className="animate-ring-fill"
            style={{ animationDelay: `${delay}ms` }}
          />
        </svg>
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <span className="font-mono text-base font-bold leading-none" style={{ color }}>
            {score}%
          </span>
        </div>
      </div>
      <span
        className="font-display text-[11px] font-semibold uppercase tracking-wider capitalize"
        style={{ color: "#8B92A8" }}
      >
        {company}
      </span>
    </div>
  );
}
