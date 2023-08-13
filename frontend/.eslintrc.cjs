module.exports = {
  root: true,
  parser: "@typescript-eslint/parser",
  extends: [
    "eslint:recommended",
    "plugin:@typescript-eslint/recommended",
    "prettier",
    "plugin:import/warnings",
    "plugin:import/typescript",
  ],
  plugins: ["svelte3", "@typescript-eslint"],
  ignorePatterns: ["*.cjs"],
  overrides: [
    {
      files: ["*.svelte"],
      processor: "svelte3/svelte3",
    },
  ],
  settings: {
    "svelte3/typescript": () => require("typescript"),
  },
  parserOptions: {
    sourceType: "module",
    ecmaVersion: 2020,
  },
  rules: {
    "@typescript-eslint/ban-ts-comment": "off",
    "@typescript-eslint/no-unused-vars": [
      "error",
      {
        argsIgnorePattern: "^_",
        ignoreRestSiblings: true,
        varsIgnorePattern: "^_",
      },
    ],
    "import/order": [
      "warn",
      {
        // https://github.com/import-js/eslint-plugin-import/blob/v2.26.0/docs/rules/order.md
        groups: ["builtin", "external", "internal", "parent", "sibling", "index", "type"],
        pathGroups: [
          {
            pattern: "$app/**",
            group: "external",
          },
          {
            pattern: "$lib/**",
            group: "internal",
          },
        ],
        alphabetize: {
          order: "asc",
          caseInsensitive: true,
        },
        "newlines-between": "always",
      },
    ],
  },
  env: {
    browser: true,
    es2017: true,
    node: true,
  },
}
