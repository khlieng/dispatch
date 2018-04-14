var path = require('path');
var webpack = require('webpack');
var MiniCssExtractPlugin = require('mini-css-extract-plugin');
var autoprefixer = require('autoprefixer');
var cssnano = require('cssnano');

module.exports = {
  mode: 'production',
  entry: ['./src/js/index'],
  output: {
    filename: 'bundle.js'
  },
  resolve: {
    alias: {
      components: path.resolve(__dirname, 'src/js/components'),
      containers: path.resolve(__dirname, 'src/js/containers'),
      state: path.resolve(__dirname, 'src/js/state'),
      utils: path.resolve(__dirname, 'src/js/utils')
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
      { test: /\.js$/, loader: 'babel-loader', exclude: /node_modules/ },
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
              plugins: [
                autoprefixer,
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
      filename: 'bundle.css'
    })
  ],
  optimization: {
    splitChunks: {
      cacheGroups: {
        styles: {
          test: /\.css$/,
          chunks: 'all'
        }
      }
    }
  }
};
