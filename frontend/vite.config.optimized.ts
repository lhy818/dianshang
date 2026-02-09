import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'
import { visualizer } from 'rollup-plugin-visualizer'
import viteCompression from 'vite-plugin-compression'
import { ViteImageOptimizer } from 'vite-plugin-image-optimizer'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => ({
  plugins: [
    react({
      // 启用Fast Refresh
      fastRefresh: true,
      // 启用JSX运行时自动导入（减少包大小）
      jsxRuntime: 'automatic',
    }),
    // 包分析器 - 只在分析模式下启用
    mode === 'analyze' && visualizer({
      open: true,
      filename: 'dist/stats.html',
      gzipSize: true,
      brotliSize: true,
    }),
    // Gzip和Brotli压缩
    viteCompression({
      verbose: true,
      disable: false,
      threshold: 10240, // 10KB以上才压缩
      algorithm: 'gzip',
      ext: '.gz',
    }),
    viteCompression({
      verbose: true,
      disable: false,
      threshold: 10240,
      algorithm: 'brotliCompress',
      ext: '.br',
    }),
    // 图片优化
    ViteImageOptimizer({
      // 图片优化配置
      png: {
        quality: 80,
      },
      jpeg: {
        quality: 80,
      },
      jpg: {
        quality: 80,
      },
      webp: {
        quality: 80,
      },
      avif: {
        quality: 70,
      },
    }),
  ].filter(Boolean),
  
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@components': path.resolve(__dirname, './src/components'),
      '@features': path.resolve(__dirname, './src/features'),
      '@hooks': path.resolve(__dirname, './src/hooks'),
      '@utils': path.resolve(__dirname, './src/utils'),
      '@types': path.resolve(__dirname, './src/types'),
    },
  },
  
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
    // 启用HMR热更新
    hmr: {
      overlay: true,
    },
  },
  
  build: {
    // 生产环境构建配置
    target: 'es2020',
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: mode === 'production', // 生产环境移除console
        drop_debugger: true,
      },
    },
    // 代码分割配置
    rollupOptions: {
      output: {
        manualChunks: {
          // 将第三方库拆分为单独的chunk
          'react-vendor': ['react', 'react-dom', 'react-router-dom'],
          'ui-vendor': ['antd', '@ant-design/icons'],
          'utils-vendor': ['axios', 'zustand', '@tanstack/react-query'],
          'charts-vendor': ['recharts'],
          'animation-vendor': ['framer-motion'],
        },
        // 文件命名策略
        chunkFileNames: 'assets/js/[name]-[hash].js',
        entryFileNames: 'assets/js/[name]-[hash].js',
        assetFileNames: 'assets/[ext]/[name]-[hash].[ext]',
      },
    },
    // 资源文件大小限制警告
    chunkSizeWarningLimit: 1000, // 1MB
    // 启用CSS代码分割
    cssCodeSplit: true,
    // 生成sourcemap（生产环境建议false）
    sourcemap: mode !== 'production',
  },
  
  // 预加载配置
  preview: {
    port: 4173,
  },
  
  // 环境变量配置
  define: {
    'process.env.NODE_ENV': JSON.stringify(mode),
  },
  
  // 优化依赖预构建
  optimizeDeps: {
    include: [
      'react',
      'react-dom',
      'react-router-dom',
      'antd',
      '@ant-design/icons',
      'axios',
    ],
    exclude: ['@babel/runtime'],
  },
}))