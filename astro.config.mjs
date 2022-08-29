import { defineConfig } from "astro/config";
import windicss from "astro-windicss";

import prefetch from "@astrojs/prefetch";

// https://astro.build/config
export default defineConfig({
  integrations: [windicss(), prefetch()]
});