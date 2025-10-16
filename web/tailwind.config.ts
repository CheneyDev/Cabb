import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './app/**/*.{ts,tsx}',
    './components/**/*.{ts,tsx}',
  ],
  darkMode: ['class'],
  theme: {
    extend: {
      colors: {
        // COSS-inspired neutrals and accents
        background: '#0A0A0A',
        foreground: '#171717',
        border: '#262626',
        primary: {
          DEFAULT: '#4F46E5',
          foreground: '#ffffff',
        },
        accent: {
          DEFAULT: '#06B6D4',
          foreground: '#0A0A0A',
        },
        muted: '#525252',
      },
      borderRadius: {
        xl: '1rem',
      },
    },
  },
  plugins: [],
}

export default config

