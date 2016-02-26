angular
  .module('app')
  .factory('themeService', themeService);

themeService.$inject = ['$http', 'commService'];
function themeService($http, commService) {
  let service = {
    list: list
  };
  return service;

  function list() {
    return $http.get('../api/backend/themes')
      .then(commService.ajaxContentCallBack);
  }
}