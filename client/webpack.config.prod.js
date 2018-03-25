var webpack = require('webpack');

module.exports = {
  mode: 'production',
  entry: [
    './src/js/index'
  ],
  output: {
    filename: 'bundle.js'
  },
  module: {
    rules: [
      { test: /\.js$/, loader: 'eslint-loader', exclude: /node_modules/, enforce: 'pre' },
      { test: /\.js$/, loader: 'babel-loader', exclude: /node_modules/ },
      { test: /\.css$/, loader: 'style-loader!css-loader' }
    ]
  },
  plugins: [
    new webpack.DefinePlugin({
      DEV: false
    })
  ]
};
