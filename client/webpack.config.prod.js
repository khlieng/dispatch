var path = require('path');
var MiniCssExtractPlugin = require('mini-css-extract-plugin');
var postcssPresetEnv = require('postcss-preset-env');
var cssnano = require('cssnano');
var TerserPlugin = require('terser-webpack-plugin');
var { InjectManifest } = require('workbox-webpack-plugin');
var HashOutputPlugin = require('webpack-plugin-hash-output');

module.exports = {
  mode: 'production',
  entry: {
    main: './js/index',
    boot: './js/boot'
  },
  output: {
    filename: '[name].[chunkhash].js',
    chunkFilename: '[name].[chunkhash].js',
    publicPath: '/'
  },
  resolve: {
    alias: {
      components: path.resolve(__dirname, 'js/components'),
      containers: path.resolve(__dirname, 'js/containers'),
      state: path.resolve(__dirname, 'js/state'),
      utils: path.resolve(__dirname, 'js/utils')
    }
  },
  module: {
    rules: [
      {
        test: /\.js$/,
        loader: 'eslint-loader',
        exclude: /node_modules/,
        enforce: 'pre',
        options: {
          fix: true
        }
      },
      {
        test: /\.js$/,
        loader: 'babel-loader',
        exclude: /node_modules/
      },
      {
        test: /\.css$/,
        use: [
          MiniCssExtractPlugin.loader,
          {
            loader: 'css-loader',
            options: {
              modules: false
            }
          },
          {
            loader: 'postcss-loader',
            options: {
              ident: 'postcss',
              plugins: [
                require('postcss-flexbugs-fixes'),
                postcssPresetEnv({
                  autoprefixer: {
                    flexbox: 'no-2009'
                  }
                }),
                cssnano({
                  discardUnused: {
                    fontFace: false
                  }
                })
              ]
            }
          }
        ]
      }
    ]
  },
  plugins: [
    new MiniCssExtractPlugin({
      filename: '[name].[contenthash].css',
      chunkFilename: '[name].[contenthash].css'
    }),
    new InjectManifest({
      swSrc: './js/sw.js',
      importWorkboxFrom: 'local',
      globDirectory: './public',
      globPatterns: ['*', 'font/*.woff2'],
      exclude: [
        /\.map$/,
        /^manifest.*\.js(?:on)?$/,
        /^boot.*\.js$/,
        /^runtime.*\.js$/
      ]
    }),
    new HashOutputPlugin()
  ],
  optimization: {
    minimizer: [
      new TerserPlugin({
        terserOptions: {
          safari10: true
        },
        cache: true,
        parallel: true
      })
    ],
    splitChunks: {
      chunks: 'all',
      cacheGroups: {
        styles: {
          test: /\.css$/,
          chunks: 'all'
        }
      }
    },
    runtimeChunk: 'single'
  }
};
