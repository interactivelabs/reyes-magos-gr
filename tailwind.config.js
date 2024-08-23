/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/**/*.{html,js,templ,go}", "./asseets/css/main.css"],
  theme: {
    extend: {
      backgroundImage: {
        "hero-paper": "url('https://static.dl-toys.com/bg_water02_sm.webp')",
        "footer-texture": "url('https://static.dl-toys.com/footer_bg.webp')",
      },
      fontFamily: {
        display: ['"Indie Flower"', "cursive"],
      },
      colors: {
        "brand-orange": "#F4673D",
        "brand-blue": "#2A3245",
      },
    },
  },
  plugins: [require("@tailwindcss/aspect-ratio")],
};
