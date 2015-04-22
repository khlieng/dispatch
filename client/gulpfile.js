var gulp = require('gulp');
var gulpif = require('gulp-if');
var minifyHTML = require('gulp-minify-html');
var minifyCSS = require('gulp-minify-css');
var autoprefixer = require('gulp-autoprefixer');
var uglify = require('gulp-uglify');
var gzip = require('gulp-gzip');
var browserify = require('browserify');
var source = require('vinyl-source-stream');
var streamify = require('gulp-streamify');
var babelify = require('babelify');
var strictify = require('strictify');
var watchify = require('watchify');

var argv = require('yargs')
    .alias('p', 'production')
    .argv;

if (argv.production) {
    process.env['NODE_ENV'] = 'production';
}

gulp.task('html', function() {
    gulp.src('src/*.html')
        .pipe(minifyHTML())
        .pipe(gulp.dest('dist'));
});

gulp.task('css', function() {
    gulp.src('src/css/*.css')
        .pipe(autoprefixer())
        .pipe(minifyCSS())
        .pipe(gulp.dest('dist/css'));
});

gulp.task('js', function() {
    return js(false);
});

function js(watch) {
    var bundler, rebundle;
    bundler = browserify('./src/js/app.js', {
        debug: !argv.production,
        cache: {},
        packageCache: {},
        fullPaths: watch
    });

    if (watch) {
        bundler = watchify(bundler);
    }

    bundler
        .transform(babelify)
        .transform(strictify);

    rebundle = function() {
        var stream = bundler.bundle();
        stream.on('error', console.log);
        return stream
            .pipe(source('bundle.js'))
            .pipe(gulpif(argv.production, streamify(uglify())))
            .pipe(gulp.dest('dist'));
    };

    bundler.on('time', function(time) {
        console.log('JS bundle: ' + time + ' ms');
    });
    bundler.on('update', rebundle);
    return rebundle();
}

gulp.task('fonts', function() {
    gulp.src('src/font/*')
        .pipe(gulp.dest('dist/font'));
});

gulp.task('gzip', ['html', 'css', 'js', 'fonts'], function() {
    gulp.src('dist/**/!(*.gz)')
        .pipe(gzip())
        .pipe(gulp.dest('dist'));
});

gulp.task('gzip:watch', function() {
    gulp.src('dist/**/*.{html,css,js}')
        .pipe(gzip())
        .pipe(gulp.dest('dist'));
});

gulp.task('watch', ['default'], function() {
    gulp.watch('dist/**/*.{html,css,js}', ['gzip:watch'])
    gulp.watch('src/*.html', ['html']);
    gulp.watch('src/css/*.css', ['css']);
    return js(true);
});

gulp.task('default', ['html', 'css', 'js', 'fonts', 'gzip']);