import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  // Served by the Go server under /app/; the documentation site owns /.
  base: '/app/',
  server: {
    port: 5173,
    // In development the Go server runs alongside; the hash with the launch
    // context stays in the browser and never reaches either server.
    proxy: { '/api': 'http://localhost:8080' },
  },
  build: { outDir: 'dist' },
});
