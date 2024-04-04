module.exports = {
  "src/**/*.{js,jsx,ts,tsx}": ["eslint"],
  "*.{js,jsx,ts,tsx,md}": ["prettier --write", "eslint --fix"],
};
