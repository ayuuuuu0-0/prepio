"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { api } from "@/lib/api";

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    const token = api.loadToken();
    if (!token) {
      router.replace("/login");
      return;
    }
    api
      .getProfile()
      .then((profile) => router.replace(profile.onboarding_completed ? "/dashboard" : "/onboarding"))
      .catch(() => router.replace("/login"));
  }, [router]);

  return (
    <main className="flex min-h-screen items-center justify-center">
      <p className="text-slate-500">Loading...</p>
    </main>
  );
}
