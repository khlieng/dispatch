var path = require('path');
var exec = require('child_process').exec;
var url = require('url');

var gulp = require('gulp');
var gutil = require('gulp-util');
var htmlmin = require('gulp-htmlmin');
var nano = require('gulp-cssnano');
var autoprefixer = require('gulp-autoprefixer');
var gzip = require('gulp-gzip');
var concat = require('gulp-concat');
var cache = require('gulp-cached');
var express = require('express');
var proxy = require('express-http-proxy');
var webpack = require('webpack');

gulp.task('html', function() {
  return gulp.src('src/*.html')
    .pipe(htmlmin({
      collapseWhitespace: true,
      removeAttributeQuotes: true
    }))
    .pipe(gulp.dest('dist'));
});

gulp.task('css', function() {
  return gulp.src(['src/css/fontello.css', 'src/css/style.css'])
    .pipe(concat('bundle.css'))
    .pipe(autoprefixer())
    .pipe(nano())
    .pipe(gulp.dest('dist'));
});

gulp.task('js', function(cb) {
  var config = require('./webpack.config.prod.js');
  var compiler = webpack(config);

  process.env['NODE_ENV'] = 'production';

  compiler.run(function(err, stats) {
    if (err) throw new gutil.PluginError('webpack', err);

    gutil.log('[webpack]', stats.toString({
      colors: true
    }));

    if (stats.hasErrors()) process.exit(1);

    cb();
  });
});

gulp.task('fonts', function() {
  return gulp.src('src/font/*')
    .pipe(gulp.dest('dist/font'));
});

gulp.task('config', function() {
  return gulp.src('../config.default.toml')
    .pipe(gulp.dest('dist/gz'));
});

function compress() {
  return gulp.src(['dist/**/!(*.gz)', '!dist/{gz,gz/**}'])
    .pipe(gzip())
    .pipe(gulp.dest('dist/gz'));
}

gulp.task('gzip', ['css', 'js', 'fonts'], compress);
gulp.task('gzip:dev', ['css', 'fonts'], compress);

gulp.task('bindata', ['gzip', 'config'], function(cb) {
  exec('go-bindata -nomemcopy -nocompress -pkg assets -o ../assets/bindata.go -prefix "dist/gz" dist/gz/...', cb);
});

gulp.task('bindata:dev', ['gzip:dev', 'config'], function(cb) {
  exec('go-bindata -debug -pkg assets -o ../assets/bindata.go -prefix "dist/gz" dist/gz/...', cb);
});

gulp.task('dev', ['css', 'fonts', 'config', 'gzip:dev', 'bindata:dev'], function() {
  gulp.watch('src/css/*.css', ['css']);

  var config = require('./webpack.config.dev.js');
  var compiler = webpack(config);
  var app = express();

  app.use(require('webpack-dev-middleware')(compiler, {
    noInfo: true,
    publicPath: config.output.publicPath
  }));

  app.use(require('webpack-hot-middleware')(compiler));

  app.use('/', express.static('dist'));

  app.use('*', proxy('localhost:1337', {
    forwardPath: function(req, res) {
      return url.parse(req.url).path;
    }
  }));

  app.listen(3000, function (err) {
    if (err) {
      console.log(err);
      return;
    }

    console.log('Listening at http://localhost:3000');
  });
});

gulp.task('build', ['css', 'js', 'fonts', 'config', 'gzip', 'bindata']);

gulp.task('default', ['dev']);
