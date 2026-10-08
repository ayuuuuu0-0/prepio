import type { Metadata } from "next";
import localFont from "next/font/local";
import "./globals.css";

// Self-hosted (src/fonts, SIL Open Font License) so builds never depend on Google Fonts.
const sora = localFont({
  src: "../fonts/sora-latin-wght-normal.woff2",
  variable: "--font-display",
  weight: "100 800",
  display: "swap",
});

const manrope = localFont({
  src: "../fonts/manrope-latin-wght-normal.woff2",
  variable: "--font-body",
  weight: "200 800",
  display: "swap",
});

const ibmPlexMono = localFont({
  src: [
    { path: "../fonts/ibm-plex-mono-latin-400-normal.woff2", weight: "400", style: "normal" },
    { path: "../fonts/ibm-plex-mono-latin-500-normal.woff2", weight: "500", style: "normal" },
    { path: "../fonts/ibm-plex-mono-latin-600-normal.woff2", weight: "600", style: "normal" },
  ],
  variable: "--font-mono",
  display: "swap",
});

export const metadata: Metadata = {
  title: "Prepio — Level Up Your Career",
  description: "A progression game for working engineers: short lessons, instant feedback, and visible mastery.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${sora.variable} ${manrope.variable} ${ibmPlexMono.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
