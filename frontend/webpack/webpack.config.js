const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");

console.log(path.join(__dirname, "../dist"));
console.log(path.resolve(__dirname, "../dist"));

module.exports = {
  entry: "./src/index.tsx",
  watch: true,
  output: {
    path: path.resolve(__dirname, "../dist"),
    filename: "bundle.js",
  },
  mode: "development",
  devtool: "source-map",
  devServer: {
    inline: true, // Enable watch and live reload
    host: "localhost",
    port: 8080,
    stats: "errors-only",
    hot: true,
    contentBase: "./dist",
    watchContentBase: true,
    proxy: {
      "/api": "http://localhost:8000",
    },
  },
  module: {
    rules: [
      {
        test: /\.(ts|tsx)$/,
        exclude: /node_modules/,
        loader: "awesome-typescript-loader",
      },
      {
        test: /\.css$/,
        use: [
          "style-loader",
          { loader: "css-loader", options: { importLoaders: 1 } },
        ],
      },
      {
        test: /\.(scss|sass)$/,
        loaders: [
          "style-loader",
          { loader: "css-loader", options: { importLoaders: 1 } },
          "sass-loader",
        ],
      },
    ],
  },
  resolve: {
    extensions: [".js", ".ts", ".tsx"],
  },
  plugins: [
    new HtmlWebpackPlugin({
      filename: "index.html", //Name of file in ./dist/
      template: "./src/index.html", //Name of template in ./src
      hash: true,
    }),
  ],
};
