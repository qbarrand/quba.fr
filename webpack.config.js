const path = require('path');

const FaviconsWebpackPlugin = require('favicons-webpack-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = {
    entry: path.resolve(__dirname, 'web-src/app.ts'),
    plugins: [
        new FaviconsWebpackPlugin(
            path.resolve(__dirname, 'web-src/img/favicon.png')
        ),
        new HtmlWebpackPlugin({
            base: { target: '_blank' },
            meta: {
                description: 'Software Engineer working at Red Hat on OpenShift. Likes Go, Linux, containers and networking. Previously: Swisscom, CERN, SNCF Réseau.',
                keywords: 'software engineer infrastructure devops sre reliability networks go golang rust containers kubernetes docker linux big-o zürich switzerland',
                viewport: 'width=device-width, initial-scale=1',
                'theme-color': '' // set at runtime
            },
            title: 'Quentin Barrand | Software Engineer',
            template: path.resolve(__dirname, 'web-src/index.html'),
        })
    ],
    // Enable sourcemaps for debugging webpack's output.
    devtool: "source-map",
    resolve: {
        extensions: [".js", ".ts"],
    },

    module: {
        rules: [
            {
                test: /\.css$/i,
                use: ['style-loader', 'css-loader'],
            },
            {
                test: /\.ts$/,
                use: 'ts-loader',
                exclude: /node_modules/,
            },
        ]
    },
    output: {
        filename: '[name].[contenthash].js',
        path: path.resolve(__dirname, 'dist'),
    },
};