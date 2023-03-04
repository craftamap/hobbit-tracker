import esbuild from 'esbuild'
import vuePlugin from 'esbuild-plugin-vue3'
import { htmlPlugin } from '@craftamap/esbuild-plugin-html'
import fs from 'node:fs/promises'
import yargs from 'yargs'
import { hideBin } from 'yargs/helpers'
import { workboxBuildPlugin } from './workbox-build.mjs'
import { manifestGeneratorPlugin } from './manifest-generator.mjs'

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
  await fs.rm('dist', { recursive: true })
  console.log('cleanFirst', 'done')
}

const noopPlugin = {
  name: 'noop',
  setup(build) {
    build.onEnd(() => { console.log('  \u001b[37;1m! skipping workboxBuildPlugin\u001b[0m') })
  },
}

const ctx = await esbuild.context({
  entryPoints: ['src/main.ts'],
  bundle: true,
  metafile: true,
  splitting: true,
  minify: options.mode === 'production',
  sourcemap: options.mode === 'development',
  format: 'esm',
  logLevel: 'info',
  outdir: 'dist/',
  loader: { '.woff2': 'file', '.woff': 'file' },
  plugins: [
    vuePlugin(),
    htmlPlugin({
      files: [
        {
          entryPoints: ['src/main.ts'],
          filename: 'index.html',
          htmlTemplate: (await fs.readFile('public/index.html')).toString(),
          scriptLoading: 'module',
        },
      ],
    }),
    {
      name: 'copy',
      setup(build) {
        build.onEnd(async() => {
          await fs.copyFile('public/favicon.png', 'dist/favicon.png')
          console.log('  copy:', 'public/favicon.png -> dist/favicon.png')

          console.log()
        })
      },
    },
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    options.mode === 'production' ? workboxBuildPlugin() : noopPlugin,
    manifestGeneratorPlugin({
      name: 'Hobbit Tracker',
      // eslint-disable-next-line @typescript-eslint/camelcase
      background_color: '#111d1f',
    })],
})

if (options.watch) {
  await ctx.watch();
} else {
  await ctx.rebuild();
  await ctx.dispose();
}
