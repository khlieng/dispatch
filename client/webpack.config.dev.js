var path = require('path');
var webpack = require('webpack');
var postcssPresetEnv = require('postcss-preset-env');

module.exports = {
  mode: 'development',
  entry: ['webpack-hot-middleware/client', './js/index'],
  output: {
    filename: 'bundle.js',
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
      { test: /\.js$/, loader: 'babel-loader', exclude: /node_modules/ },
      {
        test: /\.css$/,
        use: [
          'style-loader',
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
                require('postcss-flexbugs-fixes'),
                postcssPresetEnv({
                  autoprefixer: {
                    flexbox: 'no-2009'
                  }
                })
              ]
            }
          }
        ]
      }
    ]
  },
  plugins: [new webpack.HotModuleReplacementPlugin()]
};
