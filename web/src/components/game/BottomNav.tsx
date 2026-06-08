"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const tabs = [
  { href: "/dashboard", label: "Home", icon: "⊞" },
  { href: "/journey", label: "Journey", icon: "◈" },
  { href: "/challenge", label: "Play", icon: "▶" },
  { href: "/league", label: "League", icon: "⬡", soon: true },
  { href: "/quests", label: "Quests", icon: "◇", soon: true },
  { href: "/profile", label: "Profile", icon: "◉" },
];

/** BottomNav provides primary tab navigation on mobile-first layouts. */
export function BottomNav() {
  const path = usePathname();

  return (
    <nav
      className="fixed bottom-0 left-0 right-0 z-50"
      style={{
        background: "rgba(15,17,23,0.95)",
        borderTop: "1px solid #2E3347",
        backdropFilter: "blur(20px)",
      }}
    >
      <div className="mx-auto flex max-w-lg items-center justify-around px-1 py-2">
        {tabs.map((tab) => {
          const active = path === tab.href;
          return (
            <Link
              key={tab.href}
              href={tab.href}
              className={`relative flex flex-col items-center gap-0.5 rounded-xl px-2 py-2 transition-all ${
                tab.soon ? "opacity-40" : ""
              }`}
              style={{
                color: active ? "#7C6EF5" : "#4A5068",
                background: active ? "rgba(124,110,245,0.12)" : "transparent",
              }}
            >
              <span className="text-base leading-none">{tab.icon}</span>
              <span className="font-display text-[9px] font-semibold tracking-wide">
                {tab.label}
              </span>
              {active && (
                <span
                  className="absolute -top-px left-1/4 right-1/4 h-px rounded-full"
                  style={{ background: "#7C6EF5" }}
                />
              )}
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
