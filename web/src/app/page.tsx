import type { Metadata } from "next";
import { Landing } from "@/components/landing/Landing";
import { SessionRedirect } from "@/components/landing/SessionRedirect";

export const metadata: Metadata = {
  title: "Prepio · Level up as an engineer, five minutes a day",
  description:
    "Short, instantly graded lessons in system design, backend and production, low-level design, and DSA, with per-topic mastery, streaks, and weekly leagues.",
};

export default function Home() {
  return (
    <>
      <SessionRedirect />
      <Landing />
    </>
  );
}
