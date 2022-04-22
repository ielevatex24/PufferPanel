const path = require('path')
const { defineConfig } = require('vite')

module.exports = defineConfig({
  build: {
    lib: {
      entry: path.resolve(__dirname, 'src/index.js'),
      name: 'PufferPanel',
      fileName: (format) => `pufferpanel.${format}.js`
    },
    rollupOptions: {
      external: ['axios', 'nanoevents'],
      output: {
        globals: {
          axios: 'axios'
        }
      }
    }
  }
})
