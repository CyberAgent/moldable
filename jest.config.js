/** @type {import("jest").Config} */
module.exports = {
  roots: ["<rootDir>"],
  verbose: true,
  collectCoverage: true,
  collectCoverageFrom: ["src/**/*.{ts,tsx,js,jsx}"],
  moduleFileExtensions: ["ts", "tsx", "js", "jsx", "json", "node"],
  modulePathIgnorePatterns: ["<rootDir>/node_modules", "<rootDir>/dist", "<rootDir>/__mocks__"],
  transform: {
    "^.+\\.(t|j)sx?$": [
      "@swc/jest",
      {
        sourceMaps: true,
        jsc: {
          transform: {
            react: {
              runtime: "automatic",
            },
          },
        },
      },
    ],
  },
  transformIgnorePatterns: ["[/\\\\]node_modules[/\\\\].+\\.(js|jsx)$"],
  globals: {
    "ts-jest": {
      tsconfig: "<rootdir>/tsconfig.json",
    },
  },
  testEnvironment: "node",
};
