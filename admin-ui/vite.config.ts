import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import { TanStackRouterVite } from '@tanstack/router-vite-plugin';
import { resolve } from 'path';

export default defineConfig(({ mode }) => {
  let apiBaseURL: string;
  switch (mode) {
    case 'production':
      apiBaseURL = 'https://api.kotohiro.com';
      break;
    case 'development':
      apiBaseURL = 'https://api-dev.kotohiro.com';
      break;
    default:
      apiBaseURL = 'http://localhost:3000';
  }

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
        '/v1': apiBaseURL,
        '/auth': apiBaseURL
      }
    },
    resolve: {
      alias: {
        '@': resolve(__dirname, './src')
      }
    }
  };
});
