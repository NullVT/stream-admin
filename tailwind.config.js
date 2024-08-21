import colors from "tailwindcss/colors";
import defaultTheme from "tailwindcss/defaultTheme";

/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./frontend/**/*.{vue,js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Inter var", ...defaultTheme.fontFamily.sans],
      },
      colors: {
        primary: colors.rose[800],
        background: colors.zinc[800],
        base: colors.zinc[900],
        twitch: "#9146FF",
        youtube: "#FF0000",
      },
    },
  },
  plugins: [],
};
