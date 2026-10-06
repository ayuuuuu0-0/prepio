"use client";

import { useEffect, useState } from "react";
import { ringOffset } from "@/lib/topics";

/**
 * MasteryRing draws a topic's mastery (0-100) as a ring. A null mastery means "not started":
 * an empty ring with a dash, never a zero. The fill animates in once on mount.
 */
export function MasteryRing({ mastery, tint, size = 84 }: { mastery: number | null; tint: string; size?: number }) {
  const stroke = 8;
  const radius = (size - stroke) / 2;
  const circumference = 2 * Math.PI * radius;

  const [drawn, setDrawn] = useState<number | null>(null);
  useEffect(() => {
    const id = requestAnimationFrame(() => setDrawn(mastery));
    return () => cancelAnimationFrame(id);
  }, [mastery]);

  const label = mastery === null ? "Not started" : `Mastery ${mastery} out of 100`;

  return (
    <svg
      width={size}
      height={size}
      viewBox={`0 0 ${size} ${size}`}
      role="img"
      aria-label={label}
      className="shrink-0"
    >
      <circle cx={size / 2} cy={size / 2} r={radius} fill="none" stroke="#2E3347" strokeWidth={stroke} />
      <circle
        cx={size / 2}
        cy={size / 2}
        r={radius}
        fill="none"
        stroke={tint}
        strokeWidth={stroke}
        strokeLinecap="round"
        strokeDasharray={circumference}
        strokeDashoffset={ringOffset(drawn, circumference)}
        transform={`rotate(-90 ${size / 2} ${size / 2})`}
        style={{ transition: "stroke-dashoffset 900ms cubic-bezier(0.22, 1, 0.36, 1)", filter: `drop-shadow(0 0 6px ${tint}66)` }}
      />
      <text
        x="50%"
        y="50%"
        textAnchor="middle"
        dominantBaseline="central"
        className="font-display"
        style={{ fill: mastery === null ? "#4A5068" : "#E8EAED", fontSize: size * 0.3, fontWeight: 800 }}
      >
        {mastery === null ? "–" : mastery}
      </text>
    </svg>
  );
}
