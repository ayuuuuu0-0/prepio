"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { api } from "@/lib/api";

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    (async () => {
      const ok = await api.ensureSession();
      if (!ok) {
        router.replace("/login");
        return;
      }
      try {
        const profile = await api.getProfile();
        router.replace(profile.onboarding_completed ? "/dashboard" : "/onboarding");
      } catch {
        router.replace("/login");
      }
    })();
  }, [router]);

  return (
    <main className="flex min-h-screen items-center justify-center">
      <p className="text-slate-500">Loading...</p>
    </main>
  );
}
