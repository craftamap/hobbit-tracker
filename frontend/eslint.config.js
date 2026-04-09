import eslint from '@eslint/js';
import eslintPluginVue from 'eslint-plugin-vue';
import typescriptEslint from 'typescript-eslint';
import globals from 'globals';

export default typescriptEslint.config([
  {
    ignores: ['**/dist/**', '**/dist-ssr/**', '**/coverage/**'],
  },
  {
    extends: [
      eslint.configs.recommended,
      ...typescriptEslint.configs.recommended,
      ...eslintPluginVue.configs['flat/recommended'],
    ],
    files: ['**/*.{ts,vue}'],
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      globals: globals.browser,
      parserOptions: {
        parser: typescriptEslint.parser,
      },
    },
    rules: {
      "comma-dangle": ["error", "always-multiline"],
      "space-before-function-paren": ["error", "never"],
      "vue/no-reserved-component-names": ["off"],
    },
  },
]);

