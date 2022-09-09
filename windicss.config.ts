import { defineConfig } from "windicss/helpers";

export default defineConfig({
  extract: {
    include: ["src/**/*.{astro,tsx,ts}"],
    exclude: ["node_modules", ".git"],
  },
  theme: {
    extend: {
      colors: {
        c0: "rgb(255, 245, 228)",
        c1: "rgb(255, 196, 196)",
        c2: "rgb(238, 105, 131)",
        c3: "rgb(133, 14, 53)",
      },
    },
  },
});
