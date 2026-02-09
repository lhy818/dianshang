/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          DEFAULT: '#2C5AA0', // 品牌蓝
          light: '#4A7BC1',
          dark: '#1A3A6F',
        },
        secondary: {
          DEFAULT: '#FF6B35', // 活力橙
          light: '#FF8F5E',
          dark: '#CC552A',
        },
        success: '#4CAF50',
        warning: '#FFC107',
        error: '#F44336',
      },
      fontFamily: {
        sans: ['PingFang SC', 'SF Pro Text', 'system-ui', 'sans-serif'],
        display: ['PingFang SC', 'SF Pro Display', 'system-ui', 'sans-serif'],
        mono: ['Monaco', 'Consolas', 'monospace'],
      },
      borderRadius: {
        '4': '4px',
        '8': '8px',
        '12': '12px',
      },
      boxShadow: {
        'card': '0 2px 8px rgba(0, 0, 0, 0.08)',
        'card-hover': '0 4px 16px rgba(0, 0, 0, 0.12)',
        'modal': '0 8px 32px rgba(0, 0, 0, 0.16)',
      },
    },
  },
  plugins: [],
  corePlugins: {
    preflight: false, // 避免与Ant Design样式冲突
  },
}