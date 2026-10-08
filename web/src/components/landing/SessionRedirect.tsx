"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { api } from "@/lib/api";

/**
 * SessionRedirect sends a signed-in visitor from the landing page into the app. Visitors
 * without a session simply stay on the landing page.
 */
export function SessionRedirect() {
  const router = useRouter();

  useEffect(() => {
    let cancelled = false;
    (async () => {
      if (!(await api.ensureSession()) || cancelled) return;
      try {
        const profile = await api.getProfile();
        if (!cancelled) router.replace(profile.onboarding_completed ? "/dashboard" : "/onboarding");
      } catch {
        // A stale session just leaves the visitor on the landing page.
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [router]);

  return null;
}
