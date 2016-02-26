angular
  .module('app')
  .factory('appService', appService);

appService.$inject = ['$http', '$q', 'commService'];
function appService($http, $q, commService) {
  let d = $q.defer()
    , service = {
      init: init,
      defer: d,
      promise: d.promise
    };
  return service;

  function init() {
    return $http.get('../api/backend/app')
      .then(commService.ajaxContentCallBack);
  }
}
