const path = require('path');
const CompressionPlugin = require("compression-webpack-plugin");
const FaviconsWebpackPlugin = require('favicons-webpack-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin');
const zlib = require("zlib");

const DIST_DIR_NAME = 'dist'

module.exports = {
    entry: path.resolve(__dirname, 'web-src/app.js'),
    plugins: [
        new CompressionPlugin(), // gzip by default
        new CompressionPlugin({
            filename: "[path][base].br",
            algorithm: "brotliCompress",
            compressionOptions: {
                params: {
                    [zlib.constants.BROTLI_PARAM_QUALITY]: 11,
                },
            },
        }),
        new FaviconsWebpackPlugin(
            path.resolve(__dirname, 'web-src/img/favicon.png')
        ),
        new HtmlWebpackPlugin({
            base: { target: '_blank' },
            meta: {
                description: 'Software Engineer working at NVIDA on Software-Defined Networking. Likes Go, Linux, containers and networking. Previously: Red Hat, Swisscom, CERN, SNCF Réseau.',
                keywords: 'software engineer infrastructure devops sre reliability networks go golang rust containers kubernetes docker linux big-o zürich switzerland',
                viewport: 'width=device-width, initial-scale=1',
                'theme-color': '' // set at runtime
            },
            title: 'Quentin Barrand | Software Engineer',
            template: path.resolve(__dirname, 'web-src/index.html'),
        })
    ],
    module: {
        rules: [
            {
                test: /\.css$/i,
                use: ['style-loader', 'css-loader'],
            }
        ]
    },
    output: {
        clean: { keep: /backgrounds/},
        filename: '[name].[contenthash].js',
        path: path.resolve(__dirname, DIST_DIR_NAME),
    },
    optimization: {
        minimize: true,
    }
};
