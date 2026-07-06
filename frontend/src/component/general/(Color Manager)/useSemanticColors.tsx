"use client";

import { useTheme } from "@/component/general/(Color Manager)/ThemeController";
import { DarkSemantic, LightSemantic } from "@/component/general/(Color Manager)/SemanticColors";

export function useSemanticColors() {
  const { theme } = useTheme();

  return theme === "Dark" ? DarkSemantic : LightSemantic;
}
