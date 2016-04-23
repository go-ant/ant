var gulp = require('gulp')
  , babel = require('gulp-babel')
  , uglify = require('gulp-uglifyjs')
  , sass = require('gulp-sass')
  , htmlmin = require('gulp-htmlmin')
  , minify = require('gulp-minify-css')
  , rename = require('gulp-rename')
  , replace = require('gulp-replace')
  , clean = require('gulp-clean')
  , n2a = require('gulp-native2ascii')
  , templateCache = require('gulp-angular-templatecache')
  , po2json = require('gulp-po2json-angular-translate');

var paths = {
  html: 'core/client/*.html',
  i18n: 'content/i18n/*.po',
  vendor: [
    'core/assets/js/lib/jquery/jquery.js',
    'core/assets/js/lib/moment/moment.js',
    'core/assets/js/lib/showdown/showdown.js',
    'core/assets/js/lib/showdown/showdown-*.js',
    'core/assets/js/lib/loadimage/*.js',
    'core/assets/js/lib/angularjs/angular.js',
    'core/assets/js/lib/angularjs/angular-*.js'
  ],
  scriptsMain: [
    'core/client/main/app.js',
    'core/client/common.js',
    'core/client/app.translate.js',
    'core/client/main/**/**.js'
  ],
  scriptsInstall: [
    'core/client/install/app.js',
    'core/client/common.js',
    'core/client/app.translate.js',
    'core/client/main/user/user.service.js'
  ],
  scriptsLogin: [
    'core/client/login/app.js',
    'core/client/common.js',
    'core/client/app.translate.js',
    'core/client/main/user/user.service.js'
  ],
  templates: [
    'core/client/main/views/**/*'
  ],
  templatesCache: [
    'built/assets/js/templates.js'
  ],
  styles: [
    'core/assets/sass/**/*.scss'
  ]
};

var optsUglify = {
  compress: true,
  output: {
    ascii_only: true,
    beautify: false
  }
};

gulp.task('copy', function() {
  return gulp.src('core/assets/css/**/**')
    .pipe(gulp.dest('built/assets/css'))
});

gulp.task('html-min', function() {
  var options = {
    removeComments: true,
    collapseWhitespace: true,
    removeEmptyAttributes: true,
    removeScriptTypeAttributes: true,
    removeStyleLinkTypeAttributes: true,
    ignoreCustomFragments: [ /\{%[^%]+%\}/g ],
    minifyJS: true,
    minifyCSS: true
  };
  return gulp.src(paths.html)
    .pipe(htmlmin(options))
    .pipe(gulp.dest('built/'));
});

gulp.task('template-cache', function() {
  return gulp.src(paths.templates)
    .pipe(templateCache({
      module: 'app'
    }))
    .pipe(gulp.dest('built/assets/js/'));
});

gulp.task('scripts-vendor', function() {
  return gulp.src(paths.vendor)
    .pipe(uglify('vendor.min.js', optsUglify))
    .pipe(gulp.dest('built/assets/js'));

});

gulp.task('scripts-install', function() {
  gulp.src(paths.scriptsInstall)
    .pipe(babel())
    .pipe(uglify('install.min.js', optsUglify))
    .pipe(gulp.dest('built/assets/js'));
});

gulp.task('scripts-login', function() {
  gulp.src(paths.scriptsLogin)
    .pipe(babel())
    .pipe(uglify('login.min.js', optsUglify))
    .pipe(gulp.dest('built/assets/js'));
});

gulp.task('scripts-main', ['template-cache'], function() {
  var packageScripts = paths.scriptsMain.concat(paths.templatesCache);
  return gulp.src(packageScripts)
    .pipe(babel())
    .pipe(uglify('goant.min.js', optsUglify))
    .pipe(gulp.dest('built/assets/js'));
});

gulp.task('scripts', ['scripts-main'], function() {
  gulp.src('built/assets/js/templates.js')
    .pipe(clean({force: true}));
});

gulp.task('styles', function() {
  gulp.src('core/assets/sass/themes/*.scss')
    .pipe(sass())
    .pipe(minify({processImport: false}))
    .pipe(replace('themes/default/assets/', ''))
    .pipe(rename({ extname: '.min.css' }))
    .pipe(gulp.dest('built/assets/css/themes/'));


  gulp.src('core/assets/sass/vendor.scss')
    .pipe(sass())
    .pipe(minify({processImport: false}))
    .pipe(replace('themes/default/assets/', ''))
    .pipe(rename({ extname: '.min.css' }))
    .pipe(gulp.dest('built/assets/css/'));
});


gulp.task('i18n', function () {
  return gulp.src(paths.i18n)
    .pipe(po2json())
    .pipe(n2a({reverse: false}))
    .pipe(gulp.dest('built/assets/i18n/'));
});


gulp.task('watch', function() {
  gulp.watch(paths.html, ['html-min']);
  gulp.watch(paths.styles, ['styles']);
  gulp.watch(paths.vendor, ['scripts-vendor']);
  gulp.watch(paths.scriptsInstall, ['scripts-install']);
  gulp.watch(paths.scriptsLogin, ['scripts-login']);
  gulp.watch(paths.scriptsMain, ['scripts']);
  gulp.watch(paths.templates, ['scripts']);
  gulp.watch(paths.i18n, ['i18n']);
});

gulp.task('production', ['copy', 'html-min', 'styles', 'i18n', 'scripts-vendor', 'scripts-install', 'scripts-login', 'scripts']);

gulp.task('default', ['production', 'watch']);
