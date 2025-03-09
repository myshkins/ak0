import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react'
import * as path from 'path';

export default defineConfig({
  plugins: [react()],
  build: {
    lib: {
      name: 'prerender_components',
      entry: './prerender_components.jsx',
      formats: ['es'],
    },
    rollupOptions: {
      external: ['react', 'react-dom'] // if your script uses React
    },
    sourcemap: true, // Add this for better debugging
    logLevel: 'debug' // Add this for more detailed build output
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
      '../': path.resolve(__dirname, '../')
    }
  }
});

