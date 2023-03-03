import fs from 'fs/promises'
import path from 'path'

const defaultManifest = {
  icons: [
    {
      src: './img/icons/icon-60x60.png',
      sizes: '60x60',
      type: 'image/png',
    },
    {
      src: './img/icons/icon-76x76.png',
      sizes: '76x76',
      type: 'image/png',
    },
    {
      src: './img/icons/icon-120x120.png',
      sizes: '120x120',
      type: 'image/png',
    },
    {
      src: './img/icons/icon-144x144.png',
      sizes: '144x144',
      type: 'image/png',
    },
    {
      src: './img/icons/icon-150x150.png',
      sizes: '150x150',
      type: 'image/png',
    },
    {
      src: './img/icons/icon-152x152.png',
      sizes: '152x152',
      type: 'image/png',
    },
    {
      src: './img/icons/icon-180x180.png',
      sizes: '180x180',
      type: 'image/png',
    },
    {
      src: './img/icons/icon-192x192.png',
      sizes: '192x192',
      type: 'image/png',
    },
    {
      src: './img/icons/icon-512x512.png',
      sizes: '512x512',
      type: 'image/png',
    },
    {
      src: './img/icons/icon-maskable-192x192.png',
      sizes: '192x192',
      type: 'image/png',
      purpose: 'maskable',
    },
    {
      src: './img/icons/icon-maskable-512x512.png',
      sizes: '512x512',
      type: 'image/png',
      purpose: 'maskable',
    },
  ],
  start_url: '.',
  id: '/',
  display: 'standalone',
  background_color: '#000000',
}

export const manifestGeneratorPlugin = (options) => ({
  name: 'manifest-generator-plugin',
  setup(build) {
    build.onEnd(async() => {
      console.log('  manifest-generator-plugin:', 'generating manifest')
      const start = new Date()

      const publicOptions = {
        name: options?.name,
        short_name: options?.name,
        theme_color: options?.themeColor,
      }

      const outputManifest = Object.assign(
        {}, publicOptions, defaultManifest, options?.manifestOptions,
      )

      await fs.writeFile(path.join(build.initialOptions.outdir, 'manifest.json'), JSON.stringify(outputManifest))
      console.log(' ', path.join(build.initialOptions.outdir, 'manifest.json'))

      for (const icons of outputManifest.icons) {
        await fs.mkdir(path.dirname(path.join(build.initialOptions.outdir, icons.src)), { recursive: true })
        await fs.copyFile(path.join('./public', icons.src), path.join(build.initialOptions.outdir, icons.src))
        console.log(' ', path.join('./public', icons.src), ' -> ', path.join(build.initialOptions.outdir, icons.src))
      }
      console.log('  manifest-generator-plugin:', `done in ${new Date().getTime() - start.getTime()}ms\n`)
    })
  },
})
