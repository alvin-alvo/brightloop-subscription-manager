/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        stripe: {
          primary: '#533afd',
          'primary-deep': '#4434d4',
          'primary-press': '#2e2b8c',
          'primary-soft': '#665efd',
          'primary-subdued': '#b9b9f9',
          'brand-dark': '#1c1e54',
          ink: '#0d253d',
          'ink-secondary': '#273951',
          'ink-mute': '#64748d',
          'ink-mute-2': '#61718a',
          canvas: '#ffffff',
          'canvas-soft': '#f6f9fc',
          'canvas-cream': '#f5e9d4',
          hairline: '#e3e8ee',
          'hairline-input': '#a8c3de',
          ruby: '#ea2261',
          magenta: '#f96bee',
          lemon: '#9b6829',
          // Dark palette surfaces
          'dark-bg': '#0a101d',
          'dark-surface': '#111a2e',
          'dark-card': '#16223b',
          'dark-border': '#1e2f4f',
          'dark-text': '#f8fafc',
          'dark-mute': '#94a3b8',
        },
      },
      boxShadow: {
        'stripe-sm': '0 1px 3px rgba(0, 55, 112, 0.08)',
        'stripe-md': '0 8px 24px rgba(0, 55, 112, 0.08), 0 2px 6px rgba(0, 55, 112, 0.04)',
        'stripe-lg': '0 16px 36px rgba(0, 55, 112, 0.12), 0 4px 10px rgba(0, 55, 112, 0.06)',
        'dark-sm': '0 1px 3px rgba(0, 0, 0, 0.4)',
        'dark-md': '0 8px 24px rgba(0, 0, 0, 0.5), 0 2px 6px rgba(0, 0, 0, 0.3)',
      },
      fontFamily: {
        sans: ['Inter', 'system-ui', '-apple-system', 'sans-serif'],
      },
    },
  },
  plugins: [],
}
