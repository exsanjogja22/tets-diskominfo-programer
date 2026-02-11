import { defineConfig } from 'vite'

export default defineConfig({
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
      '/report': {
        target: 'http://localhost:8083',
        changeOrigin: true,
      },
    },
  },
})
