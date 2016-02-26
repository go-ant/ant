angular
  .module('app')
  .factory('postService', postService);

postService.$inject = ['$http', 'commService'];
function postService($http, commService) {
  let service = {
    list: list,
    get: get,
    create: create,
    edit: edit,
    remove: remove
  };
  return service;

  function list(opts) {
    return $http.get('../api/backend/posts', {params: {limit: opts.limit || 15, page: opts.page || 1}})
      .then(commService.ajaxContentCallBack);
  }

  function get(id, include) {
    return $http.get('../api/backend/posts/' + id, {params: {include: include}})
      .then(commService.ajaxContentCallBack);
  }

  function create(post) {
    return $http.post('../api/backend/posts', post)
      .then(commService.ajaxContentCallBack);
  }

  function edit(post) {
    return $http.put('../api/backend/posts/' + post.id, post)
      .then(commService.ajaxContentCallBack);
  }

  function remove(id) {
    return $http.delete('../api/backend/posts/' + id)
      .then(commService.ajaxContentCallBack);
  }

}