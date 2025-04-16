import { defineConfig } from 'vite'

export default defineConfig({
  server: {
    host: '0.0.0.0',
    port: 3000,
    strictPort: true
  },
  build: {
    outDir: 'dist',
    assetsDir: 'css',
    rollupOptions: {
      output: {
        assetFileNames: 'css/main.[ext]' // This ensures CSS gets output to css/main.css
      }
    }
  }
})
