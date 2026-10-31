import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { viteSingleFile } from 'vite-plugin-singlefile'

// Dev:   npm run dev  -> http://localhost:5173/app.html (proxies /api to the Go server on :8080)
// Build: npm run build -> single self-contained frontend/index.html, which main.go already serves at "/".
export default defineConfig({
  plugins: [vue(), viteSingleFile()],
  server: {
    port: 5173,
    proxy: { '/api': { target: 'http://localhost:8080', changeOrigin: true } },
  },
  build: { outDir: 'dist', emptyOutDir: true, rollupOptions: { input: 'app.html' } },
})
