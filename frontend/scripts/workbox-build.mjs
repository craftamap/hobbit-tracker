import workboxBuild from 'workbox-build'
import path from 'path'

export const workboxBuildPlugin = () => ({
  name: 'workbox-build',
  setup(build) {
    build.onEnd(async() => {
      const start = new Date()
      console.log('  workbox-build:', 'building')
      const result = await workboxBuild.generateSW({
        globDirectory: 'dist',
        globPatterns: ['**/*.{html,js,css,woff,woff2}'],
        swDest: 'dist/service-worker.js',
      })
      // console.log('workbox-build', result)
      if (result.warnings.length > 0) {
        console.log('  workbox-build:', result.warnings)
      }

      for (const fp of result.filePaths) {
        console.log(' ', path.relative('.', fp))
      }

      console.log('  workbox-build:', `done in ${new Date().getTime() - start.getTime()}ms\n`)
    })
  },
})
