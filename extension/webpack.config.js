const path = require('path');
const webpack = require('webpack');
const CopyPlugin = require('copy-webpack-plugin');
const dotenv = require('dotenv');

module.exports = (env) => {
    // Determine which env file to use
    let envFile = '.env';
    if (env.local) {
        envFile = '.env.local';
    } else if (env.development) {
        envFile = '.env.development';
    } else if (env.production) {
        envFile = '.env.production';
    }

    // Load environment variables from the selected file
    const envVars = dotenv.config({ path: envFile }).parsed;

    // Set browser from env var or command line arg
    const browser = env.browser || process.env.BROWSER || 'chrome';
    const manifestVersion = browser === 'firefox' ? 'v2' : 'v3';

    return {
        mode: env.production ? 'production' : 'development',
        devtool: env.production ? 'source-map' : 'inline-source-map',
        entry: {
            content: './src/content/index.js',
            background: './src/background/index.js'
        },
        output: {
            path: path.resolve(__dirname, 'dist', browser),
            filename: '[name].js',
            clean: true
        },
        module: {
            rules: [
                {
                    test: /\.js$/,
                    exclude: /node_modules/,
                    use: {
                        loader: 'babel-loader'
                    }
                },
                {
                    test: /\.css$/,
                    use: ['style-loader', 'css-loader']
                }
            ]
        },
        plugins: [
            new webpack.DefinePlugin({
                'process.env': JSON.stringify({
                    ...envVars,
                    NODE_ENV: env.production ? 'production' : 'development'
                })
            }),
            new CopyPlugin({
                patterns: [
                    {
                        from: `src/manifest.${manifestVersion}.json`,
                        to: 'manifest.json',
                        transform(content) {
                            // Modify manifest name in development to distinguish it
                            if (!env.production) {
                                const manifest = JSON.parse(content.toString());
                                manifest.name += ' (Dev)';
                                return JSON.stringify(manifest, null, 2);
                            }
                            return content;
                        }
                    }
                ]
            })
        ]
    };
};