import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import Components from 'unplugin-vue-components/vite'
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers'
import viteCompression from 'vite-plugin-compression';
import removeConsole from 'vite-plugin-remove-console';
import { resolve } from 'path'
import { nodePolyfills } from 'vite-plugin-node-polyfills'

// https://vitejs.dev/config/
export default defineConfig({
  base: './',
  plugins: [
    vue(),
    vueJsx(),
    Components({
      resolvers: [AntDesignVueResolver()]
    }),
    nodePolyfills(),
    // viteCompression({
    //   deleteOriginFile: false, //删除源文件
    //   threshold: 10240, //压缩前最小文件大小
    //   algorithm: 'gzip', //压缩算法
    //   ext: '.gz', //文件类型
    // }),
    // removeConsole()
  ],
  resolve: {
    alias: [
      { find: '@', replacement: resolve(__dirname, 'src') },
      { find: /^~/, replacement: '' },
    ]
  },
  // 开启less支持
  css: {
    preprocessorOptions: {
      less: {
        javascriptEnabled: true
      }
    }
  },
  server: {
    host: true,
    proxy: {      
      '/api/schema/': {
        target: `http://127.0.0.1:8000/schema`, // 后端接口地址
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/schema/, '')
      },
      '/api/api/': {
        target: `http://127.0.0.1:8000/api`, // 后端接口地址
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/api/, '')
      },
      '/api': {
        target: `http://127.0.0.1:8000/api`, // 后端接口地址
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
      
      
    }
  },
  build: {
    rollupOptions: {
      input: 'doc.html',
      output: {
        chunkFileNames: 'webjars/js/[name]-[hash].js',
        entryFileNames: 'webjars/js/[name]-[hash].js',
        assetFileNames: 'webjars/[ext]/[name]-[hash].[ext]',
        manualChunks: () => 'everything'
      }
    }
  }
})
