/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.{html,js,templ,go}", "./asseets/css/main.css"],
  theme: {
    extend: {
      backgroundImage: {
        "hero-paper":
          "url('https://static.dl-toys.com/img/bg_water02_sm.webp')",
        "footer-texture":
          "url('https://static.dl-toys.com/img/footer_bg.webp')",
      },
      fontFamily: {
        sans: [
          "Cabin",
          "ui-sans-serif",
          "system-ui",
          "sans-serif",
          "Apple Color Emoji",
          "Segoe UI Emoji",
          "Segoe UI Symbol",
          "Noto Color Emoji",
        ],
        display: ["Raleway", "ui-sans-serif", "system-ui", "sans-serif"],
        decorative: ["Indie Flower", "cursive"],
      },
      colors: {
        "brand-orange": "#F4673D",
        "brand-blue": "#2A3245",
      },
    },
  },
  plugins: [require("@tailwindcss/aspect-ratio")],
};
