let cssConfig = {};
const path = require("path");
const vueSrc = "./src";

if (process.env.NODE_ENV == "production") {
    cssConfig = {
        extract: {
            filename: "[name].css",
            chunkFilename: "[name].css"
        }
    };
}

module.exports = {
    css: cssConfig,
    configureWebpack: {
        output: {
            filename: "[name].js"
        },
        optimization: {
            splitChunks: false
        },
        resolve: {
            alias: {
                "@": path.resolve(__dirname, vueSrc)
            },
            extensions: ['.js', '.vue', '.json']
        }
    },
    devServer: {
        disableHostCheck: true
    }
};
