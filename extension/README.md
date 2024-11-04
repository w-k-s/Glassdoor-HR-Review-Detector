# Review Filter Extension

A browser extension that helps identify and filter potentially fake reviews on Glassdoor. Compatible with Chrome and Firefox.

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Local Development](#local-development)
- [Testing](#testing)
- [Browser-specific Development](#browser-specific-development)
- [Deployment](#deployment)
- [Environment Variables](#environment-variables)
- [Project Structure](#project-structure)

## Prerequisites

- Node.js (v14 or higher)
- npm (v6 or higher)
- Git
- Chrome and/or Firefox browser
- A code editor (VS Code recommended)

## Installation

1. Clone the repository:
2. Install dependencies:
```bash
npm install
```

3. Create environment configuration:
```bash
cp .env.example .env
```

4. Update `.env` with your API endpoint:
```
API_ENDPOINT=https://your-api-endpoint.com
BROWSER=chrome
```

## Local Development

### Building the Extension

Build the extension for development:
```bash
# For Chrome
npm run build:chrome

# For Firefox
npm run build:firefox
```

The built extension will be in `dist/chrome` or `dist/firefox` respectively.

### Loading in Chrome

1. Open Chrome and navigate to `chrome://extensions`
2. Enable "Developer mode" (top right corner)
3. Click "Load unpacked"
4. Navigate to the `dist/chrome` directory and select it
5. The extension should now be loaded and active

### Loading in Firefox

1. Open Firefox and navigate to `about:debugging`
2. Click "This Firefox" in the left sidebar
3. Click "Load Temporary Add-on"
4. Navigate to the `dist/firefox` directory and select the `manifest.json` file
5. The extension should now be loaded and active

### Development Workflow

1. Make changes to the source code
2. Run the build command for your target browser
3. The extension will need to be reloaded in the browser:
   - Chrome: Click the refresh icon on the extension card
   - Firefox: Click "Reload" next to the extension

## Local Development

### Setting up local environment

1. Create a local environment file:
```bash
cp .env.local.example .env.local
```

2. Update `.env.local` with your local API endpoint:
```
API_ENDPOINT=http://localhost:3000/api/reviews
BROWSER=chrome
NODE_ENV=development
```

### Running locally

1. Start your local API server

2. Build and watch for changes:
```bash
# For Chrome
npm run watch:chrome

# For Firefox
npm run watch:firefox
```

3. Load the extension in your browser:

For Chrome:
- Go to chrome://extensions/
- Enable Developer mode
- Click "Load unpacked"
- Select the `dist/chrome` directory

For Firefox:
- Go to about:debugging#/runtime/this-firefox
- Click "Load Temporary Add-on"
- Select the `manifest.json` file in the `dist/firefox` directory

The extension will automatically rebuild when you make changes. You'll need to:
- Chrome: Click the refresh icon on the extension card
- Firefox: Click "Reload" next to the extension

### Environment-specific builds

```bash
# Local development
npm run build:local:chrome
npm run build:local:firefox

# Development server
npm run build:dev:chrome
npm run build:dev:firefox

# Production
npm run build:prod:chrome
npm run build:prod:firefox

# Watch mode (auto-rebuilds on changes)
npm run watch:chrome
npm run watch:firefox
```

## Testing

### Running Unit Tests

```bash
# Run all unit tests
npm test

# Run unit tests in watch mode
npm test -- --watch

# Run tests with coverage
npm test -- --coverage
```

### Running E2E Tests

```bash
# Run all E2E tests
npm run test:e2e

# Run specific E2E test file
npm run test:e2e tests/e2e/reviewFilter.test.js
```

### Test Structure

- `tests/unit/`: Contains unit tests for individual components
- `tests/e2e/`: Contains end-to-end tests using Puppeteer
- `tests/mocks/`: Contains mock data and utilities for testing

## Browser-specific Development

### Chrome Development Notes

- Uses Manifest V3
- Support for modern JavaScript features
- Direct access to chrome.* APIs

### Firefox Development Notes

- Uses Manifest V2
- May require polyfills for certain features
- Uses browser.* API namespace

## Deployment

### Chrome Web Store Deployment

1. Create a production build:
```bash
npm run build:chrome -- --env production
```

2. Zip the contents of `dist/chrome`:
```bash
cd dist/chrome && zip -r ../chrome-extension.zip ./*
```

3. Submit to Chrome Web Store:
   - Go to [Chrome Developer Dashboard](https://chrome.google.com/webstore/developer/dashboard)
   - Click "New Item"
   - Upload the `chrome-extension.zip` file
   - Fill in required information
   - Submit for review

### Firefox Add-ons Deployment

1. Create a production build:
```bash
npm run build:firefox -- --env production
```

2. Zip the contents of `dist/firefox`:
```bash
cd dist/firefox && zip -r ../firefox-addon.zip ./*
```

3. Submit to Firefox Add-ons:
   - Go to [Firefox Add-ons Developer Hub](https://addons.mozilla.org/developers/)
   - Click "Submit a New Add-on"
   - Upload the `firefox-addon.zip` file
   - Fill in required information
   - Submit for review

### Environment Variables

Different environment files are used based on the build command:

- `.env.local`: Local development (localhost)
- `.env.development`: Development server
- `.env.production`: Production server

Each environment can have different settings for:
- API endpoints
- Debug flags
- Feature toggles
- Other configuration

## Project Structure

```
.
├── src/                  # Source code
│   ├── manifest.v2.json  # Firefox manifest
│   ├── manifest.v3.json  # Chrome manifest
│   ├── content/         # Content scripts
│   ├── api/            # API services
│   └── styles/         # CSS styles
├── tests/              # Test files
│   ├── e2e/           # End-to-end tests
│   ├── unit/          # Unit tests
│   └── mocks/         # Test mocks
├── dist/              # Built extensions
├── webpack.config.js  # Build configuration
├── jest.config.js     # Test configuration
└── package.json       # Project dependencies
```