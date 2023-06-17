import commonjs from '@rollup/plugin-commonjs'
import json from '@rollup/plugin-json'
import resolve from '@rollup/plugin-node-resolve'
import includePaths from 'rollup-plugin-includepaths'
import dotenv from 'rollup-plugin-dotenv'
import livereload from 'rollup-plugin-livereload'
import postcss from 'rollup-plugin-postcss'
import svelte from 'rollup-plugin-svelte'
import { terser } from 'rollup-plugin-terser'
import routify from '@roxi/routify/plugins/rollup'
import autoPreprocess from 'svelte-preprocess'
import typescript from '@rollup/plugin-typescript'
import outputManifest from 'rollup-plugin-output-manifest'
import nodePolyfills from 'rollup-plugin-polyfill-node'

const production = !process.env.ROLLUP_WATCH

const calculateNameForChunk = (chunk) => chunk.fileName.split('.').shift()

export default [
  {
    input: 'assets/application.ts',
    output: {
      dir: './public/assets',
      entryFileNames: 'application.js',
      format: 'iife',
      name: 'app',
      sourcemap: !production,
    },
    plugins: [
      svelte({
        compilerOptions: {
          // enable run-time checks when not in production
          dev: !production,
        },
        emitCss: true, // give component style to postcss() for processing
        preprocess: autoPreprocess(),
      }),

      typescript({ sourceMap: !production }),

      nodePolyfills(),

      // If you have external dependencies installed from
      // npm, you'll most likely need these plugins. In
      // some cases you'll need additional configuration -
      // consult the documentation for details:
      // https://github.com/rollup/plugins/tree/master/packages/commonjs
      resolve({
        browser: true,
        dedupe: ['svelte'],
      }),
      includePaths({
        include: {},
        paths: [
          'assets/components',
          'assets/data',
          'assets/helpers',
          'assets/pages',
          'assets/external',
        ],
        external: [],
        extensions: ['.js', '.ts'],
      }),
      commonjs(),

      json(), // adds support for importing json files
      postcss({
        extract: true, // create a css file alongside the output.file
        sourceMap: production,
        use: {
          sass: {
            includePaths: ['node_modules'],
          },
        },
      }),
      routify({
        dynamicImports: false,
        pages: 'assets/pages',
      }),
      dotenv(),

      //           minify     auto-refresh browser on changes
      production ? terser() : livereload('public/assets'),

      outputManifest({
        generate: (keyValueDecorator, seed, opt) => (chunks) =>
          chunks.reduce((json, chunk) => {
            const fileName = chunk.fileName

            /* Ensure CSS file has a `name` so it's included in manifest.json so
             * its versioned fileName is properly included in our HTML file. */
            const name = chunk.name || calculateNameForChunk(chunk)

            return {
              ...json,
              ...keyValueDecorator(name, fileName, opt),
            }
          }, seed),
      }),
    ],
    watch: {
      clearScreen: false,
    },
  },
]
