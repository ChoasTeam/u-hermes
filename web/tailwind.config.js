/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        surface: { DEFAULT: '#0a0a0b', card: '#111113', hover: '#1a1a1e' },
        border: { DEFAULT: '#1e1e22', subtle: '#30363d' },
        primary: { DEFAULT: '#6366f1', hover: '#7c3aed' },
        success: { DEFAULT: '#22c55e', dim: '#238636' },
      },
    },
  },
  plugins: [],
};
