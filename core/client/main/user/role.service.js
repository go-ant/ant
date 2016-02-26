angular
  .module('app')
  .factory('roleService', roleService);

roleService.$inject = ['$http', 'commService'];
function roleService($http, commService) {
  let service = {
    list: list,
    get: get,
    create: create,
    edit: edit,
    remove: remove,
    listPermission: listPermission
  };
  return service;


  function list() {
    return $http.get('../api/backend/roles')
      .then(commService.ajaxContentCallBack);
  }

  function get(id, include) {
    return $http.get('../api/backend/roles/' + id, { params: { include: include } })
      .then(commService.ajaxContentCallBack);
  }

  function create(role) {
    return $http.post('../api/backend/roles', {
        name: role.name,
        description: role.description,
        permissions: role.permissions.join(',')
      })
      .then(commService.ajaxContentCallBack);
  }

  function edit(role) {
    return $http.put('../api/backend/roles/' + role.id, {
        name: role.name,
        description: role.description,
        permissions: role.permissions.join(',')
      })
      .then(commService.ajaxContentCallBack);
  }

  function remove(id) {
    return $http.delete('../api/backend/roles/' + id)
      .then(commService.ajaxContentCallBack);
  }

  function listPermission() {
    return $http.get('../api/backend/permissions')
      .then(commService.ajaxContentCallBack);
  }
}