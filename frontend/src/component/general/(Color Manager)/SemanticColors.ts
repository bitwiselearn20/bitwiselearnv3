import { Palette } from "./Palette";

/* ================= DARK ================= */
export const DarkSemantic = {
  background: {
    primary: Palette.custom_black_dark,
    secondary: Palette.custom_black_light,
    accent: Palette.custom_gray,
    special: Palette.custom_blue,
  },

  hero: {
    primary: Palette.custom_green,
    primaryFaded: Palette.custom_green_faded,
    secondary: Palette.custom_blue,
    secondaryFaded: Palette.custom_blue_faded,
  },

  text: {
    primary: Palette.custom_white,
    secondary: Palette.custom_gray,
    special: Palette.custom_blue,
    black: Palette.custom_black,
  },

  border: {
    default: Palette.custom_gray,
    faded: Palette.custom_gray_faded,
    special: Palette.custom_blue,
    green: Palette.custom_green,
  },

  icon: {
    primary: Palette.custom_white,
    secondary: Palette.custom_black,
    special: Palette.custom_blue,
  },

  accent: {
    primary: Palette.custom_black_dark,
    secondary: Palette.custom_green,
    special: Palette.custom_blue,
  },
  hover: {
    special: Palette.custom_blue_faded,
    textSpecial: Palette.custom_blue_faded,
  },
};

/* ================= LIGHT ================= */
export const LightSemantic = {
  background: {
    primary: Palette.custom_light,
    secondary: Palette.custom_surface,
    accent: Palette.custom_gray,
    special: Palette.custom_blue,
  },

  hero: {
    primary: Palette.custom_green,
    primaryFaded: Palette.custom_green_faded,
    secondary: Palette.custom_blue,
    secondaryFaded: Palette.custom_blue_faded,
  },

  text: {
    primary: Palette.custom_black,
    secondary: Palette.custom_dark_gray,
    special: Palette.custom_blue,
    black: Palette.custom_black,
  },

  border: {
    default: Palette.custom_gray,
    faded: Palette.custom_gray_faded,
    special: Palette.custom_blue,
    green: Palette.custom_green,
  },

  icon: {
    primary: Palette.custom_black,
    secondary: Palette.custom_gray,
    special: Palette.custom_blue,
  },

  accent: {
    primary: Palette.custom_surface,
    secondary: Palette.custom_green,
    special: Palette.custom_blue,
  },
  hover: {
    special: Palette.custom_blue_faded,
    textSpecial: Palette.custom_blue_faded,
  },
};
