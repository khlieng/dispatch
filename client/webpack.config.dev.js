var path = require('path');
var webpack = require('webpack');

module.exports = {
  mode: 'development',
  entry: ['webpack-hot-middleware/client', './src/js/index'],
  output: {
    filename: 'bundle.js',
    publicPath: '/'
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
      { test: /\.css$/, loader: 'style-loader!css-loader' }
    ]
  },
  plugins: [new webpack.HotModuleReplacementPlugin()]
};
