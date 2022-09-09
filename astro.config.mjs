import { defineConfig } from "astro/config";
import Windicss from "vite-plugin-windicss";
import prefetch from "@astrojs/prefetch";

import solidJs from "@astrojs/solid-js";

// https://astro.build/config
export default defineConfig({
  integrations: [prefetch(), solidJs()],
  vite: {
    plugins: [Windicss()],
  },
});
