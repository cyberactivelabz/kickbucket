import type { Config } from "tailwindcss";

// Palette: deep ink-blue base, with amber for "open" findings and coral for "public".
const config: Config = {
  content: ["./app/**/*.{ts,tsx}", "./components/**/*.{ts,tsx}"],
  theme: {
    extend: {
      colors: {
        ink: "#0E1520",
        panel: "#151F2E",
        line: "#26354A",
        text: "#D5DEE8",
        mute: "#8A9BB0",
        amber: "#FFB020",
        coral: "#FF6B5E",
        mint: "#5FD3A5",
      },
      fontFamily: {
        sans: ["var(--font-sans)", "system-ui", "sans-serif"],
        mono: ["var(--font-mono)", "ui-monospace", "monospace"],
      },
    },
  },
  plugins: [],
};
export default config;
