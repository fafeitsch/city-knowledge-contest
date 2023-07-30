import type { Config } from "tailwindcss";

export default {
  content: ["./app/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      ringColor: {
        DEFAULT: "#e5e7eb"
      }
    },
  },
  plugins: [],
} satisfies Config;
