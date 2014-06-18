// Karma configuration
module.exports = function(config) {
  config.set({
    frameworks: [
      'jasmine',
      'browserify'
    ],
    files: [
      'public/javascripts/vendor/sockjs-0.3.js',
      'src/test/**/*.spec.js'
    ],
    exclude: [],
    preprocessors: {
      'src/test/**/*': [
        'browserify'
      ]
    },
    browserify: {},
    reporters: [
      'progress'
    ],
    port: 9876,
    runnerPort: 9100,
    colors: true,
    logLevel: config.LOG_WARN,
    autoWatch: true,
    captureTimeout: 60000,
    singleRun: false,
    browsers: [
      'Chrome'
    ],
  });
};
