import { defineConfig } from "astro/config";
import windicss from "astro-windicss";
import yaml from "@rollup/plugin-yaml";

// https://astro.build/config
export default defineConfig({
  integrations: [windicss()],
  vite: {
    build: {
      rollupOptions: {
        plugins: [
          yaml({
            transform(data, file) {
              return { data: { data, file } };
            },
          }),
        ],
      },
    },
  },
});
