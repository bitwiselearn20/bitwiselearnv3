import { useEffect, useState } from "react";

export function useFullscreenEnforcer(onExit: () => void) {
  const [isFullscreen, setIsFullscreen] = useState(false);

  const enterFullscreen = async () => {
    if (!document.fullscreenElement) {
      await document.documentElement.requestFullscreen();
    }
  };

  useEffect(() => {
    const handleChange = () => {
      const fs = !!document.fullscreenElement;
      setIsFullscreen(fs);

      if (!fs) {
        onExit();
      }
    };

    document.addEventListener("fullscreenchange", handleChange);
    return () => document.removeEventListener("fullscreenchange", handleChange);
  }, [onExit]);

  return { enterFullscreen, isFullscreen };
}
