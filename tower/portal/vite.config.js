import { defineConfig } from 'vite'
import postcss from '@tailwindcss/postcss'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  css: {
    postcss
  },
  server: {
    port: 3001, // portal port
    host: '0.0.0.0', // This allows connections from other hostnames
    allowedHosts: ['dev.local', 'portal.dev.local', 'localhost']
  },
  build: {
    outDir: '../portal/dist',
    emptyOutDir: true,
  },
  publicDir: 'public'
})
