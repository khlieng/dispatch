var exec = require('child_process').exec;

var gulp = require('gulp');
var gutil = require('gulp-util');
var gulpif = require('gulp-if');
var minifyHTML = require('gulp-minify-html');
var minifyCSS = require('gulp-minify-css');
var autoprefixer = require('gulp-autoprefixer');
var uglify = require('gulp-uglify');
var gzip = require('gulp-gzip');
var concat = require('gulp-concat');
var eslint = require('gulp-eslint');
var browserify = require('browserify');
var source = require('vinyl-source-stream');
var streamify = require('gulp-streamify');
var babelify = require('babelify');
var strictify = require('strictify');
var watchify = require('watchify');
var merge = require('merge-stream');
var cache = require('gulp-cached');

var argv = require('yargs')
    .alias('p', 'production')
    .argv;

if (argv.production) {
    process.env['NODE_ENV'] = 'production';
}

var deps = Object.keys(require('./package.json').dependencies);

gulp.task('html', function() {
    return gulp.src('src/*.html')
        .pipe(minifyHTML())
        .pipe(gulp.dest('dist'));
});

gulp.task('css', function() {
    return gulp.src(['src/css/fontello.css', 'src/css/style.css'])
        .pipe(concat('bundle.css'))
        .pipe(autoprefixer())
        .pipe(minifyCSS())
        .pipe(gulp.dest('dist'));
});

gulp.task('js', function() {
    return js(false);
});

function js(watch) {
    var bundler = browserify('./src/js/app.js', {
        debug: !argv.production,
        transform: [babelify, strictify],
        cache: {},
        packageCache: {},
        fullPaths: watch
    });

    bundler.external(deps);

    var rebundle = function() {
        return bundler.bundle()
            .on('error', gutil.log)
            .pipe(source('bundle.js'))
            .pipe(gulpif(argv.production, streamify(uglify())))
            .pipe(gulp.dest('dist'));
    };

    if (watch) {
        bundler = watchify(bundler);
        bundler.on('update', rebundle);
        bundler.on('log', gutil.log);
    }

    var vendorBundler = browserify({
        debug: !argv.production,
        require: deps
    });

    var vendor = vendorBundler.bundle()
        .on('error', gutil.log)
        .pipe(source('vendor.js'))
        .pipe(gulpif(argv.production, streamify(uglify())))
        .pipe(gulp.dest('dist'));
    
    return merge(rebundle(), vendor);
}

gulp.task('lint', function() {
    return gulp.src('src/js/**/*.{js,jsx}')
        .pipe(cache('lint'))
        .pipe(eslint())
        .pipe(eslint.format())
        .pipe(eslint.failOnError());
});

gulp.task('fonts', function() {
    return gulp.src('src/font/*')
        .pipe(gulp.dest('dist/font'));
});

gulp.task('config', function() {
    return gulp.src('../config.default.toml')
        .pipe(gulp.dest('dist/gz'));
});

gulp.task('gzip', ['html', 'css', 'js', 'fonts'], function() {
    return gulp.src(['dist/**/!(*.gz)', '!dist/{gz,gz/**}'])
        .pipe(gzip())
        .pipe(gulp.dest('dist/gz'));
});

gulp.task('gzip:watch', function() {
    return gulp.src('dist/**/*.{html,css,js}')
        .pipe(cache('gzip'))
        .pipe(gzip())
        .pipe(gulp.dest('dist/gz'));
});

function bindata(cb) {
    if (argv.production) {
        exec('go-bindata -nomemcopy -nocompress -pkg assets -o ../assets/bindata.go -prefix "dist/gz" dist/gz/...', cb);
    } else {
        exec('go-bindata -debug -pkg assets -o ../assets/bindata.go -prefix "dist/gz" dist/gz/...', cb);
    }
}

gulp.task('bindata', ['gzip', 'config'], bindata);
gulp.task('bindata:watch', ['gzip:watch'], bindata);

gulp.task('watch', ['default'], function() {
    gulp.watch('dist/**/*.{html,css,js}', ['gzip:watch', 'bindata:watch'])
    gulp.watch('src/*.html', ['html']);
    gulp.watch('src/css/*.css', ['css']);
    gulp.watch('src/js/**/*.{js,jsx}', ['lint']);
    return js(true);
});

gulp.task('default', ['html', 'css', 'js', 'lint', 'fonts', 'config', 'gzip', 'bindata']);