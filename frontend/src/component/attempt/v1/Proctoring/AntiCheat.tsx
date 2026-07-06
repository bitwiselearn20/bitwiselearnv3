import { useEffect, useRef } from "react";
import toast from "react-hot-toast";

export function useAntiCheatControls(started: boolean) {
  const startedRef = useRef(started);

  useEffect(() => {
    startedRef.current = started;
  }, [started]);

  useEffect(() => {
    if (!startedRef.current) return;

    const showToast = (message: string) => {
      toast.error(message, {
        duration: 800,
        position: "top-right",
        style: { background: "#000", color: "#fff" },
      });
    };

    const disableContextMenu = (e: MouseEvent) => {
      e.preventDefault();
      showToast("Right click is disabled during the test.");
    };

    const disableCopyPaste = (e: ClipboardEvent) => {
      e.preventDefault();
      showToast("Copy / Paste is disabled during the test.");
    };

    const disableShortcuts = (e: KeyboardEvent) => {
      if (!startedRef.current) return;

      const key = e.key.toLowerCase();
      const ctrlOrCmd = e.ctrlKey || e.metaKey;

      /* ---------------- COPY / PASTE ---------------- */
      if (ctrlOrCmd && ["c", "v", "x", "a"].includes(key)) {
        e.preventDefault();
        showToast("Copy / Paste is disabled during the test.");
        return;
      }

      /* ---------------- TAB NAV ---------------- */
      if (key === "tab") {
        e.preventDefault();
        showToast("Tab navigation is disabled during the test.");
        return;
      }

      /* ---------------- DEVTOOLS BLOCK ---------------- */
      const isDevToolsShortcut =
        key === "f12" ||
        (ctrlOrCmd && e.shiftKey && ["i", "j", "c"].includes(key));

      if (isDevToolsShortcut) {
        e.preventDefault();
        e.stopPropagation();
        showToast("Inspect tools are disabled during the test.");
        return;
      }
    };

    document.addEventListener("contextmenu", disableContextMenu);
    document.addEventListener("copy", disableCopyPaste);
    document.addEventListener("paste", disableCopyPaste);
    document.addEventListener("cut", disableCopyPaste);
    document.addEventListener("keydown", disableShortcuts, true);

    return () => {
      document.removeEventListener("contextmenu", disableContextMenu);
      document.removeEventListener("copy", disableCopyPaste);
      document.removeEventListener("paste", disableCopyPaste);
      document.removeEventListener("cut", disableCopyPaste);
      document.removeEventListener("keydown", disableShortcuts, true);
    };
  }, [started]);
}
