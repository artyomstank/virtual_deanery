import type { Config } from "tailwindcss";

export default {
  darkMode: ["class"],
  content: [
    "./index.html",
    "./src/**/*.{ts,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        border: "#27272a",
        background: "#09090b",
        foreground: "#fafafa",
        muted: "#18181b",
        card: "#111111",
        primary: {
          DEFAULT: "#6366f1",
          foreground: "#ffffff"
        }
      },
      borderRadius: {
        xl: "1rem",
        "2xl": "1.5rem"
      },
      boxShadow: {
        soft: "0 10px 40px rgba(0,0,0,0.25)"
      }
    }
  },
  plugins: []
} satisfies Config;