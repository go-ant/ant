//
//  Github Extension (WIP)
//  ~~strike-through~~   ->  <del>strike-through</del>
//

(function () {
  'use strict';

  var github = function () {
    return [
      {
        // strike-through
        // NOTE: showdown already replaced "~" with "~T", so we need to adjust accordingly.
        type:    'lang',
        regex:   '(~T){2}([^~]+)(~T){2}',
        replace: function (match, prefix, content) {
          return '<del>' + content + '</del>';
        }
      }
    ];
  };

  // Client-side export
  if (typeof window !== 'undefined' && window.showdown && window.showdown.extensions) {
    window.showdown.extensions.github = github;
  }
  // Server-side export
  if (typeof module !== 'undefined') {
    module.exports = github;
  }

}());