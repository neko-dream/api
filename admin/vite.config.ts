import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

// 管理画面。dev では API(localhost:3000) を /v1・/auth でプロキシし、Cookie を同一オリジン扱いにする。
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3300,
    proxy: {
      "/v1": "http://localhost:3000",
      "/auth": "http://localhost:3000",
    },
  },
});
