angular
  .module('app')
  .factory('settingService', settingService);

settingService.$inject = ['$http', 'commService'];
function settingService($http, commService) {
  let service = {
    list: list,
    edit: edit
  };
  return service;

  function list() {
    return $http.get('../api/backend/settings')
      .then(commService.ajaxContentCallBack);
  }

  function edit(setting) {
    return $http.put('../api/backend/settings', setting)
      .then(commService.ajaxContentCallBack);
  }
}