module.exports = {
  context: __dirname,
  mode: 'development',
  entry: './web/src/index.js',
  output: {
    path: __dirname + '/web/js',
    filename: 'index.js'
  },
  devtool: 'none',
  module: {
    rules: [
      {
        test: /\.js$/,
        exclude: /node_modules/,
        use: {
          loader: "babel-loader"
        }
      }
    ]
  }
}
