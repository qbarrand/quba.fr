var grunt = require('grunt');

module.exports = function(grunt) {
  grunt.initConfig({
    cssmin: {
      // options: {
      //   shorthandCompacting: false,
      //   roundingPrecision: -1
      // },
      custom: {
        files: {
          'css/bundle.min.css': [
            'css/tooltipster.css',
            'css/main.css',
            'css/custom.css'
          ]
        }
      }
    },
    uglify: {
      min: {
        files: {
          'js/bundle.min.js': [
            // jQuery and plugins
            'js/jquery.min.js',
            'js/jquery.poptrox.min.js',
            'js/jquery.tooltipster.min.js',
            'js/jquery.validate.min.js',

            // skel and plugins
            'js/skel.min.js',
            'js/skel-layout.min.js',
            'js/util.js',
            'js/main.js'
          ]
        }
      }
    },
    htmlmin: {
      dist: {
        options: {
          removeComments: true,
          collapseWhitespace: true
        },
        files: {
          'index.html': 'index_dev.html',
        }
      }
    }
  });

  grunt.loadNpmTasks('grunt-contrib-uglify');
  grunt.loadNpmTasks('grunt-contrib-cssmin');
  grunt.loadNpmTasks('grunt-contrib-htmlmin');

  grunt.registerTask('default', ['uglify', 'cssmin', 'htmlmin']);
}
