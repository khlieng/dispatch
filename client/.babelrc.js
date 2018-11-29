module.exports = {
  presets: [
    [
      '@babel/preset-env',
      {
        modules: false,
        loose: true
      }
    ],
    '@babel/preset-react'
  ],
  plugins: [
    ['@babel/plugin-proposal-class-properties', { loose: true }],
    '@babel/plugin-proposal-export-default-from',
    '@babel/plugin-proposal-export-namespace-from',
    '@babel/plugin-syntax-dynamic-import',
    'react-hot-loader/babel'
  ],
  env: {
    test: {
      plugins: ['@babel/plugin-transform-modules-commonjs']
    },
    production: {
      plugins: [
        '@babel/plugin-transform-react-inline-elements',
        '@babel/plugin-transform-react-constant-elements'
      ]
    }
  }
};
