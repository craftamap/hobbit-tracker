import { Plugin, defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import workboxBuild from 'workbox-build'

const images = {
  "regular": [60, 76, 120, 144, 150, 152, 180, 192, 512],
  "maskable": [192, 512],
} as const;

const manifestGeneratorPlugin = (): Plugin => ({
  name: 'manifestGenerator',
  generateBundle() {
    this.emitFile({
      type: "asset",
      fileName: "manifest.json",
      source: JSON.stringify({
        name: "Hobbit Tracker",
        short_name: "Hobbit Tracker",
        start_url: ".",
        display: "standalone",
        background_color: "#000000",
        icons: [
          ...images.regular.map((size) => ({
            src: `./img/icons/icon-${size}x${size}.png`,
            sizes: `${size}x${size}`,
            type: 'image/png'
          })),
          ...images.regular.map((size) => ({
            src: `./img/icons/icon-maskable-${size}x${size}.png`,
            sizes: `${size}x${size}`,
            type: 'image/png',
            purpose: 'maskable',
          }))
        ]
      })
    })
  }
})

const workboxPlugin = (): Plugin => ({
  name: 'workbox',
  writeBundle: async () => {
    const result = await workboxBuild.generateSW({
      globDirectory: 'dist',
      globPatterns: ['**/*.{html,js,css,woff,woff2}'],
      swDest: 'dist/service-worker.js',
    })
    console.log(result)
  },
});

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(), manifestGeneratorPlugin(), workboxPlugin()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        secure: false,
      },
      '/auth': {
        target: 'http://localhost:8080',
        secure: false,
      },
    },
  },
})
