// GoAnt Image Preview
//
// Manages the conversion of image markdown `![]()` from markdown into the HTML image preview
// This provides a dropzone and other interface elements for adding images

(function (extension) {
  'use strict';

  if (typeof showdown === 'object') {
    showdown.extension('antimagepreview', extension());
  } else {
    // showdown was not found so we throw
    throw Error('Could not find showdown library');
  }

}(function() {

  return [
    // ![] image syntax
    {
      type: 'lang',
      filter: function (text) {
        var imageMarkdownRegex = /^!(?:\[([^\n\]]*)\])(?:\(([^\n\]]*)\))?$/gim,
        /* regex from isURL in node-validator. Yum! */
        uriRegex = /^(?!mailto:)(?:(?:https?|ftp):\/\/)?(?:\S+(?::\S*)?@)?(?:(?:(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z\u00a1-\uffff0-9]+-?)*[a-z\u00a1-\uffff0-9]+)(?:\.(?:[a-z\u00a1-\uffff0-9]+-?)*[a-z\u00a1-\uffff0-9]+)*(?:\.(?:[a-z\u00a1-\uffff]{2,})))|localhost)(?::\d{2,5})?(?:\/[^\s]*)?$/i,
        pathRegex = /^(\/)?([^\/\0]+(\/)?)+$/i;

        return text.replace(imageMarkdownRegex, function (match, alt, src) {
          var result = ''
            , output;

          if (src && (src.match(uriRegex) || src.match(pathRegex))) {
            output = '<p><img src="' + src + '" alt="' + alt + '"></p>';
          } else {
            output = '<section class="img-zone post-img-zone">' +
                '<div class="img-desc">Add image of <strong>' + alt + '</strong></div>' +
                '<div class="progress"></div>' +
              '</section>';
          }

          return output;
        });
      }
    }
  ];
}));