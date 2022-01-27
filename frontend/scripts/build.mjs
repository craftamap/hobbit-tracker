import esbuild from 'esbuild'
import vuePlugin from 'esbuild-plugin-vue3'
import { htmlPlugin } from '@craftamap/esbuild-plugin-html'
import fs from 'node:fs/promises'
import workboxBuild from 'workbox-build'
import yargs from 'yargs'
import { hideBin } from 'yargs/helpers'

const options = yargs(hideBin(process.argv))
  .option('mode', {
    default: 'production',
    choices: ['development', 'production'],
  })
  .option('watch', {
    default: false,
    type: 'boolean',
  })
  .option('cleanFirst', {
    default: false,
    type: 'boolean',
  })
  .parse()

if (options?.cleanFirst) {
  console.log('cleanFirst')
  const r = await fs.rm('dist', { recursive: true })
  console.log('cleanFirst', r)
}

esbuild.build({
  entryPoints: ['src/main.ts'],
  bundle: true,
  metafile: true,
  splitting: true,
  minify: options.mode === 'production',
  sourcemap: options.mode === 'development',
  watch: options?.watch,
  format: 'esm',
  logLevel: 'info',
  outdir: 'dist/',
  loader: { '.woff2': 'file', '.woff': 'file' },
  plugins: [vuePlugin(), htmlPlugin({
    files: [
      {
        entryPoints: 'src/main.ts',
        filename: 'index.html',
        htmlTemplate: (await fs.readFile('public/index.html')).toString(),
        scriptLoading: 'module',
      },
    ],
  }), {
    name: 'workbox-build',
    setup(build) {
      build.onEnd(async() => {
        console.log('workbox-build', 'building')
        const result = await workboxBuild.generateSW({
          globDirectory: 'dist',
          globPatterns: ['**/*.{html,js,css,woff,woff2}'],
          swDest: 'dist/service-worker.js',
        })
        console.log('workbox-build', result)
      })
    },
  }],
})
