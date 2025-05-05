import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import { TanStackRouterVite } from '@tanstack/router-vite-plugin';
import { resolve } from 'path';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '');
  const apiBaseUrl = env.VITE_API_URL || 'http://localhost:3000';
  console.log(apiBaseUrl)
  return {
    plugins: [
      react(),
      TanStackRouterVite(),
    ],
    build: {
      outDir: '../static/admin-ui',
      emptyOutDir: true,
    },
    server: {
      proxy: {
        '/v1': apiBaseUrl,
        '/auth': apiBaseUrl
      }
    },
    resolve: {
      alias: {
        '@': resolve(__dirname, './src')
      }
    }
  };
});
