"use client";

import { Suspense } from "react";
import Attempt from "@/component/attempt/Attempt";
import { AttemptMode } from "@/component/attempt/v1/types";
import { useSearchParams } from "next/navigation";

function AttemptContent() {
  const searchParams = useSearchParams();

  const id = searchParams.get("id");
  const type = searchParams.get("type");

  if (!id || !type) {
    return (
      <div className="h-screen flex items-center justify-center text-white/70">
        Invalid attempt link
      </div>
    );
  }

  const mode: AttemptMode =
    type === "assignment" ? "ASSIGNMENT" : "ASSESSMENT";

  return <Attempt id={id} mode={mode} />;
}

export default function AttemptPage() {
  return (
    <Suspense
      fallback={
        <div className="h-screen flex items-center justify-center text-white/70">
          Loading attempt...
        </div>
      }
    >
      <AttemptContent />
    </Suspense>
  );
}
