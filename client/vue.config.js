const { defineConfig } = require("@vue/cli-service");
module.exports = defineConfig({
  devServer: {
    proxy: {
      "/shorten": {
        target: "https://e7ast1c-shrty.herokuapp.com",
        secure: false,
        changeOrigin: true,
      },
      "/expand": {
        target: "https://e7ast1c-shrty.herokuapp.com",
        secure: false,
        changeOrigin: true,
      },
    },
    // },
  },
  transpileDependencies: true,
  chainWebpack: (config) => {
    const svgRule = config.module.rule("svg");

    svgRule.uses.clear();

    svgRule
      .oneOf("inline")
      .resourceQuery(/inline/)
      .use("babel-loader")
      .loader("babel-loader")
      .end()
      .use("vue-svg-loader")
      .loader("vue-svg-loader")
      .end()
      .end()
      .oneOf("external")
      .use("file-loader")
      .loader("file-loader")
      .options({
        name: "assets/[name].[hash:8].[ext]",
      });
    config.plugin("html").tap((args) => {
      args[0].title = "Shorty";
      return args;
    });
  },
});
