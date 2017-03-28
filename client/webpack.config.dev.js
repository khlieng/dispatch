var path = require('path');
var webpack = require('webpack');

module.exports = {
  devtool: 'eval',
  entry: [
    'webpack-hot-middleware/client',
    './src/js/index'
  ],
  output: {
    path: path.join(__dirname, 'dist'),
    filename: 'bundle.js',
    publicPath: '/'
  },
  module: {
    rules: [
      { test: /\.js$/, loader: 'eslint-loader', exclude: /node_modules/, enforce: 'pre' },
      { test: /\.js$/, loader: 'babel-loader', exclude: /node_modules/ },
      { test: /\.css$/, loader: 'style-loader!css-loader' },
      { test: /node_modules/, loader: 'ify-loader' }
    ]
  },
  plugins: [
    new webpack.DefinePlugin({
      DEV: true,
    }),
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoEmitOnErrorsPlugin()
  ]
};
