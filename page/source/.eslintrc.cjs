module.exports = {
    env: {
      browser: true,
      es2021: true,
      node: true
    },
    extends: ["eslint:recommended", "plugin:vue/vue3-essential", "plugin:@typescript-eslint/recommended", "./.eslintrc-auto-import.json"],
    overrides: [],
    parser: "vue-eslint-parser",
    parserOptions: {
      ecmaVersion: 6,
      sourceType: "module",
      ecmaFeatures: {
        modules: true
      },
      requireConfigFile: false,
      parser: "@typescript-eslint/parser"
    },
    plugins: ["vue", "@typescript-eslint"],
    rules: {
      "space-before-function-paren": 0,
      "vue/multi-word-component-names": "off",
      "no-unused-vars": "off", // 未使用变量
      "no-debugger": "off",
      eqeqeq: [2, "allow-null"],
      "spaced-comment": 2, // 注释自动后面空两格
      "no-var": 2,
      "vue/padding-line-between-blocks": "error" // 块之间要隔一行
    }
  }