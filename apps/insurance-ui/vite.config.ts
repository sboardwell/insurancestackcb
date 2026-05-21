import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import path from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  base: '/',
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3000,
    host: true,
    proxy: {
      '/api/policies': {
        target: process.env.VITE_POLICIES_SERVICE_URL || 'http://localhost:8001',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/policies/, ''),
      },
      '/api/claims': {
        target: process.env.VITE_CLAIMS_SERVICE_URL || 'http://localhost:8002',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/claims/, ''),
      },
      '/api/quotes': {
        target: process.env.VITE_PRICING_SERVICE_URL || 'http://localhost:8003',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/quotes/, ''),
      },
      '/api/customers': {
        target: process.env.VITE_CUSTOMERS_SERVICE_URL || 'http://localhost:8004',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/customers/, ''),
      },
      '/api/payments': {
        target: process.env.VITE_PAYMENTS_SERVICE_URL || 'http://localhost:8005',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/payments/, ''),
      },
    },
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts',
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
    },
  },
});

