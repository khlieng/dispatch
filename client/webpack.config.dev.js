var path = require('path');
var webpack = require('webpack');

function dir(p) {
  return path.resolve(__dirname, p);
}

function jsDir(p) {
  return path.resolve(__dirname, 'src/js', p);
}

module.exports = {
  devtool: 'eval',
  entry: [
    'react-hot-loader/patch',
    'webpack-hot-middleware/client',
    './src/js/index'
  ],
  output: {
    path: dir('dist'),
    filename: 'bundle.js',
    publicPath: '/'
  },
  resolve: {
    alias: {
      components: jsDir('components'),
      containers: jsDir('containers'),
      state: jsDir('state'),
      util: jsDir('util')
    }
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
      DEV: true,
    }),
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoEmitOnErrorsPlugin()
  ]
};
