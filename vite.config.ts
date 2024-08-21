import vue from "@vitejs/plugin-vue";
import path from "path";
import { defineConfig } from "vite";

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    port: 3001,
    cors: false,
  },
  plugins: [vue()],
  resolve: {
    alias: {
      "@components/*": path.resolve(__dirname, "./frontend/components/*"),
    },
  },
  assetsInclude: ["./frontend/assets/*"],
});
