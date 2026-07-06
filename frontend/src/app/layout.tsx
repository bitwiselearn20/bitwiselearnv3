import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { Toaster } from "react-hot-toast";
import Script from "next/script";
import { ThemeProvider } from "@/component/general/(Color Manager)/ThemeController";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "BitwiseLearn — Master Coding, Placements & Your Tech Career",
  description:
    "BitwiseLearn bridges academic theory and industry reality with coding practice, certifications, mentorship, and placement training for students.",
  keywords: [
    "tech career accelerator",
    "campus placement training",
    "coding practice platform for students",
    "industry certifications AWS Azure GCP",
    "DSA practice platform",
    "college placement management system",
    "online interview preparation",
    "student mentorship program",
  ],
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" data-color-mode="dark">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <ThemeProvider>
        {children}
        </ThemeProvider>
        <Toaster position="top-right" />
        <Script
          src="https://acrobatservices.adobe.com/view-sdk/viewer.js"
          strategy="beforeInteractive"
        />
      </body>
    </html>
  );
}
