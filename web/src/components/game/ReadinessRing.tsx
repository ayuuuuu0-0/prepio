/** ReadinessRing shows company readiness as a colorful circular progress ring. */
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
        <svg className="h-24 w-24 -rotate-90" viewBox="0 0 100 100">
          <circle cx="50" cy="50" r={radius} fill="none" stroke="#E8E8E8" strokeWidth="8" />
          <circle
            cx="50"
            cy="50"
            r={radius}
            fill="none"
            stroke={color}
            strokeWidth="8"
            strokeLinecap="round"
            strokeDasharray={circumference}
            strokeDashoffset={offset}
            className="animate-ring-fill transition-all duration-1000"
            style={{ animationDelay: `${delay}ms` }}
          />
        </svg>
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <span className="font-display text-lg font-bold" style={{ color }}>
            {score}%
          </span>
        </div>
      </div>
      <span className="font-display text-xs font-bold capitalize text-[#3C3C3C]">{company}</span>
    </div>
  );
}
