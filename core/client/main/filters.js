angular
  .module('app')
  .filter('fromNow', filterFromNow)
  .filter('trustAsHtml', filterTrustAsHtml);

function filterFromNow() {
  return function(date, defValue) {
    return Date.parse(date) < 0 ? defValue: moment(date).fromNow();
  }
}

filterTrustAsHtml.$inject = ['$sce'];
function filterTrustAsHtml($sce) {
  return function (html) {
    return $sce.trustAsHtml(html);
  }
}