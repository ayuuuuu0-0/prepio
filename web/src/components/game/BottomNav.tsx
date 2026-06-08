"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const tabs = [
  { href: "/dashboard", label: "Home", icon: "🏠" },
  { href: "/journey", label: "Journey", icon: "🗺️" },
  { href: "/challenge", label: "Play", icon: "⚡" },
];

/** BottomNav provides game-style tab navigation on mobile-first layouts. */
export function BottomNav() {
  const path = usePathname();

  return (
    <nav className="fixed bottom-0 left-0 right-0 z-50 border-t-4 border-[#E5E5E5] bg-white/95 backdrop-blur">
      <div className="mx-auto flex max-w-lg justify-around py-2">
        {tabs.map((tab) => {
          const active = path === tab.href;
          return (
            <Link
              key={tab.href}
              href={tab.href}
              className={`flex flex-col items-center rounded-2xl px-4 py-2 transition ${
                active ? "bg-lime-100 text-[#58CC02]" : "text-[#777]"
              }`}
            >
              <span className="text-xl">{tab.icon}</span>
              <span className="font-display text-xs font-bold">{tab.label}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
