import react from "@vitejs/plugin-react";
import { defineConfig, loadEnv } from "vite";

// Local API is separate from the public Web dev server on port 3000.
export default defineConfig(({ mode }) => {
  const config = loadEnv(mode, process.cwd(), "");
  const apiTarget = config.LOCAL_API_URL || "http://127.0.0.1:8081";
  return {
    base: "/admin/",
    plugins: [react()],
    build: {
      outDir: "../server/static/admin-ui",
      emptyOutDir: true,
    },
    server: {
      host: "127.0.0.1",
      port: 3300,
      strictPort: true,
      proxy: {
        "/v1": apiTarget,
        "/auth": apiTarget,
      },
    },
  };
});
