import { defineConfig } from 'vite'
import postcss from '@tailwindcss/postcss'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  css: {
    postcss
  },
  build: {
    outDir: '../web/dist',
    emptyOutDir: true,
  },
  publicDir: 'public'
})
