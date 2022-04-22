import path from 'path'
import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import eslint from "vite-plugin-eslint"
import legacy from "@vitejs/plugin-legacy"
import vueI18n from '@intlify/vite-plugin-vue-i18n'
import fs from 'fs'

export default defineConfig({
  resolve: {
    alias: [
      {find: "@", replacement: path.resolve(__dirname, 'src')}
    ]
  },
  define: {
    localeList: fs.readdirSync('src/lang')
  },
  plugins: [
    vue(),
    vueI18n({
      runtimeOnly: false,
      include: path.resolve(__dirname, '@/lang/**')
    }),
    eslint(),
    legacy({
      // https://browserslist.dev/?q=bGFzdCAyIHZlcnNpb25zLCBmaXJlZm94IGVzciwgbm90IGRlYWQsIG5vdCBpZSAxMQ%3D%3D
      targets: ["last 2 versions", "firefox esr", "not dead", "not IE 11"]
    })
  ],
  server: {
    proxy: {
      '/auth': {
        target: 'https://nitori.griem.xyz',
        changeOrigin: true
      },
      '/api': {
        target: 'https://nitori.griem.xyz',
        changeOrigin: true,
        ws: true
      },
      '/proxy': {
        target: 'https://nitori.griem.xyz',
        changeOrigin: true,
        ws: true
      }
    }
  }
})
