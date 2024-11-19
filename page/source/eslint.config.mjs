import eslintJs from '@eslint/js'
import eslintPluginVue  from 'eslint-plugin-vue'

export default [
    eslintJs.configs.recommended,
    eslintPluginVue.configs.recommended,
]