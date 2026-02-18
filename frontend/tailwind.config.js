/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        mono: ['"Fira Code"', 'monospace'], // Se tiver Fira Code instalada, fica ótimo
        sans: ['Inter', 'sans-serif'],
      },
      colors: {
        // Cores personalizadas se necessário
        'cyber-black': '#050505',
      }
    },
  },
  plugins: [],
}